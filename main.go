package main

import (
	"formtest/frame"
)

const (
	baseURL = "http://localhost:8080/api"
	formID  = "form-1"
)

func main() {
	api := &frame.APIer{Data1: -1, Data2: "initial value"}

	api.Form(formID, "")
	// api.Get(baseURL)

}
