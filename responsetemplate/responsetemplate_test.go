package responsetemplate

import (
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	tpl := NewResponseTemplate("")
	if tpl.GetCode() != 423 {
		t.Error("TestConstructor: Expected response code not matching.")
	}
	if strings.Compare(tpl.GetDescription(), "Empty API response") != 0 {
		t.Error("TestConstructor: Expected response description not matching.")
	}
}

func TestGetHash(t *testing.T) {
	h := NewResponseTemplate("").GetHash()
	if v, ok := h["CODE"]; !ok || strings.Compare(v.(string), "423") != 0 {
		t.Error("TestGetHash: Expected response code not matching.")
	}
	if v, ok := h["DESCRIPTION"]; !ok || strings.Compare(v.(string), "Empty API response") != 0 {
		t.Error("TestGetHash: Expected response description not matching.")
	}
}

func TestGetQueuetime1(t *testing.T) {
	tpl := NewResponseTemplate("")
	if tpl.GetQueuetime() != 0 {
		t.Error("TestGetQueuetime1: Expected queuetime not matching")
	}
}

func TestGetQueuetime2(t *testing.T) {
	tpl := NewResponseTemplate("[RESPONSE]\r\ncode=423\r\ndescription=Empty API response\r\nqueuetime=0\r\nEOF\r\n")
	if tpl.GetQueuetime() != 0 {
		t.Error("TestGetQueuetime2: Expected queuetime not matching")
	}
}

func TestGetRuntime1(t *testing.T) {
	tpl := NewResponseTemplate("")
	if tpl.GetRuntime() != 0 {
		t.Error("TestGetRuntime1: Expected runtime not matching")
	}
}

func TestGetRuntime2(t *testing.T) {
	tpl := NewResponseTemplate("[RESPONSE]\r\ncode=423\r\ndescription=Empty API response\r\nruntime=0.12\r\nEOF\r\n")
	if tpl.GetRuntime() != 0.12 {
		t.Error("TestGetRuntime2: Expected runtime not matching")
	}
}
