package responsetranslator_test

import (
	"testing"

	"github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v4/response"
	rp "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v4/responseparser"
	rt "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v4/responsetranslator"
)

func TestTranslate(t *testing.T) {
	cmd := map[string]string{"COMMAND": "CheckDomainTransfer", "DOMAIN": "my–domain.com", "AUTH": "blablabla"}

	// Test ACL error translation
	t.Run("ACLTranslation", func(t *testing.T) {
		raw := "[RESPONSE]\r\ncode=530\r\ndescription=Authorization failed; Operation forbidden by ACL\r\nEOF\r\n"
		expected := "Authorization failed; Used Command `CheckDomainTransfer` not white-listed by your Access Control List"
		newraw := rt.Translate(raw, cmd)
		r := &response.Response{
			Raw:  newraw,
			Hash: rp.Parse(newraw),
		}
		h := r.GetHash()
		if h["DESCRIPTION"] != expected {
			t.Errorf("Expected: %s, got: %s", expected, h["DESCRIPTION"])
		}
	})

	// Test CheckDomainTransfer translations
	t.Run("CheckDomainTransferTranslation", func(t *testing.T) {
		testCases := map[string]string{
			// Add more test cases as needed
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (clientTransferProhibited)": "This Domain is locked. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (requested)":                "Registration of this Domain Name has not yet completed. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (requestedcreate)":          "Registration of this Domain Name has not yet completed. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (requesteddelete)":          "Deletion of this Domain Name has been requested. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (pendingdelete)":            "Deletion of this Domain Name is pending. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY WRONG AUTH":                        "The given Authorization Code is wrong. Initiating a Transfer is therefore impossible.",
			"Request is not available; DOMAIN TRANSFER IS PROHIBITED BY AGE OF THE DOMAIN":                 "This Domain Name is within 60 days of initial registration. Initiating a Transfer is therefore impossible.",
		}

		for input, expected := range testCases {
			raw := "[RESPONSE]\r\ncode=219\r\ndescription=" + input + "\r\nEOF\r\n"
			newraw := rt.Translate(raw, cmd)
			r := &response.Response{
				Raw:  newraw,
				Hash: rp.Parse(newraw),
			}
			h := r.GetHash()
			if h["DESCRIPTION"] != expected {
				t.Errorf("Expected: %s, got: %s", expected, h["DESCRIPTION"])
			}
		}
	})

	// Test translate function with various scenarios
	t.Run("Translate", func(t *testing.T) {
		testCases := map[string]string{
			// Add more test cases as needed
			"[RESPONSE]\r\ncode=219\r\ndescription=Request is not available; DOMAIN TRANSFER IS PROHIBITED BY STATUS (clientTransferProhibited)\r\nEOF\r\n": "This Domain is locked. Initiating a Transfer is therefore impossible.",
			"[RESPONSE]\r\ncode=505\r\ndescription=Invalid attribute value syntax; resource record [213123 A 1.2.4.5asdfa]\r\nEOF\r\n":                      "Invalid Syntax for DNSZone Resource Record: 213123 A 1.2.4.5asdfa",
			"[RESPONSE]\r\ncode=505\r\ndescription=Syntax error in Parameter DOMAIN (my–domain.de)\r\nEOF\r\n":                                              "The Domain Name my–domain.de is invalid.",
		}

		for input, expected := range testCases {
			newraw := rt.Translate(input, cmd)
			r := &response.Response{
				Raw:  newraw,
				Hash: rp.Parse(newraw),
			}
			h := r.GetHash()
			if h["DESCRIPTION"] != expected {
				t.Errorf("Expected: %s, got: %s", expected, h["DESCRIPTION"])
			}
		}
	})
}
