package feature

import "github.com/rai-project/dlframework"

func InstanceSegmentType() Option {
	return Type(dlframework.FeatureType_INSTANCESEGMENT)
}

func InstanceSegment(e *dlframework.InstanceSegment) Option {
	return func(o *dlframework.Feature) {
		InstanceSegmentType()(o)
		o.Feature = &dlframework.Feature_InstanceSegment{
			InstanceSegment: e,
		}
	}
}

func ensureInstanceSegment(o *dlframework.Feature) *dlframework.InstanceSegment {
	if o.Type != dlframework.FeatureType_INSTANCESEGMENT && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_InstanceSegment{}
	}
	iseg, ok := o.Feature.(*dlframework.Feature_InstanceSegment)
	if !ok {
		panic("expecting a classification feature")
	}
	if iseg.InstanceSegment == nil {
		iseg.InstanceSegment = &dlframework.InstanceSegment{}
	}
	InstanceSegmentType()(o)
	return iseg.InstanceSegment
}

func InstanceSegmentXmin(xmin float32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Xmin = xmin
	}
}

func InstanceSegmentXmax(xmax float32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Xmax = xmax
	}
}

func InstanceSegmentYmin(ymin float32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Ymin = ymin
	}
}

func InstanceSegmentYmax(ymax float32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Ymax = ymax
	}
}
func InstanceSegmentIndex(index int32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Index = index
	}
}

func InstanceSegmentLabel(label string) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Label = label
	}
}

func InstanceSegmentMaskType(masktype string) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.MaskType = masktype
	}
}

func InstanceSegmentHeight(height int32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Height = height
	}
}

func InstanceSegmentWidth(width int32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.Width = width
	}
}

func InstanceSegmentIntMask(mask []int32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.IntMask = mask
	}
}

func InstanceSegmentFloatMask(mask []float32) Option {
	return func(o *dlframework.Feature) {
		iseg := ensureInstanceSegment(o)
		iseg.FloatMask = mask
	}
}
