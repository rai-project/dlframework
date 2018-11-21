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
