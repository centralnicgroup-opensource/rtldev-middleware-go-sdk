// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package responsetranslator provides basic functionality to customize the API response description
package responsetranslator

import (
	"regexp"
	"strings"

	RTM "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responsetemplatemanager"
)

type ResponseTranslator struct {
}

var descriptionRegexMap = map[string]string{
	// HX | CNR?
	"Authorization failed; Operation forbidden by ACL": "Authorization failed; Used Command `{COMMAND}` not white-listed by your Access Control List",
	// CNR
	"Domain status does not allow for operation":                                          "This Domain is locked. Initiating a Transfer is therefore impossible.",
	"Authorization failed [Invalid authorization information]":                            "The given Authorization Code is wrong. Initiating a Transfer is therefore impossible.",
	"Missing required attribute; premium domain name. please provide required parameters": "Confirm the Premium pricing by providing the necessary premium domain price data.",
}

var descriptionRegexMapSkipQuote = map[string]string{
	// all cases goes here that need translation e.g.
	// HX | CNR?
	`Invalid attribute value syntax; resource record \[(.+)\]`:                  "Invalid Syntax for DNSZone Resource Record: $1",
	`Missing required attribute; CLASS(?:=| \[MUST BE )PREMIUM_([\w\+]+)[\s\]]`: "Confirm the Premium pricing by providing the parameter CLASS with the value PREMIUM_$1.",
	`Syntax error in Parameter DOMAIN \((.+)\)`:                                 "The Domain Name $1 is invalid.",
}

// Translate function for plain api response
func Translate(raw string, cmd map[string]string, phs ...map[string]string) string {
	ph := map[string]string{}
	if len(phs) > 0 {
		ph = phs[0]
	}

	httperror := ""
	newraw := raw
	if len(raw) == 0 {
		newraw = "empty"
	}
	// Hint: Empty API Response (replace {CONNECTION_URL} later)

	// curl error handling
	isHTTPError := false
	if strings.HasPrefix(newraw, "httperror|") {
		isHTTPError = true
		httperror = strings.Replace(newraw, "httperror|", "", 1)
		newraw = "httperror"
	}

	// Explicit call for a static template
	rtm := RTM.GetInstance()
	if rtm.HasTemplate(newraw) {
		// don't use getTemplate as it leads to endless loop as of again
		// creating a response instance
		newraw = rtm.Templates[newraw]
		if isHTTPError && len(httperror) > 0 {
			newraw = strings.ReplaceAll(newraw, "{HTTPERROR}", " ("+httperror+")")
		}
	}

	if rtm.HasTemplate("invalid") {
		// Missing CODE or DESCRIPTION in API Response
		pattern1 := regexp.MustCompile(`(?i)description[\s]*=`)
		pattern2 := regexp.MustCompile(`(?i)code[\s]*=`)
		pattern3 := regexp.MustCompile(`(?i)description[\s]*=\r\n`)

		if pattern1.FindString(newraw) == "" || pattern2.FindString(newraw) == "" || pattern3.FindString(newraw) != "" {
			newraw = rtm.Templates["invalid"]
		}
	}

	// Iterate through the description-to-regex mapping
	// generic API response description rewrite
	// Iterate through the description-to-regex mapping with quoted regex
	data := ""
	for regex, val := range descriptionRegexMap {
		// Escape the regex pattern and attempt to find a match
		escapedRegex := regexp.QuoteMeta(regex)
		data = FindMatch(escapedRegex, newraw, val, cmd, ph)

		// If a match is found, exit the inner loop
		if len(data) > 0 {
			newraw = data
			break
		}
	}

	// Iterate through the description-to-regex mapping without quotes
	for regex, val := range descriptionRegexMapSkipQuote {
		data = FindMatch(regex, newraw, ""+val, cmd, ph)

		// If a match is found, exit the inner loop
		if len(data) > 0 {
			newraw = data
			break
		}
	}

	pattern := regexp.MustCompile(`\{.+\}`)
	return pattern.ReplaceAllString(newraw, "")
}

func FindMatch(regex string, newraw string, val string, cmd map[string]string, ph map[string]string) string {
	// match the response for given description
	// NOTE: we match if the description starts with the given description
	// it would also match if it is followed by additional text
	ret := ""
	qregex := regexp.MustCompile(`(?i)description[\s]*=[\s]*` + regex + `([^\r\n]+)?`)

	if qregex.FindString(newraw) != "" {
		// Replace any placeholders dynamically in val using cmd and ph maps
		for key, replaceVal := range cmd {
			val = strings.ReplaceAll(val, "{"+key+"}", replaceVal)
		}
		for key, replaceVal := range ph {
			val = strings.ReplaceAll(val, "{"+key+"}", replaceVal)
		}

		// If $newraw matches $qregex, replace with "description=" + val
		tmp := qregex.ReplaceAllString(newraw, "description="+val)
		if newraw != tmp {
			ret = tmp
		}
	}

	// Skipquote entries should not replace placeholder variables
	if strings.HasPrefix(val, "SkipPregQuote") {
		return ret
	}

	// Compile the regular expression to match placeholders
	vregex := regexp.MustCompile(`\{[^}]+\}`)

	// Keep track of whether any replacements were made
	replacementsMade := false

	// Replace placeholders with corresponding values from the map
	newret := ret
	if ret == "" {
		newret = newraw
	}
	newret = vregex.ReplaceAllStringFunc(newret, func(match string) string {
		placeholder := match[1 : len(match)-1] // Extract the placeholder name

		// Check if the placeholder exists in the map
		if value, exists := ph[placeholder]; exists {
			replacementsMade = true
			return value
		}
		return "" // Remove placeholder if not found
	})

	// Return an newret if replacements were made
	if replacementsMade {
		return newret
	}

	return ret
}
