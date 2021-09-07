package frame

import (
	"reflect"
	"strings"
	"testing"

	"gonum.org/v1/gonum/spatial/r3"
)

var c = struct {
	F   string `json:"f"`
	Vec r3.Vec `json:"vec"`
	M   struct {
		H struct {
			A struct {
				L string `json:"little"`
			} `json:"a"`
		} `json:"had"`
	} `json:"mary"`
}{
	F:   "hello",
	Vec: r3.Vec{X: 1, Y: 2, Z: 3},
	M: struct {
		H struct {
			A struct {
				L string "json:\"little\""
			} "json:\"a\""
		} "json:\"had\""
	}{H: struct {
		A struct {
			L string "json:\"little\""
		} "json:\"a\""
	}{A: struct {
		L string "json:\"little\""
	}{L: "lamb"}}},
}

func TestWalk(t *testing.T) {
	dats := make([]string, 1)
	walkStruct("json", &c, "c", func(tree string, field reflect.Value) {
		t.Logf("tree:%v, field: %v", tree, field)
		dats = append(dats, tree)
	})
}
func TestWalkEdit(t *testing.T) {
	const key = "vec.X"
	walkStruct("json", &c, "obj", func(tree string, field reflect.Value) {
		if strings.HasSuffix(tree, key) {
			t.Logf("modify %s to be 1337", key)
			field.SetFloat(1337)
		}
		t.Logf("fields: %v", strings.Split(tree, "."))
	})
	t.Logf("%+v", c)
}
