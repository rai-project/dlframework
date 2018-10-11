package options

import (
	"bytes"
	"fmt"
)

type DeviceType int

const (
	CPU_DEVICE  DeviceType = iota + 1 // cpu device type
	CUDA_DEVICE DeviceType = 2        // cuda device type
)

func (d DeviceType) String() string {
	switch d {
	case CPU_DEVICE:
		return "cpu"
	case CUDA_DEVICE:
		return "cuda"
	}
	return "<<unknown_device>>"
}

type device struct {
	id         int        // device id
	deviceType DeviceType // device type
}

func (n device) ID() int {
	return n.id
}

func (n device) Type() DeviceType {
	return n.deviceType
}

func (n device) String() string {
	return fmt.Sprintf("%v:%v", n.deviceType, n.id)
}

type devices []device

func (d devices) String() string {
	if len(d) == 0 {
		return "[]"
	}
	buf := new(bytes.Buffer)
	buf.WriteString("[")
	for _, n := range d {
		buf.WriteString(n.String())
		buf.WriteString(",")
	}
	r := buf.Bytes()
	r[len(r)-1] = ']'
	return string(r)
}

type inputNode struct {
	key   string // name
	shape []int  // shape of ndarray
}

func (n inputNode) Key() string {
	return n.key
}

func (n inputNode) Shape() []int {
	return n.shape
}
