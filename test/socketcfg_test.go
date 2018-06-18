package test

import (
	"apiconnector/client/socketcfg"
	"strings"
	"testing"
)

func TestSetCredentials(t *testing.T) {
	scfg := socketcfg.Socketcfg{}
	scfg.SetCredentials("test.user", "test.passw0rd", "")
	if strings.Compare(scfg.EncodeData(), "s_login=test.user&s_pw=test.passw0rd&") != 0 {
		t.Error("TestSetCredentials: Expected credentials couldn't be set.")
	}
}

func TestSetEntity(t *testing.T) {
	scfg := socketcfg.Socketcfg{}
	scfg.SetEntity("1234")
	if strings.Compare(scfg.EncodeData(), "s_entity=1234&") != 0 {
		t.Error("TestSetEntity: Expected entity couldn't be set.")
	}
	scfg.SetEntity("54cd")
	if strings.Compare(scfg.EncodeData(), "s_entity=54cd&") != 0 {
		t.Error("TestSetEntity: Expected entity couldn't be set.")
	}
}

func TestSetSession(t *testing.T) {
	scfg := socketcfg.Socketcfg{}
	scfg.SetCredentials("test.user", "test.passw0rd", "")
	scfg.SetSession("MYAPISESSIONID")
	if strings.Compare(scfg.EncodeData(), "s_session=MYAPISESSIONID&") != 0 {
		t.Error("TestSetSession: Expected session id couldn't be set.")
	}
}

func TestSetUser(t *testing.T) {
	scfg := socketcfg.Socketcfg{}
	scfg.SetCredentials("test.user", "test.passw0rd", "")
	scfg.SetUser("hexotestman.com")
	if strings.Compare(scfg.EncodeData(), "s_login=test.user&s_pw=test.passw0rd&s_user=hexotestman.com&") != 0 {
		t.Error("TestSetUser: Expected user couldn't be set.")
	}
}
