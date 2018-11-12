package responsetemplatemanager

import (
	"fmt"
	"os"
	"strings"
	"testing"

	RT "github.com/hexonet/go-sdk/responsetemplate"
)

var rtm = GetInstance()

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

func TestGetTemplate(t *testing.T) {
	tpl := rtm.GetTemplate("IwontExist")
	if tpl.GetCode() != 500 {
		t.Error("Expected response code not matching")
	}
	if strings.Compare(tpl.GetDescription(), "Response Template not found") != 0 {
		t.Error("TestGetTemplate: Expected response description not matching")
	}
}

func TestGetTemplates(t *testing.T) {
	defaultones := []string{"404", "500", "error", "httperror", "empty", "unauthorized", "expired"}
	tpls := rtm.GetTemplates()
	for _, k := range defaultones {
		if _, ok := tpls[k]; !ok {
			t.Error(fmt.Sprintf("TestGetTemplates: Expected default template '%s' to exist.", k))
		}
	}
}

func TestIsTemplateMatchHash(t *testing.T) {
	tpl := RT.NewResponseTemplate("")
	h := tpl.GetHash()
	if !rtm.IsTemplateMatchHash(h, "empty") {
		t.Error("TestIsTemplateMatchHash: Expected hash response to match 'empty' response template.")
	}
}

func TestIsTemplateMatchPlain(t *testing.T) {
	tpl := RT.NewResponseTemplate("")
	plain := tpl.GetPlain()
	if !rtm.IsTemplateMatchPlain(plain, "empty") {
		t.Error("TestIsTemplateMatchPlain: Expected plain response to match 'empty' response template.")
	}
}
