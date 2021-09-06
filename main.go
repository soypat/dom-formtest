package main

import (
	"fmt"
	"formtest/frame"
	"reflect"

	"gonum.org/v1/gonum/spatial/r3"
	"honnef.co/go/js/dom/v2"
)

const (
	baseURL = "http://localhost:8085/api"
	formID  = "form-1"
)

func main() {
	var dats = struct {
		A string  `json:"dataStr"`
		B float64 `json:"dataNum"`
		V r3.Vec  `json:"vector"`
	}{}
	dats.A = "init"
	dats.B = -1
	api, err := frame.New("form1", baseURL, &dats)
	if err != nil {
		panic(err)
	}

	form := api.GenerateForm()
	logf("gen form")
	form.GetElementsByTagName("button")[0].AddEventListener("click", false, func(e dom.Event) {
		e.PreventDefault()
		err = api.Marshal()
		if err != nil {
			logf("error: %s", err)
		}
		logf("got data %#v", dats)
	})
	UpdateValues(api)
}

func UpdateValues(api *frame.APIer) {
	api.ForEachInput(func(he *dom.HTMLInputElement, v reflect.Value) {
		switch v.Kind() {
		case reflect.String:
			he.SetValue(v.String())
		case reflect.Float64:
			he.SetValue(fmt.Sprintf("%g", v.Float()))
		case reflect.Int:
			he.SetValue(fmt.Sprintf("%d", v.Int()))
		}
	})
}
