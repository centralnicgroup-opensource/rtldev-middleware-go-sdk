// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package responsetemplatemanager provides basic functionality to handle API response data
package responsetemplatemanager

import (
	"strings"
	"sync"

	rp "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responseparser"
)

// ResponseTemplateManager is a struct used to cover basic functionality to work with
// API response templates.
type ResponseTemplateManager struct {
	Templates map[string]string
}

var instance *ResponseTemplateManager
var once sync.Once

// GetInstance method to return the responsetemplatemanager singleton instance
func GetInstance() *ResponseTemplateManager {
	once.Do(func() {
		instance = &ResponseTemplateManager{
			Templates: map[string]string{
				"404":          generateTemplate("421", "Page not found"),
				"500":          generateTemplate("500", "Internal server error"),
				"empty":        generateTemplate("423", "Empty API response. Probably unreachable API end point {CONNECTION_URL}"),
				"error":        generateTemplate("421", "Command failed due to server error. Client should try again"),
				"expired":      generateTemplate("530", "SESSION NOT FOUND"),
				"httperror":    generateTemplate("421", "Command failed due to HTTP communication error"),
				"unauthorized": generateTemplate("530", "Unauthorized"),
				"invalid":      generateTemplate("423", "Invalid API response. Contact Support"),
			},
		}
	})
	return instance
}

// generateTemplate method to generate API a response template string
// for given code and description
func generateTemplate(code string, description string) string {
	var tmp strings.Builder
	tmp.WriteString("[RESPONSE]\r\nCODE=")
	tmp.WriteString(code)
	tmp.WriteString("\r\nDESCRIPTION=")
	tmp.WriteString(description)
	tmp.WriteString("\r\nEOF\r\n")
	return tmp.String()
}

// GenerateTemplate method to generate API a response template string
// for given code and description
func (rtm *ResponseTemplateManager) GenerateTemplate(code string, description string) string {
	return generateTemplate(code, description)
}

// AddTemplate method to add a template to the templates container
func (rtm *ResponseTemplateManager) AddTemplate(id string, plain string) *ResponseTemplateManager {
	rtm.Templates[id] = plain
	return rtm
}

// GetTemplate method to get a ResponseTemplate from templates container
func (rtm *ResponseTemplateManager) GetTemplate(id string) string {
	if rtm.HasTemplate(id) {
		return rtm.Templates[id]
	}
	return generateTemplate("500", "Response Template not found")
}

// GetTemplates method to return a map covering all available response templates
func (rtm *ResponseTemplateManager) GetTemplates() map[string]string {
	tpls := map[string]string{}
	for key := range rtm.Templates {
		tpls[key] = rtm.Templates[key]
	}
	return tpls
}

// HasTemplate method to check if given template id exists in template container
func (rtm *ResponseTemplateManager) HasTemplate(id string) bool {
	if _, ok := rtm.Templates[id]; ok {
		return true
	}
	return false
}

// IsTemplateMatchHash method to check if given API response hash matches a given
// template by code and description
func (rtm *ResponseTemplateManager) IsTemplateMatchHash(tpl2 map[string]interface{}, id string) bool {
	h := rp.Parse(rtm.GetTemplate(id))
	return ((h["CODE"] == tpl2["CODE"].(string)) &&
		(h["DESCRIPTION"] == tpl2["DESCRIPTION"].(string)))
}

// IsTemplateMatchPlain method to check if given API plain response matches a given
// template by code and description
func (rtm *ResponseTemplateManager) IsTemplateMatchPlain(plain string, id string) bool {
	tpl2 := rp.Parse(plain)
	return rtm.IsTemplateMatchHash(tpl2, id)
}
