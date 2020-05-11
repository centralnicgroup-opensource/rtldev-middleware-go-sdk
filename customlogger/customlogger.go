// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package customlogger provides functionality around debug outputs/logging of API communication
package customlogger

import (
	"fmt"

	L "github.com/hexonet/go-sdk/logger"
	R "github.com/hexonet/go-sdk/response"
)

// Logger is a struct representing logger for API communication.
type CustomLogger struct {
	*L.Logger
}

// NewLogger represents the constructor for struct Logger.
func NewCustomLogger() *CustomLogger {
	logger := &CustomLogger{}
	logger.Logger = L.NewLogger()
	return logger
}

// Log method to ouput/log api communication
func (c *CustomLogger) Log(post string, r *R.Response, errormsg ...string) {
	fmt.Printf("%s\n", r.GetCommandPlain())
	fmt.Printf("POST: %s\n", post)
	if len(errormsg) > 0 && len(errormsg[0]) > 0 {
		fmt.Printf("HTTP communication failed: %s\n", errormsg)
	}
	fmt.Println(r.GetPlain())
}
