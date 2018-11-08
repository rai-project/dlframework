package predict

import (
	"github.com/rai-project/dlframework"
	"github.com/spf13/cast"
)

func ToFeatures(fs0 interface{}) dlframework.Features {
	fs, err := cast.ToSliceE(fs0)
	if err != nil {
		panic("expecting a list")
	}
	features := make([]*dlframework.Feature, len(f))
	for ii, f := range fs {
		features[ii] = ToFeature(f)
	}
	return features
}

func ToFeature(f interface{}) dlframework.Feature {
	switch f := f.(type) {
	case dlframework.Classification, *dlframework.Classification:
		panic("unhandled classification")
	case dlframework.GeoLocation, *dlframework.GeoLocation:
		panic("unhandled geolocation")
	case dlframework.Region, *dlframework.Region:
		panic("unhandled region")
	case dlframework.Image, *dlframework.Image:
		panic("unhandled image")
	case dlframework.Text, *dlframework.Text:
		panic("unhandled text")
	case dlframework.Audio, *dlframework.Audio:
		panic("unhandled audio")
	}
}
