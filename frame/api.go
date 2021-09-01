package frame

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"

	dom "honnef.co/go/js/dom/v2"
)

type APIer struct {
	Data1 float64 `json:"data1"`
	Data2 string  `json:"data2"`
}

func (a *APIer) Get(baseurl string) error {
	u, err := url.Parse(baseurl)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "get")
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	d := json.NewDecoder(resp.Body)
	return d.Decode(a)
}

func (a *APIer) Form(formElemID string) {
	doc := dom.GetWindow().Document()
	form := doc.GetElementByID(formElemID)
	if form == nil {
		logf("form not found")
		node := newSingleHTML(`<form id="` + formElemID + `"></form>`)
		body := doc.GetElementsByTagName("body")[0]
		logf("appending %#v", node)
		body.AppendChild(node)
		form = doc.GetElementByID(formElemID)
	}
	var spf = fmt.Sprintf
	branches := getBranch("json", a)
	logf("branched: %#v\nfrom %#v", branches, *a)
	var content string
	for i := range branches {
		switch branches[i].Kind() {
		case reflect.Float64, reflect.String:
			content += spf("<div><label>%[1]v</label><input id=\"%[1]s-%[2]d\" ></input></div>", branches[i].name, i, branches[i].String())
		default:
			continue
		}
	}

	appendChildren(form, newHTML(content))

	form.AddEventListener(EventSubmit, false, func(e dom.Event) {
		logf("type %v", e.Type())
		logf("hi")
	})
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
