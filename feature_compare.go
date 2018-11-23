//go:generate go get github.com/alvaroloes/enumer
//go:generate enumer -type=FeatureCompareMethod -json

package dlframework

type FeatureCompareMethod int

const (
	FeatureCompuareAutomatic FeatureCompareMethod = iota
	FeatureCompareTextEditDistance
	FeatureCompareBoundingBoxIOU
)
