package frame

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
)

type APIer struct {
	formID  string
	baseURL url.URL
	dataPtr interface{}
	form    *formElement
}

func New(formID, baseURL string, dataptr interface{}) (*APIer, error) {
	a := APIer{}
	if dataptr == nil || reflect.Indirect(reflect.ValueOf(dataptr)) == reflect.ValueOf(dataptr) {
		return nil, errors.New("data must be pointer")
	}
	a.formID = formID

	a.dataPtr = dataptr
	u, err := url.Parse(baseURL)
	if err == nil {
		a.baseURL = *u
	}
	return &a, nil
}

func (a *APIer) PostData() error {
	urlpath := path.Join(a.baseURL.Path, "post")
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	go enc.Encode(a.dataPtr)
	_, err := http.Post(urlpath, "json", r)
	return err
}

func (a *APIer) Get() error {
	urlPath := path.Join(a.baseURL.Path, "get")
	resp, err := http.Get(urlPath)
	if err != nil {
		return err
	}
	d := json.NewDecoder(resp.Body)
	return d.Decode(a.dataPtr)
}
