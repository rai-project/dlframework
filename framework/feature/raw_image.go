package feature

import "github.com/rai-project/dlframework"

func RawImageType() Option {
	return Type(dlframework.FeatureType_RAW_IMAGE)
}

func RawImage(e *dlframework.RawImage) Option {
	return func(o *dlframework.Feature) {
		RawImageType()(o)
		o.Feature = &dlframework.Feature_RawImage{
			RawImage: e,
		}
	}
}

func ensureRawImage(o *dlframework.Feature) *dlframework.RawImage {
	if o.Type != dlframework.FeatureType_RAW_IMAGE && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_Image{}
	}
	img, ok := o.Feature.(*dlframework.Feature_RawImage)
	if !ok {
		panic("expecting an raw image feature")
	}
	if img.RawImage == nil {
		img.RawImage = &dlframework.RawImage{}
	}
	RawImageType()(o)
	return img.RawImage
}

func RawImageID(id string) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		img.ID = id
	}
}

func RawImageWidth(width int) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		img.Width = int32(width)
	}
}

func RawImageHeight(height int) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		img.Height = int32(height)
	}
}

func RawImageChannels(channels int) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		img.Channels = int32(channels)
	}
}

func RawImageFloatData(data []float32) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		img.FloatList = data
	}
}

func RawImageInt8Data(data []int8) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		buf := make([]int32, len(data))
		for ii, val := range data {
			buf[ii] = int32(val)
		}
		img.CharList = buf
	}
}

func RawImageUInt8Data(data []uint8) Option {
	return func(o *dlframework.Feature) {
		img := ensureRawImage(o)
		buf := make([]int32, len(data))
		for ii, val := range data {
			buf[ii] = int32(val)
		}
		img.CharList = buf
	}
}

func RawImageData(data interface{}) Option {
	switch v := data.(type) {
	case []int8:
		return RawImageInt8Data(v)
	case []uint8:
		return RawImageUInt8Data(v)
	case []float32:
		return RawImageFloatData(v)
	}
	panic("invalid RawImageData type")
}
