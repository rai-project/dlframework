package steps

import (
	"strings"

	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/pipeline"
)

func NewPredict(predictor predictor.Predictor) pipeline.Step {
	_, model, err := predictor.Info()
	if err != nil {
		panic(err)
	}

	outputType := model.GetOutput().GetType()

	switch strings.ToLower(outputType) {
	case "classification":
		return NewPredictImageClassification(predictor)
	case "boundingbox":
		return NewPredictObjectDetection(predictor)
	default:
		return NewPredictImageClassification(predictor)
	}
}
