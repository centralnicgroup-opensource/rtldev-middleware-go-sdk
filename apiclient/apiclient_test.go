package apiclient

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	R "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/response"
	"github.com/stretchr/testify/assert"
)

var cl = NewAPIClient()

func TestMain(m *testing.M) {

	rtm.AddTemplate(
		"login200",
		"[RESPONSE]\r\nproperty[expiration date][0] = 2024-09-19 10:52:51\r\nproperty[sessionid][0] = bb7a884b09b9a674fb4a22211758ce87\r\ndescription = Command completed successfully\r\ncode = 200\r\nqueuetime = 0.004\r\nruntime = 0.023\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"listP0",
		"[RESPONSE]\r\nproperty[total][0] = 4\r\nproperty[first][0] = 0\r\nproperty[domain][0] = cnic-ssl-test1.com\r\nproperty[domain][1] = cnic-ssl-test2.com\r\nproperty[count][0] = 2\r\nproperty[last][0] = 1\r\nproperty[limit][0] = 2\r\ndescription = Command completed successfully\r\ncode = 200\r\nqueuetime = 0\r\nruntime = 0.007\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"OK",
		rtm.GenerateTemplate("200", "Command completed successfully"),
	)
	os.Exit(m.Run())
}

func TestGetPOSTData1(t *testing.T) {
	validate := "s_command=AUTH%3Dgwrgwqg%25%26%5C44t3%2A%0ACOMMAND%3DModifyDomain"
	enc := cl.GetPOSTData(map[string]string{
		"COMMAND": "ModifyDomain",
		"AUTH":    "gwrgwqg%&\\44t3*",
	})
	if strings.Compare(enc, validate) != 0 {
		t.Error(fmt.Printf("TestGetPOSTData1: Expected encoding result '%s' not matching '%s'.", enc, validate))
	}
}

func TestGetPOSTDataSecured(t *testing.T) {
	testUser := url.QueryEscape(os.Getenv("CNR_TEST_USER"))
	validate := "s_login=" + testUser + "&s_pw=***&s_command=COMMAND%3DCheckAuthentication%0APASSWORD%3D%2A%2A%2A%0ASUBUSER%3D" + testUser
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	enc := cl.GetPOSTData(map[string]string{
		"COMMAND":  "CheckAuthentication",
		"SUBUSER":  os.Getenv("CNR_TEST_USER"),
		"PASSWORD": os.Getenv("CNR_TEST_PASSWORD"),
	}, true)
	if strings.Compare(enc, validate) != 0 {
		t.Error(fmt.Printf("TestGetPOSTDataSecured: Expected encoding result not matching\n\n%s\n%s.", enc, validate))
	}
	cl.SetCredentials("", "")
}

func TestEnableDebugMode(_ *testing.T) {
	cl.EnableDebugMode()
}

func TestDisableDebugMode(_ *testing.T) {
	cl.DisableDebugMode()
}

func TestRequestFlattenCommand(t *testing.T) {
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.UseOTESystem()
	r := cl.Request(map[string]interface{}{
		"COMMAND": "CheckDomains",
		"DOMAiN":  []string{"example.com", "example.net"},
	})
	if !r.IsSuccess() || r.GetCode() != 200 || r.GetDescription() != "Command completed successfully" {
		t.Error("TestRequestFlattenCommand: Expected response to succeed.")
	}
	cmd := r.GetCommand()
	val1, exists1 := cmd["DOMAIN0"]
	val2, exists2 := cmd["DOMAIN1"]
	_, exists3 := cmd["DOMAIN"]
	_, exists4 := cmd["DOMAiN"]
	if !exists1 || !exists2 || exists3 || exists4 {
		t.Error("TestRequestFlattenCommand: DOMAIN parameter flattening not working (keys).")
	}
	if val1 != "example.com" || val2 != "example.net" {
		t.Error("TestRequestFlattenCommand: DOMAIN parameter flattening not working (vals).")
	}
}

func TestAutoIDNConvertCommand(t *testing.T) {
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.UseOTESystem()
	r := cl.Request(map[string]interface{}{
		"COMMAND": "CheckDomains",
		"DOMAIN":  []string{"example.com", "dömäin.example", "example.net"},
	})
	if !r.IsSuccess() || r.GetCode() != 200 || r.GetDescription() != "Command completed successfully" {
		t.Error("TestRequestFlattenCommand: Expected response to succeed." + strconv.Itoa(r.GetCode()) + r.GetDescription())
	}
	cmd := r.GetCommand()

	val1, exists1 := cmd["DOMAIN0"]
	val2, exists2 := cmd["DOMAIN1"]
	val3, exists3 := cmd["DOMAIN2"]
	_, exists4 := cmd["DOMAIN"]
	_, exists5 := cmd["DOMAiN"]
	if !exists1 || !exists2 || !exists3 || exists4 || exists5 {
		t.Error("TestRequestFlattenCommand: DOMAIN parameter flattening not working (keys).")
	}
	if val1 != "example.com" || val2 != "xn--dmin-moa0i.example" || val3 != "example.net" {
		t.Error("TestRequestFlattenCommand: DOMAIN parameter flattening not working (vals).")
	}
	// reset to defaults for following tests
	cl.SetCredentials("", "")
	cl.UseLIVESystem()
}

func TestGetURL(t *testing.T) {
	url := cl.GetURL()
	if strings.Compare(url, CNR_CONNECTION_URL_LIVE) != 0 {
		t.Error("TestGetURL: Expected url not matching.")
	}
}

func TestGetUserAgent(t *testing.T) {
	uaexpected := "GO-SDK (" + runtime.GOOS + "; " + runtime.GOARCH + "; rv:" + cl.GetVersion() + ") go/" + runtime.Version()
	ua := cl.GetUserAgent()
	if strings.Compare(ua, uaexpected) != 0 {
		t.Error("TestGetUserAgent: Expected user-agent not matching.")
	}
}

func TestSetUserAgent(t *testing.T) {
	uaid := "WHMCS"
	uarv := "7.7.0"
	uaexpected := uaid + " (" + runtime.GOOS + "; " + runtime.GOARCH + "; rv:" + uarv + ") go-sdk/" + cl.GetVersion() + " go/" + runtime.Version()
	ua := cl.SetUserAgent(uaid, uarv).GetUserAgent()
	if strings.Compare(ua, uaexpected) != 0 {
		t.Error("TestGetUserAgent: Expected user-agent not matching.")
	}
}

func TestSetUserAgentModules(t *testing.T) {
	uaid := "WHMCS"
	uarv := "7.7.0"
	mods := []string{"reg/2.6.2", "ssl/7.2.2", "dc/8.2.2"}
	uaexpected := uaid + " (" + runtime.GOOS + "; " + runtime.GOARCH + "; rv:" + uarv + ") reg/2.6.2 ssl/7.2.2 dc/8.2.2 go-sdk/" + cl.GetVersion() + " go/" + runtime.Version()
	ua := cl.SetUserAgent(uaid, uarv, mods).GetUserAgent()
	if strings.Compare(ua, uaexpected) != 0 {
		t.Error("TestGetUserAgent: Expected user-agent not matching.")
	}
}

func TestSetURL(t *testing.T) {
	url := cl.SetURL(CNR_CONNECTION_URL_PROXY).GetURL()
	if strings.Compare(CNR_CONNECTION_URL_PROXY, url) != 0 {
		t.Error("TestSetURL: Expected url not matching.")
	}
	cl.SetURL(CNR_CONNECTION_URL_LIVE)
}

func TestSetCredentials1(t *testing.T) {
	cl.SetCredentials("myaccountid", "mypassword")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_login=myaccountid&s_pw=mypassword&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetCredentials1: Expected post data string not matching.")
	}
}

func TestSetCredentials2(t *testing.T) {
	cl.SetCredentials("", "")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetCredentials2: Expected post data string not matching.")
	}
}

func TestSetRoleCredentials1(t *testing.T) {
	cl.SetRoleCredentials("myaccountid", "myroleid", "mypassword")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_login=myaccountid%3Amyroleid&s_pw=mypassword&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetRoleCredentials1: Expected post data string not matching.")
	}
}

func TestSetRoleCredentials2(t *testing.T) {
	cl.SetRoleCredentials("", "", "")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetRoleCredentials2: Expected post data string not matching.")
	}
}

func TestSetProxy(t *testing.T) {
	cl.SetProxy("127.0.0.1")
	val, err := cl.GetProxy()
	if err != nil || val != "127.0.0.1" {
		t.Error("TestSetProxy: proxy not matching expected value")
	}
	cl.SetProxy("")
}

func TestSetReferer(t *testing.T) {
	cl.SetReferer("https://www.centralnicreseller.com/")
	val, err := cl.GetReferer()
	if err != nil || val != "https://www.centralnicreseller.com/" {
		t.Error("TestSetReferer: referer not matching expected value")
	}
	cl.SetReferer("")
}

func TestUseHighPerformanceConnectionSetup(t *testing.T) {
	cl.UseHighPerformanceConnectionSetup()
	val := cl.GetURL()
	if val != CNR_CONNECTION_URL_PROXY {
		t.Error("TestUseHighPerformanceConnectionSetup: couldn't activate high performance connection setup")
	}
}

func TestDefaultConnectionSetup(t *testing.T) {
	cl.UseDefaultConnectionSetup()
	val := cl.GetURL()
	if val != CNR_CONNECTION_URL_LIVE {
		t.Error("TestDefaultConnectionSetup: couldn't activate default connection setup")
	}
}

func TestAccountStatus(t *testing.T) {
	cl.UseOTESystem()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.EnableDebugMode()
	cmd := map[string]interface{}{}
	cmd["COMMAND"] = "StatusAccount"
	r := cl.Request(cmd)
	if r.GetDescription() == "Authorization failed" {
		t.Error("TestAccountStatus: Please make sure correct credentials are provided")
	}

	if !r.IsSuccess() {
		t.Error("TestAccountStatus: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestAccountStatus: Expected record not to be nil.")
	}
}

func TestLogin(t *testing.T) {
	cl.UseOTESystem()
	cl.EnableDebugMode()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	r := cl.Login()
	if !r.IsSuccess() {
		t.Error("TestLogin2: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLogin2: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("SESSIONID")
	if err != nil || d == "" {
		t.Error("TestLogin2: Expected session not to be empty.")
	}
}

func TestLogin3(t *testing.T) {
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), "WRONGPASSWORD")
	cl.EnableDebugMode()
	r := cl.Login()
	if !r.IsError() {
		t.Error("TestLogin3: Expected response to be an error case.")
	}
}

/**
 * Make sure session is Cleaned up if password is provided after session is saved
 */
func TestLogin4(t *testing.T) {
	sessionobj := map[string]interface{}{
		"socketcfg": map[string]string{
			"login":   "myaccount",
			"session": "abc",
		},
	}
	// Initialize the first APIClient instance
	cl.ReuseSession(sessionobj).SetCredentials("myaccountid", "password").EnableDebugMode()

	// Prepare the command map
	cmd := map[string]string{
		"COMMAND": "StatusAccount",
	}

	// Get the POST data from the second APIClient instance
	tmp := cl.GetPOSTData(cmd)

	// Validate the result
	if strings.Compare(tmp, "s_login=myaccountid&s_pw=password&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestLogin4: Expected post data string not matching." + tmp)
	}
}

// validate against mocked API response [login failed; http timeout] // need mocking
// validate against mocked API response [login succeeded; no session returned] // need mocking

func TestLogout1(t *testing.T) {
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.UseOTESystem()
	cl.EnableDebugMode()
	r := cl.Login()
	if !r.IsSuccess() {
		t.Error("TestLogout1: Expected response to be a success case.")
	}
	r = cl.Logout()
	if !r.IsSuccess() {
		t.Error("TestLogout1: Expected response to be a success case.")
	}
}

func TestSaveAndReuseSession(t *testing.T) {
	// Initialize the first APIClient instance
	sessionobj := make(map[string]interface{})
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.Login()
	cl.SaveSession(sessionobj)

	// Initialize the second APIClient instance
	cl2 := NewAPIClient()
	cl2.ReuseSession(sessionobj)

	// Prepare the command map
	cmd := map[string]string{
		"COMMAND": "StatusAccount",
	}

	// Get the POST data from the second APIClient instance
	tmp := cl2.GetPOSTData(cmd)

	// Validate the result
	if !strings.Contains(tmp, "s_sessionid") {
		t.Error("TestSaveReuseSession: Expected post data string to contain session ID.")
	}
}

// validate against mocked API response [200 < r.statusCode > 299] // need mocking
// validate against mocked API response [200 < r.statusCode > 299, no debug] // need mocking

func TestRequestNextResponsePage1(t *testing.T) { // nolint: gocyclo
	r := R.NewResponse(
		rtm.GetTemplate("listP0"),
		map[string]string{
			"COMMAND": "QueryDomainList",
			"LIMIT":   "2",
		},
	)
	cl.UseOTESystem()
	cl.DisableDebugMode()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	nr, err := cl.RequestNextResponsePage(r)
	if err != nil {
		t.Error(err)
		t.Error("TestRequestNextResponsePage1: Expected not to run into error.")
	}
	if !r.IsSuccess() {
		t.Error("TestRequestNextResponsePage1: Expected response (r) to be a success case.")
	}
	if !nr.IsSuccess() {
		t.Error("TestRequestNextResponsePage1: Expected response (nr) to be a success case.")
	}
	if r.GetRecordsLimitation() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected limitation (r) not matching.")
	}
	if nr.GetRecordsLimitation() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected limitation (nr) not matching.")
	}
	if r.GetRecordsCount() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected count (r) not matching.")
	}
	if nr.GetRecordsCount() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected count (nr) not matching.")
	}
	f, err := r.GetFirstRecordIndex()
	if err != nil || f != 0 {
		t.Error("TestRequestNextResponsePage1: Expected first (r) not matching.")
	}
	l, err := r.GetLastRecordIndex()
	if err != nil || l != 1 {
		t.Error("TestRequestNextResponsePage1: Expected last (r) not matching.")
	}
	f, err = nr.GetFirstRecordIndex()
	if err != nil || f != 2 {
		t.Error("TestRequestNextResponsePage1: Expected first (nr) not matching.")
	}
	l, err = nr.GetLastRecordIndex()
	if err != nil || l != 3 {
		t.Error("TestRequestNextResponsePage1: Expected last (nr) not matching.")
	}
}

func TestRequestNextResponsePage2(t *testing.T) {
	r := R.NewResponse(
		rtm.GetTemplate("listP0"),
		map[string]string{
			"COMMAND": "QueryDomainList",
			"LIMIT":   "2",
			"FIRST":   "0",
			"LAST":    "1",
		},
	)
	_, err := cl.RequestNextResponsePage(r)
	if err == nil {
		t.Error("TestRequestNextResponsePage2: Expected error to be returned as parameter LAST is in use.")
	}
}

func TestRequestNextResponsePage3(t *testing.T) { // nolint: gocyclo
	r := R.NewResponse(
		rtm.GetTemplate("listP0"),
		map[string]string{
			"COMMAND": "QueryDomainList",
			"LIMIT":   "2",
		},
	)
	nr, err := cl.RequestNextResponsePage(r)
	if err != nil {
		t.Error("TestRequestNextResponsePage3: Expected not to run into error.")
	}
	if !r.IsSuccess() {
		t.Error("TestRequestNextResponsePage3: Expected response (r) to be a success case.")
	}
	if !nr.IsSuccess() {
		t.Error("TestRequestNextResponsePage3: Expected response (nr) to be a success case.")
	}
	if r.GetRecordsLimitation() != 2 {
		t.Error("TestRequestNextResponsePage3: Expected limitation (r) not matching.")
	}
	if nr.GetRecordsLimitation() != 2 {
		t.Error("TestRequestNextResponsePage3: Expected limitation (nr) not matching.")
	}
	if r.GetRecordsCount() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected count (r) not matching.")
	}
	if nr.GetRecordsCount() != 2 {
		t.Error("TestRequestNextResponsePage1: Expected count (nr) not matching.")
	}
	f, err := r.GetFirstRecordIndex()
	if err != nil || f != 0 {
		t.Error("TestRequestNextResponsePage1: Expected first (r) not matching.")
	}
	l, err := r.GetLastRecordIndex()
	if err != nil || l != 1 {
		t.Error("TestRequestNextResponsePage1: Expected last (r) not matching.")
	}
	f, err = nr.GetFirstRecordIndex()
	if err != nil || f != 2 {
		t.Error("TestRequestNextResponsePage1: Expected first (nr) not matching.")
	}
	l, err = nr.GetLastRecordIndex()
	if err != nil || l != 3 {
		t.Error("TestRequestNextResponsePage1: Expected last (nr) not matching.")
	}
}

func TestRequestAllResponsePages(t *testing.T) {
	nr := cl.RequestAllResponsePages(map[string]string{
		"COMMAND": "QuerySSLCertList",
		"FIRST":   "0",
		"LIMIT":   "1000",
	})
	if len(nr) == 0 {
		t.Error("TestRequestAllResponsePages: Expected count of pages not matching.")
	}
}

func TestSetUserView(t *testing.T) {
	cl.SetUserView("julia")
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cmd := map[string]interface{}{}
	cmd["COMMAND"] = "StatusAccount"
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestResetUserView: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLogin2: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("REGISTRAR")
	if err != nil {
		t.Errorf("Failed to get data by key 'REGISTRAR': %v", err)
	}
	assert.Equal(t, "julia", d)
}

func TestResetUserView(t *testing.T) {
	cl.SetUserView("julia")
	cl.ResetUserView()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cmd := map[string]interface{}{}
	cmd["COMMAND"] = "StatusAccount"
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestResetUserView: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLogin2: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("REGISTRAR")
	if err != nil {
		t.Errorf("Failed to get data by key 'REGISTRAR': %v", err)
	}
	assert.NotEqual(t, "julia", d)
}
