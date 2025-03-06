// Package apiclient provides a client for communicating with the HEXONET backend API.
//
// This package allows two types of communication:
// - Session-based communication: Used for building custom frontends and supports 2FA (Two-Factor Authentication).
// - Sessionless communication: Used for simple command requests.
//
// The package includes the following features:
// - Connection setup: The package supports three connection setups: high performance, default, and OT&E (demo system).
// - API request: The package provides methods for making API requests and handling responses.
// - Session management: The package allows for session login, logout, and session reuse.
// - Debug mode: The package supports enabling and disabling debug mode for logging and output.
// - Proxy configuration: The package allows for setting and retrieving proxy configurations for API communication.
// - User agent customization: The package provides methods for customizing the user agent header.
// - Command parameter handling: The package includes methods for flattening command parameters and automatically converting IDN (Internationalized Domain Name) values to punycode.
// - Pagination support: The package includes methods for requesting next response pages and retrieving all response pages for a given query.
//
// For more information on the available commands, refer to the HEXONET API documentation: https://github.com/hexonet/hexonet-api-documentation/tree/master/API
//
// Example usage:
//     // Create a new APIClient instance
//     client := apiclient.NewAPIClient()
//
//     // Set credentials for API communication
//     client.SetCredentials("username", "password")
//
//     // Make an API request
//     response := client.Request(map[string]interface{}{
//         "COMMAND": "StatusAccount",
//     })
//
//     // Check if the request was successful
//     if response.IsSuccess() {
//         // Process the response data
//         // ...
//     } else {
//         // Handle the error
//         // ...
//     }
//
//     // Close the API session
//     client.Logout()
//
// Note: This package is based on the HEXONET API documentation and is subject to change. Please refer to the documentation for the most up-to-date information.
// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package apiclient contains all you need to communicate with the insanely fast HEXONET backend API.
package apiclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	IDN "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/idntranslator"
	LG "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/logger"
	R "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/response"
	RTM "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responsetemplatemanager"
	SC "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/socketconfig"
)

// CNR_CONNECTION_URL_PROXY represents the url used for the high performance connection setup
const CNR_CONNECTION_URL_PROXY = "http://127.0.0.1/api/call.cgi" //nolint

// CNR_CONNECTION_URL_LIVE represents the url used for the default connection setup
const CNR_CONNECTION_URL_LIVE = "https://api.rrpproxy.net/api/call.cgi" //nolint

// CNR_CONNECTION_URL_OTE represents the url used for the OT&E (demo system) connection setup
const CNR_CONNECTION_URL_OTE = "https://api-ote.rrpproxy.net/api/call.cgi" //nolint

var rtm = RTM.GetInstance()

// APIClient is the entry point class for communicating with the insanely fast HEXONET backend api.
// It allows two ways of communication:
// * session based communication
// * sessionless communication
//
// A session based communication makes sense in case you use it to
// build your own frontend on top. It allows also to use 2FA
// (2 Factor Auth) by providing "otp" in the config parameter of
// the login method.
// A sessionless communication makes sense in case you do not need
// to care about the above and you have just to request some commands.
//
// Possible commands can be found at https://github.com/hexonet/hexonet-api-documentation/tree/master/API
type APIClient struct {
	socketTimeout time.Duration
	socketURL     string
	socketConfig  *SC.SocketConfig
	debugMode     bool
	curlopts      map[string]string
	ua            string
	logger        LG.ILogger
	subUser       string
	roleSeparator string
	client        *http.Client
}

// RequestOptions represents the options for an API request.
type RequestOptions struct {
	SetUserView bool // SetUserView indicates whether to set a data view to a given subuser.
}

// NewRequestOptions creates a new instance of RequestOptions with default values.
func NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		SetUserView: true,
	}
}

// NewAPIClient represents the constructor for struct APIClient.
func NewAPIClient() *APIClient {
	cl := &APIClient{
		debugMode:     false,
		socketTimeout: 300 * time.Second,
		socketURL:     CNR_CONNECTION_URL_LIVE,
		socketConfig:  SC.NewSocketConfig(),
		curlopts:      map[string]string{},
		ua:            "",
		logger:        nil,
		roleSeparator: ":",
		client:        &http.Client{},
	}
	cl.UseLIVESystem()
	cl.SetDefaultLogger()
	return cl
}

// SetDefaultLogger method to use the default mechanism for debug mode outputs
func (cl *APIClient) SetDefaultLogger() *APIClient {
	cl.logger = LG.NewLogger()
	return cl
}

// SetCustomLogger method to use a custom mechanism for debug mode outputs/logging
func (cl *APIClient) SetCustomLogger(logger LG.ILogger) *APIClient {
	cl.logger = logger
	return cl
}

// SetProxy method to set a proxy to use for API communication
func (cl *APIClient) SetProxy(proxy string) *APIClient {
	if len(proxy) == 0 {
		delete(cl.curlopts, "PROXY")
	} else {
		cl.curlopts["PROXY"] = proxy
	}
	return cl
}

// GetProxy method to get the configured proxy to use for API communication
func (cl *APIClient) GetProxy() (string, error) {
	val, exists := cl.curlopts["PROXY"]
	if exists {
		return val, nil
	}
	return "", errors.New("no proxy configuration available")
}

// SetReferer method to set a value for HTTP Header `Referer` to use for API communication
func (cl *APIClient) SetReferer(referer string) *APIClient {
	if len(referer) == 0 {
		delete(cl.curlopts, "REFERER")
	} else {
		cl.curlopts["REFERER"] = referer
	}
	return cl
}

// GetReferer method to get configured HTTP Header `Referer` value
func (cl *APIClient) GetReferer() (string, error) {
	val, exists := cl.curlopts["REFERER"]
	if exists {
		return val, nil
	}
	return "", errors.New("no configuration available for HTTP Header `Referer`")
}

// EnableDebugMode method to enable Debug Output to logger
func (cl *APIClient) EnableDebugMode() *APIClient {
	cl.debugMode = true
	return cl
}

// DisableDebugMode method to disable Debug Output to logger
func (cl *APIClient) DisableDebugMode() *APIClient {
	cl.debugMode = false
	return cl
}

// SetUserView method to set a data view to a given subuser
func (cl *APIClient) SetUserView(uid string) *APIClient {
	cl.subUser = uid
	return cl
}

// ResetUserView method to reset data view back from subuser to user
func (cl *APIClient) ResetUserView() *APIClient {
	cl.subUser = ""
	return cl
}

// UseHighPerformanceConnectionSetup to activate high performance conneciton setup
func (cl *APIClient) UseHighPerformanceConnectionSetup() *APIClient {
	cl.SetURL(CNR_CONNECTION_URL_PROXY)
	return cl
}

// GetPOSTData method to Serialize given command for POST request
// including connection configuration data
func (cl *APIClient) GetPOSTData(cmd map[string]string, secured ...bool) string {
	data := cl.socketConfig.GetPOSTData()
	if len(secured) > 0 && secured[0] {
		re := regexp.MustCompile("s_pw=[^&]+")
		data = re.ReplaceAllString(data, "s_pw=***")
	}
	var tmp strings.Builder
	keys := []string{}
	for key := range cmd {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		val := cmd[key]
		tmp.WriteString(key)
		tmp.WriteString("=")
		val = strings.ReplaceAll(val, "\r", "")
		val = strings.ReplaceAll(val, "\n", "")
		tmp.WriteString(val)
		tmp.WriteString("\n")
	}
	str := tmp.String()
	if len(secured) > 0 && secured[0] {
		re := regexp.MustCompile("PASSWORD=[^\n]+")
		str = re.ReplaceAllString(str, "PASSWORD=***")
	}
	if tmp.String() == "" {
		return strings.TrimSuffix(data, "&")
	}
	str = strings.TrimSuffix(str, "\n")
	return strings.Join([]string{
		data,
		url.QueryEscape("s_command"),
		"=",
		url.QueryEscape(str),
	}, "")
}

// GetURL method to get the API connection url that is currently set
func (cl *APIClient) GetURL() string {
	return cl.socketURL
}

// SetUserAgent method to customize user-agent header (useful for tools that use our SDK)
func (cl *APIClient) SetUserAgent(str string, rv string, modules ...[]string) *APIClient {
	mods := ""
	if len(modules) > 0 {
		for i := 0; i < len(modules[0]); i++ {
			mods += modules[0][i] + " "
		}
	}
	cl.ua = str + " (" + runtime.GOOS + "; " + runtime.GOARCH + "; rv:" + rv + ") " + mods + "go-sdk/" + cl.GetVersion() + " go/" + runtime.Version()
	return cl
}

// GetUserAgent method to return the user agent string
func (cl *APIClient) GetUserAgent() string {
	if len(cl.ua) == 0 {
		cl.ua = "GO-SDK (" + runtime.GOOS + "; " + runtime.GOARCH + "; rv:" + cl.GetVersion() + ") go/" + runtime.Version()
	}
	return cl.ua
}

// GetVersion method to get current module version
func (cl *APIClient) GetVersion() string {
	return "5.0.10"
}

// SaveSession method to apply data to a session for later reuse
// Please save/update that map into user session
func (cl *APIClient) SaveSession(sessionobj map[string]interface{}) *APIClient {
	sessionobj["socketcfg"] = map[string]string{
		"session": cl.socketConfig.GetSession(),
		"login":   cl.socketConfig.GetLogin(),
	}
	return cl
}

// ReuseSession method to reuse given configuration out of a user session
// to rebuild and reuse connection settings
func (cl *APIClient) ReuseSession(sessionobj map[string]interface{}) *APIClient {
	if sessionobj == nil || sessionobj["socketcfg"] == nil {
		return cl
	}
	cfg, ok := sessionobj["socketcfg"].(map[string]string)
	if !ok || cfg["login"] == "" || cfg["session"] == "" {
		return cl
	}
	cl.SetCredentials(cfg["login"])
	cl.socketConfig.SetSession(cfg["session"])
	return cl
}

// SetURL method to set another connection url to be used for API communication
func (cl *APIClient) SetURL(value string) *APIClient {
	cl.socketURL = value
	return cl
}

// SetPersistent method sets the API connection to use a persistent session
func (cl *APIClient) SetPersistent() *APIClient {
	cl.socketConfig.SetPersistent()
	return cl
}

// SetCredentials method to set Credentials to be used for API communication
func (cl *APIClient) SetCredentials(params ...string) *APIClient {
	if len(params) > 0 {
		cl.socketConfig.SetLogin(params[0])
	}
	if len(params) > 1 {
		cl.socketConfig.SetPassword(params[1])
	}
	return cl
}

// SetRoleCredentials method to set Role User Credentials to be used for API communication
func (cl *APIClient) SetRoleCredentials(params ...string) *APIClient {
	if len(params) > 0 {
		uid := params[0]
		if len(params) > 1 && len(params[1]) > 0 {
			role := params[1]
			uid = uid + cl.roleSeparator + role
		}
		if len(params) > 2 {
			pw := params[2]
			return cl.SetCredentials(uid, pw)
		}
		return cl.SetCredentials(uid)
	}
	return cl
}

// Login method to perform API login to start session-based communication
// 1st parameter: one time password
func (cl *APIClient) Login() *R.Response {
	cl.SetPersistent()
	rr := cl.Request(make(map[string]interface{}), &RequestOptions{SetUserView: false})
	cl.socketConfig.SetSession("")
	if rr.IsSuccess() {
		col := rr.GetColumn("SESSIONID")
		if col != nil {
			cl.socketConfig.SetSession(col.GetData()[0])
		}
	}
	return rr
}

// Logout method to perform API logout to close API session in use
func (cl *APIClient) Logout() *R.Response {
	rr := cl.Request(map[string]interface{}{
		"COMMAND": "StopSession",
	}, &RequestOptions{SetUserView: false})
	if rr.IsSuccess() {
		cl.socketConfig.SetSession("")
	}
	return rr
}

// Request method to perform API request using the given command
func (cl *APIClient) Request(cmd map[string]interface{}, opts ...*RequestOptions) *R.Response {
	// Use default RequestOptions if opts is not available
	options := NewRequestOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// Check if SetUserView option is enabled and subUser is set
	if (options.SetUserView) && (len(cl.subUser) > 0) {
		cmd["SUBUSER"] = cl.subUser
	}

	// flatten nested api command bulk parameters
	newcmd := cl.flattenCommand(cmd)
	// auto convert umlaut names to punycode
	newcmd = cl.autoIDNConvert(newcmd)

	// request command to API
	cfg := map[string]string{
		"CONNECTION_URL": cl.socketURL,
	}
	if cl.debugMode {
		fmt.Println("Connecting to: " + cfg["CONNECTION_URL"])
	}
	data := cl.GetPOSTData(newcmd, false)
	secured := cl.GetPOSTData(newcmd, true)

	val, err := cl.GetProxy()
	cl.client.Timeout = cl.socketTimeout

	if err == nil {
		if proxyconfigurl, parsingerr := url.Parse(val); parsingerr == nil {
			cl.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyconfigurl)}
		} else if cl.debugMode {
			fmt.Println("Not able to parse configured Proxy URL: " + val)
		}
	}
	req, err := http.NewRequest("POST", cfg["CONNECTION_URL"], strings.NewReader(data))
	if err != nil {
		tpl := rtm.GetTemplate("httperror")
		r := R.NewResponse(tpl, newcmd, cfg)
		if cl.debugMode {
			cl.logger.Log(secured, r, err.Error())
		}
		return r
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Expect", "")
	req.Header.Set("User-Agent", cl.GetUserAgent())
	val, err = cl.GetReferer()
	if err != nil {
		req.Header.Add("Referer", val)
	}
	resp, err2 := cl.client.Do(req)
	if err2 != nil {
		tpl := rtm.GetTemplate("httperror")
		r := R.NewResponse(tpl, newcmd, cfg)
		if cl.debugMode {
			cl.logger.Log(secured, r, err2.Error())
		}
		return r
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			tpl := rtm.GetTemplate("httperror")
			r := R.NewResponse(tpl, newcmd, cfg)
			if cl.debugMode {
				cl.logger.Log(secured, r, err.Error())
			}
			return r
		}
		r := R.NewResponse(string(response), newcmd, cfg)
		if cl.debugMode {
			cl.logger.Log(secured, r)
		}
		return r
	}
	tpl := rtm.GetTemplate("httperror")
	r := R.NewResponse(tpl, newcmd, cfg)
	if cl.debugMode {
		cl.logger.Log(secured, r)
	}
	return r
}

// RequestNextResponsePage method to request the next page of list entries for the current list query
// Useful for lists
func (cl *APIClient) RequestNextResponsePage(rr *R.Response) (*R.Response, error) {
	mycmd := map[string]interface{}{}
	for key, val := range rr.GetCommand() {
		mycmd[key] = val
	}
	if _, ok := mycmd["LAST"]; ok {
		return nil, errors.New("parameter LAST in use. Please remove it to avoid issues in requestNextPage")
	}
	first := 0
	if v, ok := mycmd["FIRST"]; ok {
		first, _ = fmt.Sscan("%s", v) //nolint:errcheck
	}
	total := rr.GetRecordsTotalCount()
	limit := rr.GetRecordsLimitation()
	first += limit
	if first < total {
		mycmd["FIRST"] = fmt.Sprintf("%d", first)
		mycmd["LIMIT"] = fmt.Sprintf("%d", limit)
		return cl.Request(mycmd), nil
	}
	return nil, errors.New("could not find further existing pages")
}

// RequestAllResponsePages method to request all pages/entries for the given query command
// Use this method with caution as it requests all list data until done.
func (cl *APIClient) RequestAllResponsePages(cmd map[string]string) []R.Response {
	var err error
	responses := []R.Response{}
	mycmd := map[string]interface{}{}
	mycmd["FIRST"] = "0"
	for k, v := range cmd {
		mycmd[k] = v
	}
	rr := cl.Request(mycmd)
	tmp := rr
	for {
		responses = append(responses, *tmp)
		tmp, err = cl.RequestNextResponsePage(tmp)
		if err != nil {
			break
		}
	}
	return responses
}

// UseDefaultConnectionSetup to activate default conneciton setup (the default anyways)
func (cl *APIClient) UseDefaultConnectionSetup() *APIClient {
	cl.SetURL(CNR_CONNECTION_URL_LIVE)
	return cl
}

// UseOTESystem method to set OT&E System for API communication
func (cl *APIClient) UseOTESystem() *APIClient {
	cl.SetURL(CNR_CONNECTION_URL_OTE)
	return cl
}

// UseLIVESystem method to set LIVE System for API communication
// Usage of LIVE System is active by default.
func (cl *APIClient) UseLIVESystem() *APIClient {
	cl.SetURL(CNR_CONNECTION_URL_LIVE)
	return cl
}

// flattenCommand method to translate all command parameter names to uppercase
func (cl *APIClient) flattenCommand(cmd map[string]interface{}) map[string]string {
	newcmd := map[string]string{}
	if len(cmd) == 0 {
		return newcmd
	}
	for key, val := range cmd {
		newKey := strings.ToUpper(key)
		if reflect.TypeOf(val).Kind() == reflect.Slice {
			v := val.([]string)
			for idx, str := range v {
				str = strings.ReplaceAll(str, "\r", "")
				str = strings.ReplaceAll(str, "\n", "")
				newcmd[newKey+strconv.Itoa(idx)] = str
			}
		} else {
			val := val.(string)
			val = strings.ReplaceAll(val, "\r", "")
			val = strings.ReplaceAll(val, "\n", "")
			newcmd[newKey] = val
		}
	}
	return newcmd
}

// autoIDNConvert method to translate all whitelisted parameter values to punycode, if necessary
func (cl *APIClient) autoIDNConvert(cmd map[string]string) map[string]string {
	if len(cmd) == 0 {
		return cmd
	}
	// don't convert for convertidn command to avoid endless loop
	pattern := regexp.MustCompile(`(?i)^CONVERTIDN$`)
	mm := pattern.MatchString(cmd["COMMAND"])
	if mm {
		return cmd
	}
	keys := []string{}
	pattern = regexp.MustCompile(`(?i)^(DOMAIN|NAMESERVER|DNSZONE)([0-9]*)$`)
	for key := range cmd {
		mm = pattern.MatchString(key)
		if mm {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return cmd
	}
	toconvert := []string{}
	idxs := []string{}
	pattern = regexp.MustCompile(`\r|\n`)
	idnpattern := regexp.MustCompile(`(?i)[^a-z0-9. -]+`)
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		val := pattern.ReplaceAllString(cmd[key], "")
		mm = idnpattern.MatchString(val)
		if mm {
			toconvert = append(toconvert, val)
			idxs = append(idxs, key)
		}
	}
	if len(toconvert) == 0 {
		return cmd
	}
	r := IDN.Convert(toconvert)

	for idx, pc := range r {
		cmd[idxs[idx]] = pc.PUNYCODE
	}
	return cmd
}
