// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package responseparser provides functionality to cover API response
// data parsing and serializing.
package responseparser

import (
	"regexp"
	"strings"
)

// Parse method to return plain API response parsed into hash format
func Parse(r string) map[string]interface{} {
	hash := make(map[string]interface{})
	tmp := strings.Split(strings.ReplaceAll(r, "\r", ""), "\n")
	p1 := regexp.MustCompile(`^([^\=]*[^\t\= ])[\t ]*=[\t ]*(.*)$`)
	p2 := regexp.MustCompile(`(?i)^property\[([^\]]*)\]\[([0-9]+)\]`)
	properties := make(map[string][]string)
	for _, row := range tmp {
		m := p1.MatchString(row)
		if m {
			groups := p1.FindStringSubmatch(row)
			property := strings.ToUpper(groups[1])
			mm := p2.MatchString(property)
			if mm {
				groups2 := p2.FindStringSubmatch(property)
				key := strings.ReplaceAll(strings.ToUpper(groups2[1]), "\\s", "")
				// idx2 := strconv.Atoi(groups2[2])
				list := make([]string, len(properties[key]))
				copy(list, properties[key])
				pat := regexp.MustCompile(`[\t ]*$`)
				rep1 := "${1}$2"
				list = append(list, pat.ReplaceAllString(groups[2], rep1))
				properties[key] = list
			} else {
				val := groups[2]
				if len(val) > 0 {
					pat := regexp.MustCompile(`[\t ]*$`)
					hash[property] = pat.ReplaceAllString(val, "")
				}
			}
		}
	}
	if len(properties) > 0 {
		hash["PROPERTY"] = properties
	}
	return hash
}
