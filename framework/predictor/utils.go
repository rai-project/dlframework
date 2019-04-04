package predictor

import (
	"reflect"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func flattenFloat32Slice(data interface{}) []float32 {
	// when data is a scalar
	rval := reflect.ValueOf(data)
	typ := rval.Type()
	if typ.Kind() != reflect.Array &&
		typ.Kind() != reflect.Slice &&
		typ.Kind() != reflect.Interface {
		ddata := indirect(data)
		if e, ok := ddata.(float32); ok {
			return []float32{e}
		}
		switch s := ddata.(type) {
		case bool:
			if s {
				return []float32{float32(1)}
			}
			return []float32{float32(0)}
		case int:
			return []float32{float32(s)}
		case uint:
			return []float32{float32(s)}
		case int8:
			return []float32{float32(s)}
		case uint8:
			return []float32{float32(s)}
		case int16:
			return []float32{float32(s)}
		case uint16:
			return []float32{float32(s)}
		case int32:
			return []float32{float32(s)}
		case uint32:
			return []float32{float32(s)}
		case int64:
			return []float32{float32(s)}
		case uint64:
			return []float32{float32(s)}
		case float32:
			return []float32{float32(s)}
		case float64:
			return []float32{float32(s)}
		case uintptr:
			return []float32{float32(s)}
		}
		panic(errors.Errorf("unable to convert %v of kind %v", pp.Sprint(data), typ.Kind().String()))
	}

	// no we know data is a slice
	res := []float32{}
	for ii := 0; ii < rval.Len(); ii++ {
		val := rval.Index(ii)
		fval := flattenFloat32Slice(val.Interface())
		res = append(res, fval...)
	}
	return res
}

func flattenInt32Slice(data interface{}) []int32 {
	// when data is a scalar
	rval := reflect.ValueOf(data)
	typ := rval.Type()
	if typ.Kind() != reflect.Array &&
		typ.Kind() != reflect.Slice &&
		typ.Kind() != reflect.Interface {
		ddata := indirect(data)
		if e, ok := ddata.(int32); ok {
			return []int32{e}
		}
		switch s := ddata.(type) {
		case bool:
			if s {
				return []int32{int32(1)}
			}
			return []int32{int32(0)}
		case int:
			return []int32{int32(s)}
		case uint:
			return []int32{int32(s)}
		case int8:
			return []int32{int32(s)}
		case uint8:
			return []int32{int32(s)}
		case int16:
			return []int32{int32(s)}
		case uint16:
			return []int32{int32(s)}
		case int32:
			return []int32{int32(s)}
		case uint32:
			return []int32{int32(s)}
		case int64:
			return []int32{int32(s)}
		case uint64:
			return []int32{int32(s)}
		case float32:
			return []int32{int32(s)}
		case float64:
			return []int32{int32(s)}
		case uintptr:
			return []int32{int32(s)}
		}
		panic(errors.Errorf("unable to convert %v of kind %v", pp.Sprint(data), typ.Kind().String()))
	}

	// no we know data is a slice
	res := []int32{}
	for ii := 0; ii < rval.Len(); ii++ {
		val := rval.Index(ii)
		fval := flattenInt32Slice(val.Interface())
		res = append(res, fval...)
	}
	return res
}
