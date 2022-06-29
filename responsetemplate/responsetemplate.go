// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package responsetemplate provides basic functionality to handle API response data
package responsetemplate

import (
	"strconv"

	rp "github.com/hexonet/go-sdk/v3/responseparser"
)

// ResponseTemplate is a struct used to cover basic functionality to work with
// API response data (or hardcoded API response data).
type ResponseTemplate struct {
	Raw  string
	Hash map[string]interface{}
}

// NewResponseTemplate represents the constructor for struct ResponseTemplate.
func NewResponseTemplate(raw string) *ResponseTemplate {
	if len(raw) == 0 {
		raw = "[RESPONSE]\r\nCODE=423\r\nDESCRIPTION=Empty API response. Probably unreachable API end point {CONNECTION_URL}\r\nEOF\r\n"
	}
	rt := &ResponseTemplate{
		Raw:  raw,
		Hash: rp.Parse(raw),
	}
	_, exists1 := rt.Hash["CODE"]
	_, exists2 := rt.Hash["DESCRIPTION"]
	if !exists1 || !exists2 {
		rt.Raw = "[RESPONSE]\r\nCODE=423\r\nDESCRIPTION=Invalid API response. Contact Support\r\nEOF\r\n"
		rt.Hash = rp.Parse(rt.Raw)
	}
	return rt
}

// GetCode method to return the API response code
func (rt *ResponseTemplate) GetCode() int {
	h := rt.GetHash()
	c, err := strconv.Atoi(h["CODE"].(string))
	if err == nil {
		return c
	}
	return 421
}

// GetDescription method to return the API response description
func (rt *ResponseTemplate) GetDescription() string {
	h := rt.GetHash()
	return h["DESCRIPTION"].(string)
}

// GetPlain method to return raw API response
func (rt *ResponseTemplate) GetPlain() string {
	return rt.Raw
}

// GetQueuetime method to return API response queuetime
func (rt *ResponseTemplate) GetQueuetime() float64 {
	h := rt.GetHash()
	if val, ok := h["QUEUETIME"]; ok {
		f, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return f
		}
	}
	return 0.00
}

// GetHash method to return API response in hash format
func (rt *ResponseTemplate) GetHash() map[string]interface{} {
	return rt.Hash
}

// GetRuntime method to return API response runtime
func (rt *ResponseTemplate) GetRuntime() float64 {
	h := rt.GetHash()
	if val, ok := h["RUNTIME"]; ok {
		f, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return f
		}
	}
	return 0.00
}

// IsError method to check if API response represents an error case
func (rt *ResponseTemplate) IsError() bool {
	c := rt.GetCode()
	return (c >= 500 && c <= 599)
}

// IsSuccess method to check if API response represents a success case
func (rt *ResponseTemplate) IsSuccess() bool {
	c := rt.GetCode()
	return (c >= 200 && c <= 299)
}

// IsTmpError method to check if current API response represents a temporary error case
func (rt *ResponseTemplate) IsTmpError() bool {
	c := rt.GetCode()
	return (c >= 400 && c <= 499)
}

// IsPending method to check if current operation is returned as pending
func (rt *ResponseTemplate) IsPending() bool {
	h := rt.GetHash()
	if val, ok := h["PENDING"]; ok {
		if val.(string) == "1" {
			return true
		}
	}
	return false
}
