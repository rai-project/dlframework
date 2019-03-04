package agent

import "github.com/rai-project/dlframework/framework/predictor"

type PredictorLifetime struct {
	Predictor      *predictor.Predictor
	ReferenceCount uint64
}
