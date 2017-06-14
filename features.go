package dlframework

import (
	"errors"
	"sort"

	"github.com/mohae/deepcopy"
)

func TopFeatures(features0 []*PredictionFeature, k int) ([]*PredictionFeature, error) {
	features, ok := deepcopy.Copy(features0).([]*PredictionFeature)
	if !ok {
		return nil, errors.New("unable to perform a copy of the features vectors")
	}
	sort.Slice(features, func(ii, jj int) bool {
		return features[ii].Probability < features[jj].Probability
	})
	ll := len(features)
	if ll <= k {
		return features, nil
	}
	return features[ll-k-1 : ll-1], nil
}
