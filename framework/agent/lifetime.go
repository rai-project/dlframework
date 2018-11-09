package agent

type PredictorLifetime struct {
	Predictor      *Predictor
	ReferenceCount uint64
}
