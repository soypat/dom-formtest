package frame

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	dom "honnef.co/go/js/dom/v2"
)

type APIer struct {
	formID  string
	baseURL url.URL
	dataPtr interface{}
	form    dom.Element
}

func New(label, baseURL string, dataptr interface{}) (*APIer, error) {
	a := APIer{}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if dataptr == nil || reflect.Indirect(reflect.ValueOf(dataptr)) == reflect.ValueOf(dataptr) {
		return nil, errors.New("data must be pointer")
	}
	a.formID = label
	a.baseURL = *u
	a.dataPtr = dataptr
	return &a, nil
}

func (a *APIer) marshalForm() error {
	inputs := dom.GetWindow().Document().GetElementByID(a.formID).GetElementsByTagName("input")
	nameValues := make(map[string]string)
	for _, input := range inputs {
		input := input.(*dom.HTMLInputElement) // cast to input element to have access to special methods
		nameValues[input.GetAttribute("name")] = input.Value()
	}

	return nil
}

func (a *APIer) PostForm() error {
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

func (a *APIer) IDs() (ids []string) {
	branches := a.branch()
	for i := range branches {
		switch branches[i].Kind() {
		case reflect.Float64, reflect.String:
			ids = append(ids, branches[i].name+"-n"+strconv.Itoa(i))
		default:
			continue
		}
	}
	return ids
}

func (a *APIer) GenerateForm() dom.Element {
	doc := dom.GetWindow().Document()
	form := doc.GetElementByID(a.formID)

	if form == nil { // create form element if not exist
		node := newSingleHTML(`<form id="` + a.formID + `"></form>`)
		body := doc.GetElementsByTagName("body")[0]
		body.AppendChild(node)
		form = doc.GetElementByID(a.formID)
	}
	// create inner input elements of form.
	appendChildren(form, newHTML(a.inputsHTML()))

	form.AddEventListener(EventSubmit, false, func(e dom.Event) {
		e.PreventDefault() // prevent navigation away from page
		a.marshalForm()
	})
	a.form = form
	return form
}

// Creates inner input elements
// input elements will have id and name attributes defined.
func (a APIer) inputsHTML() string {
	var spf = fmt.Sprintf
	branches := a.branch()
	var content string
	ids := a.IDs()
	for i := range branches {
		switch branches[i].Kind() {
		case reflect.Float64, reflect.String:
			content += spf("<div><label>%[1]v</label><input id=\"%[2]s\" name=\"%[1]\" ></input></div>", branches[i].name, ids[i])
		default:
			continue
		}
	}
	content += "<div><button type=\"submit\">Submit</button></div>"
	return content
}

func (a APIer) branch() []branch {
	return getBranch("json", a.dataPtr)
}

func newSingleHTML(html string) dom.Node {
	div := dom.GetWindow().Document().CreateElement("div")
	trimmed := strings.TrimSpace(html)
	div.SetInnerHTML(trimmed)
	return div.FirstChild()
}

func newHTML(html string) []dom.Node {
	div := dom.GetWindow().Document().CreateElement("div")
	trimmed := strings.TrimSpace(html)
	div.SetInnerHTML(trimmed)
	return div.ChildNodes()
}

func appendChildren(target dom.Node, newChildren []dom.Node) {
	for _, v := range newChildren {
		target.AppendChild(v)
	}
}
