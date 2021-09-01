package frame

import (
	"errors"
	"fmt"
	"math"
	"reflect"

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
func getBranch(tagname string, ptr interface{}) []branch {
	main := make([]branch, 0)
	v := reflect.ValueOf(ptr)
	if v.IsNil() {
		panic("nil value passed to getBranch")
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		panic("getBranch called on non-struct interface")
	}
	a := v.Type()
	for i := 0; i < a.NumField(); i++ {
		field := a.Field(i)
		tag := field.Tag.Get(tagname)
		if tag == "" {
			continue
		}
		var b branch
		switch field.Type.Kind() {
		case reflect.Struct:
			b = branch{name: tag, number: math.NaN(), branches: getBranch(tagname, v.Field(i))}
		case reflect.Bool:
			b = branch{name: tag, number: math.NaN()}
		case reflect.String:
			b = branch{name: tag, number: math.NaN(), str: v.Field(i).String()}
		case reflect.Int:
			b = branch{name: tag, number: float64(v.Field(i).Int())}
		case reflect.Float64:
			b = branch{name: tag, number: v.Field(i).Float()}
		default:
			continue
		}
		main = append(main, b)
	}
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
