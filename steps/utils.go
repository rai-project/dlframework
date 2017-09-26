package steps

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"sync"

	"github.com/facebookgo/stack"
	"github.com/fatih/color"
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

func tryBase64DecodeString(q string) []byte {
	s, err := base64.StdEncoding.DecodeString(q)
	if err != nil {
		return []byte(q)
	}
	return s
}

func tryBase64DecodeBytes(q []byte) []byte {
	s, err := base64.StdEncoding.DecodeString(string(q))
	if err != nil {
		return q
	}
	return s
}

func onPanic(step string) {
	if r := recover(); r != nil {
		var err error
		switch r := r.(type) {
		case error:
			err = r
		default:
			err = fmt.Errorf("%v", r)
		}
		stack := stack.Callers(4)
		log.WithError(err).WithField("step", step).Errorf("[%s] %v\n", color.RedString("PANIC RECOVER"), stack)
	}
}

// Merge different channels in one channel
// https://github.com/tmrts/go-patterns/blob/master/messaging/fan_in.md
func merge(cs ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup

	out := make(chan interface{})

	// Start an send goroutine for each input channel in cs. send
	// copies values from c to out until c is closed, then calls wg.Done.
	send := func(c <-chan interface{}) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go send(c)
	}

	// Start a goroutine to close out once all the send goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
