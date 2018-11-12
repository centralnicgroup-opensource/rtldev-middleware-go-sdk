package responseparser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	plain := "[RESPONSE]\r\nCODE=421\r\nDESCRIPTION=\r\nEOF\r\n"
	plain = strings.Replace(plain, "\r\nDESCRIPTION=", "", 1)
	fmt.Println(plain)
	parsed := Parse(plain)
	if v, ok := parsed["DESCRIPTION"]; ok {
		if v != "" {
			t.Error("TestParse: Expected description to be empty")
		}
	} else {
		t.Error("TestParse: Expected description to exist")
	}
}

func TestSerialize1(t *testing.T) {
	r := Parse("[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nEOF\r\n")
	r["PROPERTY"] = map[string][]string{
		"DOMAIN": {"mydomain1.com", "mydomain2.com", "mydomain3.com"},
		"RATING": {"1", "2", "3"},
		"SUM":    {"3"},
	}
	serialized := Serialize(r)
	expected := "[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nPROPERTY[DOMAIN][0]=mydomain1.com\r\nPROPERTY[DOMAIN][1]=mydomain2.com\r\nPROPERTY[DOMAIN][2]=mydomain3.com\r\nPROPERTY[RATING][0]=1\r\nPROPERTY[RATING][1]=2\r\nPROPERTY[RATING][2]=3\r\nPROPERTY[SUM][0]=3\r\nEOF\r\n"
	if strings.Compare(serialized, expected) != 0 {
		t.Error("TestSerialize1: Expected string not matching serialized format.")
	}
}

func TestSerialize2(t *testing.T) {
	expected := "[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nEOF\r\n"
	serialized := Serialize(Parse(expected))
	if strings.Compare(serialized, expected) != 0 {
		t.Error("TestSerialize2: Expected string not matching serialized format.")
	}
}

func TestSerialize3(t *testing.T) {
	// this case shouldn't happen, otherwise we have an API-side issue
	h := Parse("[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nEOF\r\n")
	delete(h, "CODE")
	delete(h, "DESCRIPTION")
	serialized := Serialize(h)
	expected := "[RESPONSE]\r\nEOF\r\n"
	if strings.Compare(serialized, expected) != 0 {
		t.Error("ETestSerialize3: xpected string not matching serialized format.")
	}
}

func TestSerialize4(t *testing.T) {
	h := Parse("[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nEOF\r\n")
	h["QUEUETIME"] = "0"
	h["RUNTIME"] = "0.12"
	serialized := Serialize(h)
	expected := "[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	if strings.Compare(serialized, expected) != 0 {
		t.Error("TestSerialize4: Expected string not matching serialized format.")
	}
}
