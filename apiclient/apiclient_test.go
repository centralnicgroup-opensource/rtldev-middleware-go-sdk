package apiclient

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	R "github.com/hexonet/go-sdk/response"
)

var cl = NewAPIClient()

func TestMain(m *testing.M) {

	rtm.AddTemplate(
		"login200",
		"[RESPONSE]\r\nPROPERTY[SESSION][0]=h8JLZZHdF2WgWWXlwbKWzEG3XrzoW4yshhvtqyg0LCYiX55QnhgYX9cB0W4mlpbx\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.169\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"listP0",
		"[RESPONSE]\r\nPROPERTY[TOTAL][0]=2701\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[DOMAIN][0]=0-60motorcycletimes.com\r\nPROPERTY[DOMAIN][1]=0-be-s01-0.com\r\nPROPERTY[COUNT][0]=2\r\nPROPERTY[LAST][0]=1\r\nPROPERTY[LIMIT][0]=2\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.023\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"OK",
		rtm.GenerateTemplate("200", "Command completed successfully"),
	)
	os.Exit(m.Run())
}

func TestGetPOSTData1(t *testing.T) {
	validate := "s_entity=54cd&s_command=AUTH%3Dgwrgwqg%25%26%5C44t3%2A%0ACOMMAND%3DModifyDomain"
	enc := cl.GetPOSTData(map[string]string{
		"COMMAND": "ModifyDomain",
		"AUTH":    "gwrgwqg%&\\44t3*",
	})
	if strings.Compare(enc, validate) != 0 {
		t.Error(fmt.Printf("TestGetPOSTData1: Expected encoding result '%s' not matching '%s'.", enc, validate))
	}
}

func TestEnableDebugMode(t *testing.T) {
	cl.EnableDebugMode()
}

func TestDisableDebugMode(t *testing.T) {
	cl.DisableDebugMode()
}

func TestRequestFlattenCommand(t *testing.T) {
	cl.SetCredentials("test.user", "test.passw0rd")
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
	cl.SetCredentials("test.user", "test.passw0rd")
	cl.UseOTESystem()
	r := cl.Request(map[string]interface{}{
		"COMMAND": "CheckDomains",
		"DOMAiN":  []string{"example.com", "dömäin.example", "example.net"},
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

func TestGetSession1(t *testing.T) {
	cl.Logout()
	session, err := cl.GetSession()
	if err == nil || session != "" {
		t.Error("TestGetSession1: Expected no session, but found one.")
	}
}

func TestGetSesssion2(t *testing.T) {
	sessid := "testSessionID12345678"
	cl.SetSession(sessid)
	session, err := cl.GetSession()
	if err != nil {
		t.Error("TestGetSession2: Expected not to run into error.")
	}
	if strings.Compare(session, sessid) != 0 {
		t.Error("TestGetSession2: Expected session id not matching.")
	}
	cl.SetSession("")
}

func TestGetURL(t *testing.T) {
	url := cl.GetURL()
	if strings.Compare(url, ISPAPI_CONNECTION_URL) != 0 {
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

func TestSetURL(t *testing.T) {
	url := cl.SetURL(ISPAPI_CONNECTION_URL_PROXY).GetURL()
	if strings.Compare(ISPAPI_CONNECTION_URL_PROXY, url) != 0 {
		t.Error("TestSetURL: Expected url not matching.")
	}
	cl.SetURL(ISPAPI_CONNECTION_URL)
}

func TestSetOTP1(t *testing.T) {
	cl.SetOTP("12345678")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_otp=12345678&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetOTP1: Expected post data string not matching.")
	}
}

func TestSetOTP2(t *testing.T) {
	cl.SetOTP("")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetOTP2: Expected post data string not matching.")
	}
}

func TestSetSession1(t *testing.T) {
	cl.SetSession("12345678")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_session=12345678&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSession1: Expected post data string not matching.")
	}
}

func TestSetSession2(t *testing.T) {
	cl.SetRoleCredentials("myaccountid", "myrole", "mypassword")
	cl.SetOTP("12345678")
	cl.SetSession("12345678")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_session=12345678&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSession2: Expected post data string not matching.")
	}
}

func TestSetSession3(t *testing.T) {
	cl.SetSession("")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSession3: Expected post data string not matching.")
	}
}

func TestSaveReuseSession(t *testing.T) {
	sessionobj := map[string]interface{}{}
	cl.SetSession("12345678")
	cl.SaveSession(sessionobj)
	cl2 := NewAPIClient()
	cl2.ReuseSession(sessionobj)
	tmp := cl2.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_session=12345678&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSaveReuseSession: Expected post data string not matching.")
	}
	cl.SetSession("")
}

func TestSetRemoteIPAddress1(t *testing.T) {
	cl.SetRemoteIPAddress("10.10.10.10")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_remoteaddr=10.10.10.10&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetRemoteIPAddress1: Expected post data string not matching.")
	}
}

func TestSetRemoteIPAddress2(t *testing.T) {
	cl.SetRemoteIPAddress("")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetRemoteIPAddress2: Expected post data string not matching.")
	}
}

func TestSetCredentials1(t *testing.T) {
	cl.SetCredentials("myaccountid", "mypassword")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_login=myaccountid&s_pw=mypassword&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetCredentials1: Expected post data string not matching.")
	}
}

func TestSetCredentials2(t *testing.T) {
	cl.SetCredentials("", "")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetCredentials2: Expected post data string not matching.")
	}
}

func TestSetRoleCredentials1(t *testing.T) {
	cl.SetRoleCredentials("myaccountid", "myroleid", "mypassword")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_login=myaccountid%21myroleid&s_pw=mypassword&s_command=COMMAND%3DStatusAccount") != 0 {
		t.Error("TestSetRoleCredentials1: Expected post data string not matching.")
	}
}

func TestSetRoleCredentials2(t *testing.T) {
	cl.SetRoleCredentials("", "", "")
	tmp := cl.GetPOSTData(map[string]string{
		"COMMAND": "StatusAccount",
	})
	if strings.Compare(tmp, "s_entity=54cd&s_command=COMMAND%3DStatusAccount") != 0 {
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
	cl.SetReferer("https://www.hexonet.net/")
	val, err := cl.GetReferer()
	if err != nil || val != "https://www.hexonet.net/" {
		t.Error("TestSetReferer: referer not matching expected value")
	}
	cl.SetReferer("")
}

func TestUseHighPerformanceConnectionSetup(t *testing.T) {
	cl.UseHighPerformanceConnectionSetup()
	val := cl.GetURL()
	if val != ISPAPI_CONNECTION_URL_PROXY {
		t.Error("TestUseHighPerformanceConnectionSetup: couldn't activate high performance connection setup")
	}
}

func TestDefaultConnectionSetup(t *testing.T) {
	cl.UseDefaultConnectionSetup()
	val := cl.GetURL()
	if val != ISPAPI_CONNECTION_URL {
		t.Error("TestDefaultConnectionSetup: couldn't activate default connection setup")
	}
}

func TestLogin1(t *testing.T) {
	cl.UseOTESystem()
	cl.SetCredentials("test.user", "test.passw0rd")
	cl.SetRemoteIPAddress("1.2.3.4")
	cl.EnableDebugMode()
	r := cl.Login()
	if !r.IsSuccess() {
		t.Error("TestLogin1: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLogin1: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("SESSION")
	if err != nil || d == "" {
		t.Error("TestLogin1: Expected session not to be empty.")
	}
}

func TestLogin2(t *testing.T) {
	cl.UseOTESystem()
	cl.SetRoleCredentials("test.user", "testrole", "test.passw0rd")
	cl.SetRemoteIPAddress("1.2.3.4")
	r := cl.Login()
	if !r.IsSuccess() {
		t.Error("TestLogin2: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLogin2: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("SESSION")
	if err != nil || d == "" {
		t.Error("TestLogin2: Expected session not to be empty.")
	}
}

func TestLogin3(t *testing.T) {
	cl.SetCredentials("test.user", "WRONGPASSWORD")
	cl.SetRemoteIPAddress("1.2.3.4")
	r := cl.Login()
	if !r.IsError() {
		t.Error("TestLogin3: Expected response to be an error case.")
	}
}

// validate against mocked API response [login failed; http timeout] // need mocking
// validate against mocked API response [login succeeded; no session returned] // need mocking

func TestLoginExtended(t *testing.T) {
	cl.UseOTESystem()
	cl.SetCredentials("test.user", "test.passw0rd")
	cl.SetRemoteIPAddress("1.2.3.4")
	r := cl.LoginExtended(map[string]string{
		"TIMEOUT": "60",
	})
	if !r.IsSuccess() {
		t.Error("TestLoginExtended: Expected response to be a success case.")
	}
	rec := r.GetRecord(0)
	if rec == nil {
		t.Error("TestLoginExtended: Expected record not to be nil.")
	}
	d, err := rec.GetDataByKey("SESSION")
	if err != nil && d == "" {
		t.Error("TestLoginExtended: Expected session not to be empty.")
	}
}

func TestLogout1(t *testing.T) {
	r := cl.Logout()
	if !r.IsSuccess() {
		t.Error("TestLogout1: Expected response to be a success case.")
	}
}

func TestLogout2(t *testing.T) {
	tpl := R.NewResponse(
		rtm.GetTemplate("login200").GetPlain(),
		map[string]string{
			"COMMAND": "StartSession",
		},
	)
	rec := tpl.GetRecord(0)
	sessid, err := rec.GetDataByKey("SESSION")
	if err != nil {
		t.Error("TestLogout2: Expected not run into error.")
	}
	cl.EnableDebugMode()
	cl.SetSession(sessid)
	r := cl.Logout()
	if !r.IsError() {
		t.Error("TestLogout2: Expected response to be an error case.")
	}
}

// validate against mocked API response [200 < r.statusCode > 299] // need mocking
// validate against mocked API response [200 < r.statusCode > 299, no debug] // need mocking

func TestRequestNextResponsePage1(t *testing.T) {
	r := R.NewResponse(
		rtm.GetTemplate("listP0").GetPlain(),
		map[string]string{
			"COMMAND": "QueryDomainList",
			"LIMIT":   "2",
			"FIRST":   "0",
		},
	)
	cl.UseOTESystem()
	cl.SetRoleCredentials("test.user", "testrole", "test.passw0rd")
	cl.SetRemoteIPAddress("1.2.3.4")
	nr := cl.Login()
	if !nr.IsSuccess() {
		t.Error("TestRequestNextResponsePage1: Expected login response to be a success case.")
	}
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
		rtm.GetTemplate("listP0").GetPlain(),
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

func TestRequestNextResponsePage3(t *testing.T) {
	cl.DisableDebugMode()
	r := R.NewResponse(
		rtm.GetTemplate("listP0").GetPlain(),
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
	cl.SetUserView("hexotestman.com")
	cmd := map[string]interface{}{}
	cmd["COMMAND"] = "GetUserIndex"
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestSetUserView: Expected response to be a success case.")
	}
}

func TestResetUserView(t *testing.T) {
	cl.ResetUserView()
	cmd := map[string]interface{}{}
	cmd["COMMAND"] = "GetUserIndex"
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestResetUserView: Expected response to be a success case.")
	}
}
