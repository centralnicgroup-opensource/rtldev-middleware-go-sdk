// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package socketconfig provides apiconnector client connection settings
package socketconfig

import (
	"net/url"
	"strings"
)

// SocketConfig is a struct representing connection settings used as POST data for http request against the insanely fast HEXONET backend API.
type SocketConfig struct {
	login      string
	pw         string
	session    string
	persistent string
}

// NewSocketConfig represents the constructor for struct SocketConfig.
func NewSocketConfig() *SocketConfig {
	sc := &SocketConfig{
		login:      "",
		persistent: "",
		pw:         "",
		session:    "",
	}
	return sc
}

// GetPOSTData method to return the struct data ready to submit within
// POST request of type "application/x-www-form-urlencoded"
func (s *SocketConfig) GetPOSTData() string {
	var tmp strings.Builder
	if len(s.login) > 0 {
		tmp.WriteString(url.QueryEscape("s_login"))
		tmp.WriteString("=")
		tmp.WriteString(url.QueryEscape(s.login))
		tmp.WriteString("&")
	}
	if len(s.pw) > 0 {
		tmp.WriteString(url.QueryEscape("s_pw"))
		tmp.WriteString("=")
		tmp.WriteString(url.QueryEscape(s.pw))
		tmp.WriteString("&")
	}
	if len(s.session) > 0 {
		tmp.WriteString(url.QueryEscape("s_sessionid"))
		tmp.WriteString("=")
		tmp.WriteString(url.QueryEscape(s.session))
		tmp.WriteString("&")
	}
	if len(s.persistent) > 0 {
		tmp.WriteString(url.QueryEscape("persistent"))
		tmp.WriteString("=")
		tmp.WriteString(url.QueryEscape(s.persistent))
		tmp.WriteString("&")
	}
	return tmp.String()
}

// GetSession method to return the session id currently in use.
func (s *SocketConfig) GetSession() string {
	return s.session
}

// GetLogin method to return the login id currently in use.
func (s *SocketConfig) GetLogin() string {
	return s.login
}

// SetLogin method to set username to use for api communication
func (s *SocketConfig) SetLogin(value string) *SocketConfig {
	s.login = value
	return s
}

// Persistent method for session to use for api communication
func (s *SocketConfig) SetPersistent() *SocketConfig {
	s.session = ""
	s.persistent = "1"
	return s
}

// SetPassword method to set password to use for api communication
func (s *SocketConfig) SetPassword(value string) *SocketConfig {
	s.session = ""
	s.pw = value
	return s
}

// SetSession method to set a API session id to use for api communication instead of credentials
// which is basically required in case you plan to use session based communication or if you want to use 2FA
func (s *SocketConfig) SetSession(sessionid string) *SocketConfig {
	s.pw = ""
	s.persistent = ""
	s.session = sessionid
	return s
}
