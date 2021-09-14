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
type Input struct {
	div *dom.HTMLDivElement
}

func (el Input) MainDiv() *dom.HTMLDivElement {
	return el.div
}

func (el Input) Input() *dom.HTMLInputElement {
	return el.div.GetElementsByTagName("input")[0].(*dom.HTMLInputElement)
}

func (el Input) Span() *dom.HTMLSpanElement {
	return el.div.GetElementsByTagName("span")[0].(*dom.HTMLSpanElement)
}

func (el Input) Name() string {
	name, _ := getNameFromID(el.div.ID())
	return name
}

type Field struct {
	div *dom.HTMLDivElement
}

func (el Field) MainDiv() *dom.HTMLDivElement {
	return el.div
}

func (el Field) Label() *dom.HTMLLabelElement {
	return el.div.GetElementsByTagName("label")[0].(*dom.HTMLLabelElement)
}

func (el Field) Inputs() (inputs []*dom.HTMLInputElement) {
	for _, v := range el.div.ChildNodes() {
		v, ok := v.(*dom.HTMLInputElement)
		if ok {
			inputs = append(inputs, v)
		}
	}
	return inputs
}

func (el Field) Name() string {
	name, _ := getNameFromID(el.div.ID())
	return name
}

// Marshal writes form data to apier data.
func (a *APIer) Marshal() error {
	a.ForEachInput(func(div Input, field reflect.Value) {
		el := div.Input()
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

func (a *APIer) Form() *dom.HTMLFormElement {
	return a.form
}

// ForEachInput iterates over `div` elements containing a `label` and a `input` element.
func (a APIer) ForEachInput(f func(Input, reflect.Value)) {
	doc := dom.GetWindow().Document()
	walkStruct("json", a.dataPtr, a.formID, func(tree string, field reflect.Value) {
		k := field.Kind()
		if !isInputtable(k) {
			return
		}
		id, _ := getIDFromTree(tree)
		f(Input{doc.GetElementByID(id).(*dom.HTMLDivElement)}, field)
	})
}

// ForEachField iterates over `div` elements containing a `label` followed by `div` elements
// representing the inputs.
func (a APIer) ForEachField(f func(Field, reflect.Value)) {
	doc := dom.GetWindow().Document()
	walkStruct("json", a.dataPtr, a.formID, func(tree string, field reflect.Value) {
		k := field.Kind()
		if k == reflect.Ptr { // follow pointer.
			field = reflect.Indirect(field)
			k = field.Kind()
		}
		if k != reflect.Struct {
			return
		}
		id, _ := getIDFromTree(tree)
		f(Field{doc.GetElementByID(id).(*dom.HTMLDivElement)}, field)
	})
}

func (a APIer) newInnerFormHTML() dom.Node {
	maindiv := newSingleHTML("<div></div>")
	doc := dom.GetWindow().Document()
	walkStruct("json", a.dataPtr, a.formID, func(tree string, field reflect.Value) {

		id, parent := getIDFromTree(tree)
		var parentDiv dom.Node
		if parent != "" {
			parentDiv = doc.GetElementByID(parent)
		} else {
			parentDiv = maindiv
		}
		switch field.Kind() {
		case reflect.Struct:
			parentDiv.AppendChild(newSingleHTML(fmt.Sprintf("<div id=\"%s\"><label></label></div>", id)))
		case reflect.Float64, reflect.String, reflect.Int:
			parentDiv.AppendChild(newSingleHTML(fmt.Sprintf("<div id=\"%s\"><input></input><span></span></div>", id)))

		}
	})
	maindiv.AppendChild(newSingleHTML("<div><button type=\"submit\">Submit</button></div>"))
	return maindiv
}

func getIDFromTree(tree string) (string, string) { // return current tree ID and parent's ID
	sp := strings.Split(tree, ".")
	if len(sp) < 2 {
		return strings.Join(sp, "."), ""
	}
	return tree, strings.Join(sp[:len(sp)-1], ".")
}

func isInputtable(v reflect.Kind) bool {
	return v == reflect.Float64 || v == reflect.Int || v == reflect.String
}

func getNameFromID(id string) (name, parentName string) {
	sp := strings.Split(id, ".")
	if len(sp) < 2 {
		return sp[0], ""
	}
	return sp[len(sp)-1], sp[len(sp)-2]
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
