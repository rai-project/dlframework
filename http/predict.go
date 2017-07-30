package http

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predictor"
)

func PredictorPredictHandler(params predictor.PredictParams) middleware.Responder {
	return middleware.NotImplemented("operation predictor.Predict has not yet been implemented")
}
