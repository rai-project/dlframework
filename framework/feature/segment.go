package feature

import "github.com/rai-project/dlframework"

func SegmentType() Option {
	return Type(dlframework.FeatureType_SEGMENT)
}

func Segment(e *dlframework.Segment) Option {
	return func(o *dlframework.Feature) {
		SegmentType()(o)
		o.Feature = &dlframework.Feature_Segment{
			Segment: e,
		}
	}
}

func ensureSegment(o *dlframework.Feature) *dlframework.Segment {
	if o.Type != dlframework.FeatureType_SEGMENT && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_Segment{}
	}
	seg, ok := o.Feature.(*dlframework.Feature_Segment)
	if !ok {
		panic("expecting a classification feature")
	}
	if seg.Segment == nil {
		seg.Segment = &dlframework.Segment{}
	}
	SegmentType()(o)
	return seg.Segment
}

func SegmentIndex(index int32) Option {
	return func(o *dlframework.Feature) {
		seg := ensureSegment(o)
		seg.Index = index
	}
}

func SegmentLabel(label string) Option {
	return func(o *dlframework.Feature) {
		seg := ensureSegment(o)
		seg.Label = label
	}
}

func SegmentData(data []byte) Option {
	return func(o *dlframework.Feature) {
		seg := ensureSegment(o)
		seg.Data = data
	}
}

func SegmentHeight(height int32) Option {
	return func(o *dlframework.Feature) {
		seg := ensureSegment(o)
		seg.Height = height
	}
}

func SegmentWidth(width int32) Option {
	return func(o *dlframework.Feature) {
		seg := ensureSegment(o)
		seg.Width = width
	}
}
