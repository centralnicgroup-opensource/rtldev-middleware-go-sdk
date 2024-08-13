package responsetranslator_test

import (
	"testing"

	"github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/response"
	rp "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responseparser"
	rt "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responsetranslator"
)

func TestTranslate(t *testing.T) {
	cmd := map[string]string{"COMMAND": "CheckDomainTransfer", "DOMAIN": "google.com", "AUTH": "blablabla"}

	// Test ACL error translation
	t.Run("ACLTranslation", func(t *testing.T) {
		raw := "[RESPONSE]\r\ncode = 530\r\ndescription = Authorization failed; Operation forbidden by ACL\r\nEOF\r\n"
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

	t.Run("CheckDomainTransferTranslation", func(t *testing.T) {
		testCases := map[string]string{
			"Domain status does not allow for operation":               "This Domain is locked. Initiating a Transfer is therefore impossible.",
			"Authorization failed [Invalid authorization information]": "The given Authorization Code is wrong. Initiating a Transfer is therefore impossible.",
		}

		for input, expected := range testCases {
			raw := "[RESPONSE]\r\ncode = 219\r\ndescription = " + input + "\r\nEOF\r\n"
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
			"[RESPONSE]\r\ncode = 505\r\ndescription = Invalid attribute value syntax; resource record [213123 A 1.2.4.5asdfa]\r\nEOF\r\n": "Invalid Syntax for DNSZone Resource Record: 213123 A 1.2.4.5asdfa",
			"[RESPONSE]\r\ncode = 505\r\ndescription = Syntax error in Parameter DOMAIN (my–domain.de)\r\nEOF\r\n":                         "The Domain Name my–domain.de is invalid.",
		}

		for input, expected := range testCases {
			newraw := rt.Translate(input, cmd)
			r := &response.Response{
				Raw:  newraw,
				Hash: rp.Parse(newraw),
			}
			h := r.GetHash()

			// Debug output
			// fmt.Printf("Input: %s\nExpected: %s\nGot: %s\n", input, expected, h["DESCRIPTION"])

			if h["DESCRIPTION"] != expected {
				t.Errorf("Expected: %s, got: %s", expected, h["DESCRIPTION"])
			}
		}
	})
}
