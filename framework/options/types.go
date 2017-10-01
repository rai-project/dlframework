package options

type device struct {
	id         int        // device id
	deviceType DeviceType // device type
}

type inputNode struct {
	key   string // name
	shape []uint // shape of ndarray
}

type DeviceType int

const (
	CPU_DEVICE DeviceType = iota + 1 // cpu device type
	GPU_DEVICE                       // gpu device type
)
