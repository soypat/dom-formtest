package main

import (
	"formtest/frame"
)

const (
	baseURL = "http://localhost:8080/api"
	formID  = "form-1"
)

func main() {
	var dats = struct {
		A string  `json:"dataStr"`
		B float64 `json:"dataNum"`
	}{}
	api, err := frame.New("form1", baseURL, &dats)
	if err != nil {
		panic(err)
	}

	api.GenerateForm()
	// api.Get(baseURL)

}
