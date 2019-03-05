package feature

import "github.com/rai-project/dlframework"

func BoundingBoxType() Option {
	return Type(dlframework.FeatureType_BOUNDINGBOX)
}

func BoundingBox(e *dlframework.BoundingBox) Option {
	return func(o *dlframework.Feature) {
		BoundingBoxType()(o)
		o.Feature = &dlframework.Feature_BoundingBox{
			BoundingBox: e,
		}
	}
}

func ensureBoundingBox(o *dlframework.Feature) *dlframework.BoundingBox {
	if o.Type != dlframework.FeatureType_BOUNDINGBOX && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_BoundingBox{}
	}
	bbox, ok := o.Feature.(*dlframework.Feature_BoundingBox)
	if !ok {
		panic("expecting a classification feature")
	}
	if bbox.BoundingBox == nil {
		bbox.BoundingBox = &dlframework.BoundingBox{}
	}
	BoundingBoxType()(o)
	return bbox.BoundingBox
}

func BoundingBoxXmin(xmin float32) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Xmin = xmin
	}
}

func BoundingBoxXmax(xmax float32) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Xmax = xmax
	}
}

func BoundingBoxYmin(ymin float32) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Ymin = ymin
	}
}

func BoundingBoxYmax(ymax float32) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Ymax = ymax
	}
}
func BoundingBoxIndex(index int32) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Index = index
	}
}

func BoundingBoxLabel(label string) Option {
	return func(o *dlframework.Feature) {
		bbox := ensureBoundingBox(o)
		bbox.Label = label
	}
}
