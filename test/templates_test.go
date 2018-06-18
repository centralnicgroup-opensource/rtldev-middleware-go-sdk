package test

import (
	"apiconnector/response/hashresponse"
	"strings"
	"testing"
)

func TestTemplateGetter(t *testing.T) {
	r := hashresponse.NewHashResponse(hashresponse.NewTemplates().Get("expired"))
	if r.Code() != 530 || strings.Compare(r.Description(), "SESSION NOT FOUND") != 0 {
		t.Error("TestTemplateGetter: API return code or description doesn't match expectation.")
	}
}

func TestTemplateGetterParsed(t *testing.T) {
	r := hashresponse.NewTemplates().GetParsed("expired")
	if strings.Compare(r["CODE"].(string), "530") != 0 || strings.Compare(r["DESCRIPTION"].(string), "SESSION NOT FOUND") != 0 {
		t.Error("TestTemplateGetterParsed: API return code or description doesn't match expectation.")
	}
}

func TestTemplateSetter(t *testing.T) {
	tplmgr := hashresponse.NewTemplates()
	tplmgr.Set("unauthorized", "[RESPONSE]\r\ncode=530\r\ndescription=Unauthorized\r\nTRANSLATIONKEY=FAPI.530.UNAUTHORIZED\r\nEOF\r\n")
	r := hashresponse.NewHashResponse(tplmgr.Get("unauthorized"))
	if r.Code() != 530 || strings.Compare(r.Description(), "Unauthorized") != 0 {
		t.Error("TestTemplateSetter: API return code or description doesn't match expectation.")
	}

	tplmgr.SetParsed("unauthorized2", tplmgr.GetParsed("unauthorized"))
	r = hashresponse.NewHashResponse(tplmgr.Get("unauthorized2"))
	if r.Code() != 530 || strings.Compare(r.Description(), "Unauthorized") != 0 {
		t.Error("TestTemplateSetter: API return code or description doesn't match expectation.")
	}
}

func TestTemplateContainer(t *testing.T) {
	tpls := hashresponse.NewTemplates().GetAll()
	i := len(tpls)
	if i != 4 {
		t.Error("TestTemplateContainer: Wrong amount of default templates.")
	}
	for key := range tpls {
		switch key {
		case "empty", "error", "expired", "commonerror":
			i--
		}
	}
	if i != 0 {
		t.Error("TestTemplateContainer: Check for default template keys failed.")
	}
}

func TestTemplateMatch(t *testing.T) {
	tplmgr := hashresponse.NewTemplates()
	r1 := tplmgr.GetParsed("unauthorized")
	if !tplmgr.MatchParsed(r1, "unauthorized") {
		t.Error("TestTemplateMatch: parsed response did not match template \"unauthorized\".")
	}

	r2 := tplmgr.Get("unauthorized")
	if !tplmgr.Match(r2, "unauthorized") {
		t.Error("TestTemplateMatch: response did not match template \"unauthorized\".")
	}
}
