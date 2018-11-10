package agent

import "github.com/rai-project/dlframework/framework/predict"

type PredictorLifetime struct {
	Predictor      *predict.Predictor
	ReferenceCount uint64
}
