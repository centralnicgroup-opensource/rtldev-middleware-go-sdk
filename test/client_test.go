package test

import (
	"apiconnector/client"
	"strings"
	"testing"
)

func TestReadme1(t *testing.T) {
	allok := true
	cl := client.NewClient()
	cl.SetCredentials("test.user", "test.passw0rd", "")
	cl.UseOTESystem()
	r := cl.Login()
	if r.IsSuccess() {
		cmd := map[string]string{
			"COMMAND": "StatusAccount",
		}
		r = cl.Request(cmd)
		if r.IsSuccess() {
			r = cl.Logout()
			if !r.IsSuccess() {
				allok = false
			}
		} else {
			allok = false
		}
	} else {
		allok = false
	}
	if !allok {
		t.Error("TestReadme1: Something went wront.")
	}
}

func TestReadme2(t *testing.T) {
	allok := true
	cl := client.NewClient()
	cl.SetCredentials("test.user", "test.passw0rd", "")
	cl.UseOTESystem()
	cmd := map[string]string{
		"COMMAND": "StatusAccount",
	}
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		allok = false
	}
	if !allok {
		t.Error("TestReadme2: Something went wront.")
	}
}

func TestApiSession(t *testing.T) {
	cl := client.NewClient()
	cl.SetCredentials("test.user", "test.passw0rd", "")
	cl.UseOTESystem()
	r := cl.Login()
	if !r.IsSuccess() {
		t.Error("TestApiSession: Login failed.")
	}
	cmd := map[string]string{
		"COMMAND": "GetUserIndex",
	}
	r = cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestApiSession: Command Request failed.")
	}
	v, _ := r.GetColumnIndex("PARENTUSERINDEX", 0)
	if strings.Compare(v, "199") != 0 {
		t.Error("TestApiSession: Got wrong PARENTUSERINDEX value.")
	}
	v, _ = r.GetColumnIndex("USERINDEX", 0)
	if strings.Compare(v, "659") != 0 {
		t.Error("TestApiSession: Got wrong USERINDEX value.")
	}
	r = cl.Logout()
	if !r.IsSuccess() {
		t.Error("TestApiSession: Logout failed.")
	}
}

func TestSessionlessRequest(t *testing.T) {
	cl := client.NewClient()
	cl.SetCredentials("test.user", "test.passw0rd", "")
	cl.UseOTESystem()
	cmd := map[string]string{
		"COMMAND": "GetUserIndex",
	}
	r := cl.Request(cmd)
	if !r.IsSuccess() {
		t.Error("TestApiSession: Command Request failed.")
	}
	v, _ := r.GetColumnIndex("PARENTUSERINDEX", 0)
	if strings.Compare(v, "199") != 0 {
		t.Error("TestApiSession: Got wrong PARENTUSERINDEX value.")
	}
	v, _ = r.GetColumnIndex("USERINDEX", 0)
	if strings.Compare(v, "659") != 0 {
		t.Error("TestApiSession: Got wrong USERINDEX value.")
	}
}

func TestEncodeData(t *testing.T) {
	cmd := map[string]string{
		"COMMAND": "GetUserIndex",
	}
	cl := client.NewClient()
	cl.SetCredentials("test.user", "test.passw0rd", "")
	cl.UseOTESystem()
	d := cl.EncodeData(&cl.Socketcfg, cmd)
	if strings.Compare(d, "s_login=test.user&s_pw=test.passw0rd&s_entity=1234&s_command=COMMAND%3DGetUserIndex%0A") != 0 {
		t.Error("TestEncodeData: encoded data string wrong")
	}
}

func TestApiUrlGetterAndSetter(t *testing.T) {
	cl := client.NewClient()
	url := cl.Getapiurl()
	if strings.Compare(url, "https://coreapi.1api.net/api/call.cgi") != 0 {
		t.Error("TestApiUrlGetterAndSetter: wrong default apiurl value")
	}
	cl.Setapiurl("http://1api.de")
	url2 := cl.Getapiurl()
	if strings.Compare(url2, "http://1api.de") != 0 {
		t.Error("TestApiUrlGetterAndSetter: wrong apiurl value")
	}
	cl.Setapiurl(url)
	if strings.Compare(url, cl.Getapiurl()) != 0 {
		t.Error("TestApiUrlGetterAndSetter: wrong apiurl value")
	}
}
