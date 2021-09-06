// +build js

package frame

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	dom "honnef.co/go/js/dom/v2"
)

type formElement = dom.HTMLFormElement

// Marshal writes form data to apier data.
func (a *APIer) Marshal() error {
	a.ForEachInput(func(el *dom.HTMLInputElement, field reflect.Value) {
		switch field.Kind() {
		case reflect.String:
			field.SetString(el.Value())
		case reflect.Float64:
			val, err := strconv.ParseFloat(el.Value(), 64)
			if err == nil {
				field.SetFloat(val)
			}
		case reflect.Int:
			val, err := strconv.Atoi(el.Value())
			if err == nil {
				field.SetInt(int64(val))
			}
		}
	})
	logf("marshal data %#v", a.dataPtr)
	return nil
}

func (a *APIer) GenerateForm() *dom.HTMLFormElement {
	doc := dom.GetWindow().Document()
	form := doc.GetElementByID(a.formID)

	if form == nil { // create form element if not exist
		node := newSingleHTML(`<form id="` + a.formID + `"></form>`)
		body := doc.GetElementsByTagName("body")[0]
		body.AppendChild(node)
		form = doc.GetElementByID(a.formID)
	}

	// create inner input elements of form.
	appendChildren(form, a.newInnerFormHTML().ChildNodes())

	a.form = form.(*dom.HTMLFormElement)
	return a.form
}

func (a APIer) ForEachInput(f func(*dom.HTMLInputElement, reflect.Value)) {
	doc := dom.GetWindow().Document()
	walkStruct("json", a.dataPtr, a.formID, func(tree string, field reflect.Value) {
		k := field.Kind()
		if k != reflect.Int && k != reflect.String && k != reflect.Float64 {
			return
		}
		id, _ := getIDFromTree(tree)
		f(doc.GetElementByID(id).(*dom.HTMLInputElement), field)
	})
}

func (a APIer) newInnerFormHTML() dom.Node {
	maindiv := newSingleHTML("<div></div>")
	doc := dom.GetWindow().Document()
	walkStruct("json", a.dataPtr, a.formID, func(tree string, field reflect.Value) {
		id, parent := getIDFromTree(tree)
		names := strings.Split(tree, ".")
		var parentDiv dom.Node
		if parent != "" {
			parentDiv = doc.GetElementByID(parent)
		} else {
			parentDiv = maindiv
		}
		if field.Kind() == reflect.Struct {
			parentDiv.AppendChild(newSingleHTML(fmt.Sprintf("<div id=\"%s\"></div>", id)))
		} else {
			parentDiv.AppendChild(newSingleHTML(fmt.Sprintf("<div><label>%[1]v</label><input id=\"%[2]s\"></input></div>", names[len(names)-1], id)))
		}
	})
	maindiv.AppendChild(newSingleHTML("<div><button type=\"submit\">Submit</button></div>"))
	return maindiv
}

func getIDFromTree(tree string) (string, string) { // return current tree ID and parent's ID
	sp := strings.Split(tree, ".")
	if len(sp) < 2 {
		return strings.Join(sp, "-"), ""
	}
	return strings.Join(sp, "-"), strings.Join(sp[:len(sp)-1], "-")
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
