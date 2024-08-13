package main

import (
	"fmt"
	"os"

	CL "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/apiclient"
)

func main() {
	// sessionless API communication
	fmt.Println("--- SESSIONLESS API COMMUNICATION ---")
	cl := CL.NewAPIClient()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.UseOTESystem()
	cl.EnableDebugMode()
	r := cl.Request(map[string]interface{}{
		"COMMAND": "StatusAccount",
	})
	if r.IsSuccess() {
		fmt.Println("Command succeeded.")
	} else {
		fmt.Println("Command failed.")
	}
	fmt.Println()

	// session based API communication
	fmt.Println("--- SESSION BASED API COMMUNICATION ---")
	cl = CL.NewAPIClient()
	cl.SetCredentials(os.Getenv("CNR_TEST_USER"), os.Getenv("CNR_TEST_PASSWORD"))
	cl.UseOTESystem()
	cl.EnableDebugMode()
	r = cl.Login()
	if r.IsSuccess() {
		fmt.Println("Login succeeded.")
		r = cl.Request(map[string]interface{}{
			"COMMAND": "StatusAccount",
		})
		if r.IsSuccess() {
			fmt.Println("Command succeeded.")
			r = cl.Logout()
			if r.IsSuccess() {
				fmt.Println("Logout succeeded.")
			} else {
				fmt.Println("Logout failed.")
			}
		} else {
			fmt.Println("Command failed.")
		}
	} else {
		fmt.Println("Login failed.")
	}
}
