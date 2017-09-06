package steps

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// toSliceE casts an interface to a []interface{} type.
// if the input is a scalar then it's promoted to a slice
func toSlice(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch v := i.(type) {
	case string, int, int64, int32, int16, int8, uint, uint64,
		uint32, uint16, uint8, float64, float32, bool:
		return append(s, v), nil
	case []interface{}:
		return append(s, v...), nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, errors.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

// toFloat32Slice casts an interface to a []float32 type.
func toFloat32Slice(i interface{}) ([]float32, error) {
	if i == nil {
		return []float32{}, fmt.Errorf("unable to cast %#v of type %T to []float32", i, i)
	}

	switch v := i.(type) {
	case []float32:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float32, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat32E(s.Index(j).Interface())
			if err != nil {
				return []float32{}, errors.Errorf("unable to cast %#v of type %T to []float32", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float32{}, errors.Errorf("unable to cast %#v of type %T to []float32", i, i)
	}
}
