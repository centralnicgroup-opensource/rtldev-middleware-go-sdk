package idntranslator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/net/idna"
	"golang.org/x/text/unicode/norm"
)

// Row represents a row in the translation result.
type Row struct {
	IDN      string
	PUNYCODE string
}

// interfaceToStringSlice converts the input interface to a slice of strings.
func interfaceToStringSlice(input interface{}) []string {
	switch v := input.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return nil
	}
}

// Convert converts a domain string or a slice of domain strings between Unicode and Punycode formats.
func Convert(domainOrDomains interface{}) []Row {
	domains := interfaceToStringSlice(domainOrDomains)

	var translated []Row

	for _, domain := range domains {
		idn, punycode := handleConversion(domain)
		translated = append(translated, Row{IDN: idn, PUNYCODE: punycode})
	}

	return translated
}

// handleConversion handles conversion of a keyword between Unicode and Punycode formats.
func handleConversion(keyword string) (string, string) {
	if keyword == "" {
		return "", ""
	}

	return ToUnicode(keyword), ToASCII(keyword)
}

// ToUnicode converts a domain string to Unicode format.
func ToUnicode(asciiString string, transitionalProcessing ...bool) string {
	decodedKeyword := decodeUnicodeEscapes(asciiString)
	// Define the IDNA options
	opts := idna.New(
		idna.MapForLookup(),
		idna.Transitional(isTransitionalProcessing(asciiString, transitionalProcessing...)), // Map ß -> ss
		idna.StrictDomainName(false)) // Set more permissive ASCII rules.

	// Convert the Unicode string to Punycode using the specified options
	unicode, err := opts.ToUnicode(decodedKeyword)
	if err != nil {
		// Handle the error appropriately
		return asciiString // Return the original string if conversion fails
	}
	return unicode
}

// ToASCII converts a Unicode string to Punycode format.
func ToASCII(unicodeString string, transitionalProcessing ...bool) string {
	unicodeString = ToUnicode(unicodeString, transitionalProcessing...)
	// Define the IDNA options
	opts := idna.New(
		idna.MapForLookup(),
		idna.Transitional(isTransitionalProcessing(unicodeString, transitionalProcessing...)), // Map ß -> ss
	)

	// Convert the Unicode string to Punycode using the specified options
	punycode, err := opts.ToASCII(unicodeString)
	if err != nil {
		// Handle the error appropriately
		return unicodeString // Return the original string if conversion fails
	}

	return punycode
}

// DecodeUnicodeEscapes decodes Unicode escape sequences in a string, normalizes it, and converts it to lowercase.
func decodeUnicodeEscapes(unicodeString string) string {
	decoded := decodeUnicodeEscapeSequences(unicodeString)
	normalized := normalizeAndLowerCase(decoded)
	return normalized
}

// isTransitionalProcessing checks if the provided top-level domain (TLD) is non-transitional.
func isTransitionalProcessing(keyword string, transitionalProcessing ...bool) bool {
	if len(transitionalProcessing) > 0 {
		return transitionalProcessing[0]
	}

	transitionalTLDs := []string{"art", "be", "ca", "de", "fr", "pm", "re", "swiss", "tf", "wf", "yt"}
	regex := `\.(` + strings.Join(transitionalTLDs, "|") + `)\.?`
	re := regexp.MustCompile(regex)
	return re.MatchString(strings.ToLower(keyword))
}

// encodeHexToUnicode encodes hexadecimal escape sequences to Unicode escape sequences in a string.
// It replaces occurrences of double backslashes with single backslashes,
// and converts hexadecimal escape sequences like \xFC to their corresponding Unicode escape sequences like \u00FC.
func encodeHexToUnicode(inputString string) string {
	inputString = regexp.MustCompile(`\\{2}`).ReplaceAllString(inputString, "\\")
	// Define a regular expression to match hexadecimal escape sequences like \xXX
	reHex := regexp.MustCompile(`\\x([0-9a-fA-F]{2})`)
	// Replace hexadecimal escape sequences with their corresponding Unicode escape sequences
	encodedString := reHex.ReplaceAllStringFunc(inputString, func(match string) string {
		// Extract the hexadecimal value
		hexValue := match[2:]

		// Parse the hexadecimal value to an integer
		// Note: Error handling omitted for brevity
		intValue, err := strconv.ParseInt(hexValue, 16, 64)
		if err != nil {
			// Handle the error appropriately
			return match // Return the original string if conversion fails
		}
		// Convert the integer to a Unicode escape sequence
		return fmt.Sprintf("\\u%04X", intValue)
	})

	return encodedString
}

// decodeUnicodeEscapeSequences decodes Unicode escape sequences in a string.
// It decodes Unicode escape sequences like \u00FC to their corresponding characters like ü.
func decodeUnicodeEscapeSequences(unicodeString string) string {
	// Regular expression to match Unicode escape sequences and surrogate pairs
	reUnicode := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	decoded := encodeHexToUnicode(unicodeString)
	decoded = reUnicode.ReplaceAllStringFunc(decoded, func(match string) string {
		hexValue := match[2:]
		code, err := strconv.ParseUint(hexValue, 16, 32)
		if err != nil {
			return match // Return the original string if parsing fails
		}
		r := rune(code)
		if utf16.IsSurrogate(r) {
			return match // Leave surrogate pairs to be processed together
		}
		return string(r)
	})

	// Process surrogate pairs: Surrogate pairs are used in UTF-16 encoding to represent characters outside the Basic Multilingual Plane (BMP).
	// These characters are represented by pairs of 16-bit code units called surrogates.
	// In Unicode escape sequences, surrogate pairs are represented as two consecutive escape sequences: \udXXX\udYYY.
	// This regular expression captures these surrogate pair patterns.
	reSurrogatePair := regexp.MustCompile(`\\u(d[89ab][0-9a-fA-F]{2})\\u(d[c-f][0-9a-fA-F]{2})`)
	decoded = reSurrogatePair.ReplaceAllStringFunc(decoded, func(match string) string {
		// Extract the hexadecimal values for the surrogate pair
		// The first value represents the high surrogate, and the second represents the low surrogate
		r1, err := strconv.ParseUint(match[2:6], 16, 32)
		if err != nil {
			// If parsing fails, keep the original string
			return match
		}
		r2, err := strconv.ParseUint(match[8:12], 16, 32)
		if err != nil {
			// If parsing fails, keep the original string
			return match
		}
		// Combine the two code points into a single Unicode character
		// This is necessary because certain characters are represented by pairs of code points
		runeValue := utf16.DecodeRune(rune(uint16(r1)), rune(uint16(r2)))
		// If the resulting character is invalid, keep the original surrogate pair
		if runeValue == utf8.RuneError {
			return match
		}
		// Return the decoded Unicode character
		return string(runeValue)
	})

	return decoded
}

// normalizeAndLowerCase normalizes the string using NFC normalization form and converts it to lowercase.
func normalizeAndLowerCase(input string) string {
	// Normalize the string using NFC normalization form
	normalized := norm.NFC.String(input)

	// replace full width characters with normalized e.g. ＡＢＣ -> abc
	normalized = replaceFullWidthChars(normalized)

	// Convert to lowercase
	lowercase := strings.ToLower(normalized)

	return lowercase
}

// ReplaceFullWidthChars replaces full-width characters with their corresponding normal-width counterparts.
func replaceFullWidthChars(str string) string {
	var sb strings.Builder
	const fullWidthOffset = 0xfee0
	for _, r := range str {
		switch {
		case r >= 0xFF01 && r <= 0xFF5E:
			// Map full-width characters to their corresponding normal-width characters
			sb.WriteRune(r - fullWidthOffset)
		case r == '｡':
			// Replace full-width dot character with the regular dot character
			sb.WriteRune('.')
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
