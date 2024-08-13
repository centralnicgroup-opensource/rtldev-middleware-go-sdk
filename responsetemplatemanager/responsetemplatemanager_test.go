package responsetemplatemanager

import (
	"os"
	"strings"
	"testing"

	RP "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responseparser"
)

var rtm = GetInstance()

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

func TestGetTemplate(t *testing.T) {
	tpl := RP.Parse(rtm.GetTemplate("IwontExist"))
	if tpl["CODE"].(string) != "500" {
		t.Error("Expected response code not matching")
	}
	if strings.Compare(tpl["DESCRIPTION"].(string), "Response Template not found") != 0 {
		t.Error("TestGetTemplate: Expected response description not matching")
	}
}

func TestGetTemplates(t *testing.T) {
	defaultones := []string{"404", "500", "error", "httperror", "empty", "unauthorized", "expired"}
	tpls := rtm.GetTemplates()
	for _, k := range defaultones {
		if _, ok := tpls[k]; !ok {
			t.Errorf("TestGetTemplates: Expected default template '%s' to exist.", k)
		}
	}
}

func TestIsTemplateMatchHash(t *testing.T) {
	tpl := rtm.GetTemplate("empty")
	h := RP.Parse(tpl)
	if !rtm.IsTemplateMatchHash(h, "empty") {
		t.Error("TestIsTemplateMatchHash: Expected hash response to match 'empty' response template.")
	}
}

func TestIsTemplateMatchPlain(t *testing.T) {
	plain := rtm.GetTemplate("empty")
	if !rtm.IsTemplateMatchPlain(plain, "empty") {
		t.Error("TestIsTemplateMatchPlain: Expected plain response to match 'empty' response template.")
	}
}
