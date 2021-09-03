package frame

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"gonum.org/v1/gonum/mat"
)

const symTag = "symbol"

// Converts a struct with float64 field `symbol` tags to gonum's Vector representation
func SymToVec(a interface{}) *mat.VecDense {
	return tovec(getBranch(symTag, reflect.ValueOf(a)))
}

// MarshalSym populates a pointer to a `symbol` filled struct with float64 values.
func MarshalSym(vec []float64, a interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(a))
	if v.Kind() != reflect.Struct {
		return errors.New("MarshalSym called on non-struct interface")
	}
	n := len(vec)
	counts := 0
	e := v
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		fieldT := tp.Field(i)
		tag := fieldT.Tag.Get(symTag)
		if tag == "" {
			continue
		}
		switch fieldT.Type.Kind() {
		case reflect.Float64:
			if counts > n-1 {
				return errors.New("MarshalSym got struct with less symbol fields than float slice length")
			}
			e.Field(i).SetFloat(vec[counts])
			counts++
		}
	}
	return nil
}

func SymIdx(syms, target interface{}) []int {
	b := getBranch(symTag, reflect.ValueOf(syms))
	t := getBranch(symTag, reflect.ValueOf(target))
	idxs := make([]int, len(b))
	for i := range b {
		idxs[i] = -1
		for j := range t {
			if b[i].name == t[j].name {
				idxs[i] = j
				break
			}
		}
	}
	return idxs
}

// data structure to store reduced struct data.
type branch struct {
	name     string
	number   float64
	str      string
	branches []branch
}

func (b branch) Kind() reflect.Kind {
	switch {
	case b.number == b.number:
		return reflect.Float64
	case len(b.branches) > 0:
		return reflect.Struct
	case b.str != "" || b.name != "": // this is crap
		return reflect.String
	}
	return reflect.Invalid
}

func (b branch) String() string {
	var i interface{}
	switch b.Kind() {
	case reflect.Float64:
		i = b.number
	case reflect.String:
		i = b.str
	}
	return fmt.Sprintf("%v", i)
}

func walkStruct(tagname string, ptr interface{}, tree string, f func(tree string, field reflect.Value)) {
	v := reflect.ValueOf(ptr)
	vderef := reflect.Indirect(v)
	if v.Type().Kind() != reflect.Struct && (v.IsNil() || vderef == v) {
		panic("nil value passed to getBranch or not a pointer")
	}
	if vderef.Kind() != reflect.Struct {
		panic("getBranch called on non-struct interface:" + v.Kind().String())
	}
	a := vderef.Type()
	for i := 0; i < a.NumField(); i++ {
		fieldT := a.Field(i)
		tag := fieldT.Tag.Get(tagname)
		switch tag {
		case "-":
			continue
		case "":
			tag = fieldT.Name
		}
		fieldV := vderef.Field(i)
		switch fieldT.Type.Kind() {
		case reflect.Struct:
			if reflect.Indirect(fieldV) == fieldV {
				// is not a pointer
				walkStruct(tagname, fieldV.Addr().Interface(), tree+"."+tag, f)
			} else {
				// is a pointer to a struct
				walkStruct(tagname, fieldV.Interface(), tree+"."+tag, f)
			}
		default:
			f(tree+"."+tag, fieldV)
		}
	}
}

func getBranch(tagname string, ptr interface{}) []branch {
	main := make([]branch, 0)
	walkStruct(tagname, ptr, "", func(tree string, field reflect.Value) {
		var b branch
		tags := strings.Split(tree, ".")
		if len(tags) == 0 {
			panic("zero length tags")
		}
		tag := tags[len(tags)-1]
		switch field.Kind() {
		case reflect.Struct:
			b = branch{name: tag, number: math.NaN(), branches: getBranch(tagname, field)}
		case reflect.Bool:
			b = branch{name: tag, number: math.NaN()}
		case reflect.String:
			b = branch{name: tag, number: math.NaN(), str: field.String()}
		case reflect.Int:
			b = branch{name: tag, number: float64(field.Int())}
		case reflect.Float64:
			b = branch{name: tag, number: field.Float()}
		default:
			// do not add
		}
		main = append(main, b)
	})

	return main
}

func tovec(branches []branch) *mat.VecDense {
	n := len(branches)
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		data[i] = branches[i].number
	}
	return mat.NewVecDense(n, data) //mat.NewDense(n, 1, data)
}
