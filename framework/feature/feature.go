package feature

import (
	"github.com/rai-project/dlframework"
)

func New(opts ...Option) *dlframework.Feature {
	feature := &dlframework.Feature{
		Type:     dlframework.FeatureType_UNKNOWN,
		Metadata: map[string]string{},
	}
	for _, o := range opts {
		o(feature)
	}
	return feature
}
