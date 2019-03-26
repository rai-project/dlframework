package feature

import "github.com/rai-project/dlframework"

func SemanticSegmentType() Option {
	return Type(dlframework.FeatureType_SEMANTICSEGMENT)
}

func SemanticSegment(e *dlframework.SemanticSegment) Option {
	return func(o *dlframework.Feature) {
		SemanticSegmentType()(o)
		o.Feature = &dlframework.Feature_SemanticSegment{
			SemanticSegment: e,
		}
	}
}

func ensureSemanticSegment(o *dlframework.Feature) *dlframework.SemanticSegment {
	if o.Type != dlframework.FeatureType_SEMANTICSEGMENT && !isUnknownType(o) {
		panic("unexpected feature type")
	}
	if o.Feature == nil {
		o.Feature = &dlframework.Feature_SemanticSegment{}
	}
	sseg, ok := o.Feature.(*dlframework.Feature_SemanticSegment)
	if !ok {
		panic("expecting a SemanticSegment feature")
	}
	if sseg.SemanticSegment == nil {
		sseg.SemanticSegment = &dlframework.SemanticSegment{}
	}
	SemanticSegmentType()(o)
	return sseg.SemanticSegment
}

func SemanticSegmentHeight(height int32) Option {
	return func(o *dlframework.Feature) {
		sseg := ensureSemanticSegment(o)
		sseg.Height = height
	}
}

func SemanticSegmentWidth(width int32) Option {
	return func(o *dlframework.Feature) {
		sseg := ensureSemanticSegment(o)
		sseg.Width = width
	}
}

func SemanticSegmentIntMask(mask []int32) Option {
	return func(o *dlframework.Feature) {
		sseg := ensureSemanticSegment(o)
		sseg.IntMask = mask
	}
}
