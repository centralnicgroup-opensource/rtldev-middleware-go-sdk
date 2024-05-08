// https://pkg.go.dev/golang.org/x/net/idna

// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package idntranslator provides basic functionality to customize the API response description
package idntranslator

type IdnTranslatorRow struct {
	Idn string
	Punycode string
}

// Convert function for converting a domain to idn + punycode
func Convert(domains []string, options map[string]string) []IdnTranslatorRow {
	translated := []IdnTranslatorRow{};

	for idx, domain range domains {
		translated[idx] = IdnTranslatorRow{
			Idn: IdnTranslator.toUnicode(domain, options),
			Punycode: IdnTranslator.toASCII(domain, options),
		}
	}

	return translated;
}

func ToUnicode(domain string, options map[string]string) string {
	idn := domain
	return idn
}

func ToASCII(domain string, options map[string]string) string {
	ascii := domain
	return ascii
}
