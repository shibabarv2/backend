package main

import (
	"errors"
	"os"
	"shiba-backend/structs"
	"strings"
)

// IsProduction is a validator function to check if production is true. (IF PRODUCTION IS FALSE, EVERY ROUTE WILL NOT REQUIRE A ADMIN KEY!!)
// This won't be in use until I completely integrate this.
func IsProduction() bool {

	if os.Getenv("production") == "true" {
		structs.IsTesting = false
		return true
	}

	structs.IsTesting = true
	return false

}

// CheckCorrectVariables is a validator function to check if any variables in the .env file are missing.
// This won't be in use until I completely integrate this.
func CheckCorrectVariables() []error {

	var Errors []error

	// CHECK 1: Admin key
	if len(os.Getenv("ADMIN_KEY")) > 9 {
		Errors = append(Errors, errors.New("Admin key is missing or is not longer than 9 characters."))
	}

	// CHECK 2: MongoDB key
	if len(os.Getenv("MONGODBURI")) > 1 {
		Errors = append(Errors, errors.New("MongoDB URI is missing."))
	}

	// CHECK 3: API url
	if len(os.Getenv("API_URL")) > 1 {
		Errors = append(Errors, errors.New("API url is missing."))
	}

	if !strings.Contains(os.Getenv("API_URL"), "https") {
		Errors = append(Errors, errors.New("API url does not use the https:// protocol. To disable this please edit the config.go CheckCorrectVariables function."))
	}

	return Errors
}
