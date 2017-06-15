package types

import "github.com/chewxy/hm"

type Fresh struct {
	count int
}

const letters = `abcdefghijklmnopqrstuvwxyz`

func NewFresh() *Fresh {
	return &Fresh{
		count: 0,
	}
}

func (f *Fresh) Fresh() hm.TypeVariable {
	defer func() {
		f.count++
	}()
	lettersLen := len(letters)
	c := f.count
	retVal := letters[c%lettersLen]
	for c >= 0 {
		retVal = retVal + letters[c%lettersLen]
		c = c / lettersLen
	}
	return hm.TypeVariable(retVal)
}
