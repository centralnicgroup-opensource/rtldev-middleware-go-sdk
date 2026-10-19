package idntranslator_test

import (
	"testing"

	"github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/idntranslator"

	"github.com/stretchr/testify/assert"
)

type Row struct {
	IDN      string
	PUNYCODE string
}

func TestConvert(t *testing.T) {
	tests := []struct {
		domain   string
		expected []Row
	}{
		{
			domain: "",
			expected: []Row{
				{IDN: "", PUNYCODE: ""},
			},
		},
		{
			domain: "münchen.de",
			expected: []Row{
				{IDN: "münchen.de", PUNYCODE: "xn--mnchen-3ya.de"},
			},
		},
		{
			domain: "faß.de",
			expected: []Row{
				{IDN: "faß.de", PUNYCODE: "xn--fa-hia.de"},
			},
		},
		{
			domain: "faß.com",
			expected: []Row{
				{IDN: "faß.com", PUNYCODE: "fass.com"},
			},
		},
		{
			domain: "日本.co.jp",
			expected: []Row{
				{IDN: "日本.co.jp", PUNYCODE: "xn--wgv71a.co.jp"},
			},
		},
		{
			domain: "xn--wg8h.ws",
			expected: []Row{
				{IDN: "🌐.ws", PUNYCODE: "xn--wg8h.ws"},
			},
		},
		{
			domain: "ＡＢＣ・日本.co.jp",
			expected: []Row{
				{IDN: "abc・日本.co.jp", PUNYCODE: "xn--abc-rs4b422ycvb.co.jp"},
			},
		},
	}

	for _, test := range tests {
		result := idntranslator.Convert(test.domain) // Passing nil for options, as they're not used in this test
		assert.Equal(t, len(test.expected), len(result))
		for i := range test.expected {
			assert.Equal(t, test.expected[i].IDN, result[i].IDN)
			assert.Equal(t, test.expected[i].PUNYCODE, result[i].PUNYCODE)
		}
	}
}

func TestConvertBulk(t *testing.T) {
	// Define an array of domain names to test
	domains := []string{
		"münchen.de",
		"xn--mnchen-3ya.de",
		"🌐.ws",
		"xn--wg8h.ws",
		"😊.com",
		"xn--o28h.com",
		"🎉.net",
		"xn--dk8h.net",
	}

	// Call the Convert function
	convertedDomains := idntranslator.Convert(domains) // Passing nil for options, as they're not used in this test

	// Check if the converted domains have the correct values
	expected := []Row{
		{IDN: "münchen.de", PUNYCODE: "xn--mnchen-3ya.de"},
		{IDN: "münchen.de", PUNYCODE: "xn--mnchen-3ya.de"},
		{IDN: "🌐.ws", PUNYCODE: "xn--wg8h.ws"},
		{IDN: "🌐.ws", PUNYCODE: "xn--wg8h.ws"},
		{IDN: "😊.com", PUNYCODE: "xn--o28h.com"},
		{IDN: "😊.com", PUNYCODE: "xn--o28h.com"},
		{IDN: "🎉.net", PUNYCODE: "xn--dk8h.net"},
		{IDN: "🎉.net", PUNYCODE: "xn--dk8h.net"},
	}

	for i, row := range expected {
		assert.Equal(t, row.IDN, convertedDomains[i].IDN)
		assert.Equal(t, row.PUNYCODE, convertedDomains[i].PUNYCODE)
	}
}

func TestToASCIIWithTransitional(t *testing.T) {
	expected := map[string]string{
		"fass.de":          "fass.de",
		"₹.com":            "xn--yzg.com",
		"𑀓.com":            "xn--n00d.com",
		"öbb.at":           "xn--bb-eka.at",
		"ÖBB.at":           "xn--bb-eka.at",
		"ȡog.de":           "xn--og-09a.de",
		"☕.de":             "xn--53h.de",
		"I♥NY.de":          "xn--iny-zx5a.de",
		"ＡＢＣ・日本.co.jp":     "xn--abc-rs4b422ycvb.co.jp",
		"日本｡co｡jp":         "xn--wgv71a.co.jp",
		"日本｡co．jp":         "xn--wgv71a.co.jp",
		"x\u0327\u0301.de": "xn--x-xbb7i.de",
		"x\u0301\u0327.de": "xn--x-xbb7i.de",
		"عربي.de":          "xn--ngbrx4e.de",
		"نامهای.de":        "xn--mgba3gch31f.de",
		"fäß.de":           "xn--fss-qla.de",
		"faß.de":           "fass.de",
		"xn--fa-hia.de":    "fass.de",
		"σόλος.gr":         "xn--wxaikc6b.gr",
		"Σόλος.gr":         "xn--wxaikc6b.gr",
		"نامه\u200Cای.de":  "xn--mgba3gch31f.de",
		"☃-⌘.com":          "xn----dqo34k.com",
	}

	transitionalProcessing := true
	for domain, punycode := range expected {
		result := idntranslator.ToASCII(domain, transitionalProcessing)
		assert.Equal(t, punycode, result)
	}
}

func TestToASCIIWithoutTransitional(t *testing.T) {
	expected := map[string]string{
		"σόλος.gr":           "xn--wxaijb9b.gr",
		"Σόλος.gr":           "xn--wxaijb9b.gr",
		"ΣΌΛΟΣ.grﻋﺮﺑﻲ.de":    "xn--wxaikc6b.xn--gr-gtd9a1b0g.de",
		"fäß.de":             "xn--f-qfao.de",
		"faß.de":             "xn--fa-hia.de",
		"xn--bb-eka.at":      "xn--bb-eka.at",
		"XN--BB-EKA.AT":      "xn--bb-eka.at",
		"fass.de":            "fass.de",
		"not=std3":           "not=std3",
		"öbb.at":             "xn--bb-eka.at",
		"₹.com":              "xn--yzg.com",
		"𑀓.com":              "xn--n00d.com",
		"ÖBB.at":             "xn--bb-eka.at",
		"ȡog.de":             "xn--og-09a.de",
		"☕.de":               "xn--53h.de",
		"I♥NY.de":            "xn--iny-zx5a.de",
		"ＡＢＣ・日本.co.jp":       "xn--abc-rs4b422ycvb.co.jp",
		"日本｡co｡jp":           "xn--wgv71a.co.jp",
		"日本｡co．jp":           "xn--wgv71a.co.jp",
		"x\u0327\u0301.de":   "xn--x-xbb7i.de",
		"x\u0301\u0327.de":   "xn--x-xbb7i.de",
		"عربي.de":            "xn--ngbrx4e.de",
		"نامهای.de":          "xn--mgba3gch31f.de",
		"でする5秒前-majikoi.com": "xn--5-majikoi-z83h7ezr1858a1v9e.com",
		"Tạisaohọkh\\xf4ngthểchỉn\\xf3itiếngViệt.com": "xn--tisaohkhngthchnitingvit-kjcr8268qyxafd2f1b9g.com",
		"porquénopuedensimplementehablarenespañol":    "xn--porqunopuedensimplementehablarenespaol-fmd56a",
		"安室奈美恵-with-super-monkeys.de":                 "xn---with-super-monkeys-pc58ag80a8qai00g7n9n.de",
	}

	transitionalProcessing := false
	for domain, punycode := range expected {
		result := idntranslator.ToASCII(domain, transitionalProcessing)
		assert.Equal(t, punycode, result)
	}
}

func TestToUnicode(t *testing.T) {
	expected := map[string]string{
		"öbb.at":           "öbb.at",
		"Öbb.at":           "öbb.at",
		"ÖBB.at":           "öbb.at",
		"O\u0308bb.at":     "öbb.at",
		"xn--bb-eka.at":    "öbb.at",
		"faß.de":           "faß.de",
		"fass.de":          "fass.de",
		"xn--fa-hia.de":    "faß.de",
		"not=std3":         "not=std3",
		"\\ud83d\\udca9":   "💩",
		"\\ud87e\\udcca":   "𣀊",
		"fäß.de":           "fäß.de",
		"₹.com":            "₹.com",
		"𑀓.com":            "𑀓.com",
		"a‌b":              "a‌b",
		"ȡog.de":           "ȡog.de",
		"☕.de":             "☕.de",
		"I♥NY.de":          "i♥ny.de",
		"ＡＢＣ・日本.co.jp":     "abc・日本.co.jp",
		"日本｡co｡jp":         "日本.co.jp",
		"日本｡co．jp":         "日本.co.jp",
		"x\u0327\u0301.de": "x̧́.de",
		"x\u0301\u0327.de": "x̧́.de",
		"σόλος.gr":         "σόλος.gr",
		"Σόλος.gr":         "σόλος.gr",
		"ΣΌΛΟΣ.gr":         "σόλοσ.gr",
		"عربي.de":          "عربي.de",
		"نامهای.de":        "نامهای.de",
		"نامه\u200Cای.de":  "نامه‌ای.de",
		"\u3067\u3059\u308B5\u79D2\u524D-MajiKoi.com":                                                                                          "でする5秒前-majikoi.com",
		"\u5B89\u5BA4\u5948\u7F8E\u6075-with-SUPER-MONKEYS.de":                                                                                 "安室奈美恵-with-super-monkeys.de",
		"T\u1EA1isaoh\u1ECDkh\\xF4ngth\u1EC3ch\u1EC9n\\xF3iti\u1EBFngVi\u1EC7t":                                                                "tạisaohọkhôngthểchỉnóitiếngviệt",
		"Porqu\\xE9nopuedensimplementehablarenEspa\\xF1ol":                                                                                     "porquénopuedensimplementehablarenespañol",
		"\u05DC\u05DE\u05D4\u05D4\u05DD\u05E4\u05E9\u05D5\u05D8\u05DC\u05D0\u05DE\u05D3\u05D1\u05E8\u05D9\u05DD\u05E2\u05D1\u05E8\u05D9\u05EA": "למההםפשוטלאמדבריםעברית",
	}

	transitionalProcessing := false
	for domain, idn := range expected {
		result := idntranslator.ToUnicode(domain, transitionalProcessing)
		assert.Equal(t, idn, result)
	}
}

func TestToASCIIDefaultProcessingByTLD(t *testing.T) {
	// Default (no explicit argument): transitional for every TLD but the listed ones.
	assert.Equal(t, "fass.com", idntranslator.ToASCII("faß.com"))

	// The listed TLDs keep ß.
	for _, tld := range []string{"art", "be", "ca", "de", "swiss", "fr", "pm", "re", "tf", "wf", "yt"} {
		assert.Equal(t, "xn--fa-hia."+tld, idntranslator.ToASCII("faß."+tld), tld)
	}

	// Case-insensitive.
	assert.Equal(t, "xn--fa-hia.de", idntranslator.ToASCII("faß.DE"))

	// The match is anchored to the last label.
	assert.Equal(t, "fass.de.com", idntranslator.ToASCII("faß.de.com"))
	assert.Equal(t, "fass.dealer", idntranslator.ToASCII("faß.dealer"))
	assert.Equal(t, "fass.design", idntranslator.ToASCII("faß.design"))

	// An explicit argument overrides the TLD policy.
	assert.Equal(t, "fass.de", idntranslator.ToASCII("faß.de", true))
	assert.Equal(t, "xn--fa-hia.com", idntranslator.ToASCII("faß.com", false))
}

func TestConvertDefaultProcessingByTLD(t *testing.T) {
	// ToUnicode never maps ß (x/net forces non-transitional there), so IDN keeps it.
	assert.Equal(t, []Row{{IDN: "faß.com", PUNYCODE: "fass.com"}}, toTestRows(idntranslator.Convert("faß.com")))
	assert.Equal(t, []Row{{IDN: "faß.de", PUNYCODE: "xn--fa-hia.de"}}, toTestRows(idntranslator.Convert("faß.DE")))
	assert.Equal(t, []Row{{IDN: "faß.de.com", PUNYCODE: "fass.de.com"}}, toTestRows(idntranslator.Convert("faß.de.com")))
	assert.Equal(t, []Row{{IDN: "faß.dealer", PUNYCODE: "fass.dealer"}}, toTestRows(idntranslator.Convert("faß.dealer")))

	for _, tld := range []string{"art", "be", "ca", "de", "swiss", "fr", "pm", "re", "tf", "wf", "yt"} {
		assert.Equal(t, []Row{{IDN: "faß." + tld, PUNYCODE: "xn--fa-hia." + tld}},
			toTestRows(idntranslator.Convert("faß."+tld)), tld)
	}
}

func toTestRows(rows []idntranslator.Row) []Row {
	out := make([]Row, len(rows))
	for i, r := range rows {
		out[i] = Row{IDN: r.IDN, PUNYCODE: r.PUNYCODE}
	}
	return out
}
