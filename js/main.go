package main

import (
	"fmt"
	"reflect"

	"github.com/soypat/dom-formtest/js/frame"

	"honnef.co/go/js/dom/v2"
)

const (
	baseURL = "http://localhost:8085/api"
	formID  = "form-1"
)

func main() {
	var dats = Parameters{}
	// dats.A = "init"
	// dats.B = -1
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
	FormatForm(api)
}

func UpdateValues(api *frame.APIer) {
	api.ForEachInput(func(he frame.Input, v reflect.Value) {
		ip := he.Input()
		switch v.Kind() {
		case reflect.String:
			ip.SetValue(v.String())
		case reflect.Float64:
			ip.SetValue(fmt.Sprintf("%g", v.Float()))
		case reflect.Int:
			ip.SetValue(fmt.Sprintf("%d", v.Int()))
		}
	})
}

func FormatForm(api *frame.APIer) {
	api.Form().Class().Add("section")
	api.ForEachField(func(he frame.Field, v reflect.Value) {
		he.MainDiv().Class().Set([]string{"field"})

		he.Label().Class().Set([]string{"label"})
		he.Label().SetTextContent(fmt.Sprintf("%s [%v]", he.Name(), v.Type().String()))
	})

	api.ForEachInput(func(el frame.Input, v reflect.Value) {
		el.MainDiv().Class().Set([]string{"control", "block"})
		el.Input().Class().Set([]string{"input", "level-item", "is-small"})
		el.Span().Class().Set([]string{"help"})
		el.Span().SetTextContent(el.Name())
	})
}
