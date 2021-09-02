package frame

import (
	"encoding/json"
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

const tagname = "json"

type APIer struct {
	Data1 float64 `json:"data1"`
	Data2 string  `json:"data2"`
}

func (a *APIer) PostJSON(baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "get")
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	go enc.Encode(a)
	_, err = http.Post(u.String(), "json", r)
	return err
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

func (a *APIer) IDs() (ids []string) {
	branches := getBranch(tagname, a)
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

func (a *APIer) Form(formElemID, baseURL string) {
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
	ids := a.IDs()
	for i := range branches {
		switch branches[i].Kind() {
		case reflect.Float64, reflect.String:
			content += spf("<div><label>%[1]v</label><input id=\"%[2]s\" ></input></div>", branches[i].name, ids[i], branches[i].String())
		default:
			continue
		}
	}
	logf("try")
	content += "<div><button type=\"submit\">Submit</button></div>"
	logf("%#v", newHTML(content))
	appendChildren(form, newHTML(content))

	form.AddEventListener(EventSubmit, false, func(e dom.Event) {
		a.
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
