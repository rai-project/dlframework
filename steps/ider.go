package steps

type IDer interface {
	GetId() string
	GetData() interface{}
}

type IDWrapper struct {
	ID   string
	data interface{}
}

func (w IDWrapper) GetId() string {
	return w.ID
}

func (w IDWrapper) GetData() interface{} {
	return w.data
}

func NewIDWrapper(id string, data interface{}) IDer {
	return &IDWrapper{
		ID:   id,
		data: data,
	}
}
