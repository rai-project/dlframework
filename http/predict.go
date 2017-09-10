package http

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/httpapi/client/predict"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predictor"
	"github.com/rai-project/grpc"
)

type PredictHandler struct{
  agents syncmap.Map
}

func getBody(s, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

func (p *PredictHandler) Open(params predict.OpenParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.Open has not yet been implemented")
}

func (p *PredictHandler) Close(params predict.CloseParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.CloseHandler has not yet been implemented")
}

func (p *PredictHandler) Reset(params predict.ResetParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.Reset has not yet been implemented")
}

func (p *PredictHandler) getAgent(params predict.ImagesParams)  (*webmodels.DlframeworkAgent, error) {
  return nil, nil
}

func (p *PredictHandler) Images(params predict.ImagesParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.Images has not yet been implemented")
})


func (p *PredictHandler) URLs(func(params predict.UrlsParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.Urls has not yet been implemented")
}

func (p *PredictHandler) Dataset(params predict.DatasetParams) middleware.Responder {
  return middleware.NotImplemented("operation predict.Dataset has not yet been implemented")
}

func PredictorPredictHandler(params predictor.PredictParams) middleware.Responder {

	frameworkName := strings.ToLower(getBody(params.Body.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getBody(params.Body.FrameworkVersion, "*"))
	modelName := strings.ToLower(getBody(params.Body.ModelName, "*"))
	modelVersion := strings.ToLower(getBody(params.Body.ModelVersion, "*"))

	agents, err := models.agents(frameworkName, frameworkVersion, modelName, modelVersion)
	if err != nil {
		return NewError("Predictor", err)
	}

	if len(agents) == 0 {
		return NewError("Predictor",
			errors.Errorf("unable to find agents for framework=%s:%s model=%s:%s",
				frameworkName, frameworkVersion, modelName, modelVersion,
			))
	}

	agent := agents[rand.Intn(len(agents))]
	serverAddress := fmt.Sprintf("%s:%s", agent.Host, agent.Port)

	ctx := params.HTTPRequest.Context()
	conn, err := grpc.DialContext(ctx, dl.PredictorServiceDescription, serverAddress)
	if err != nil {
		return NewError("Predictor", errors.Wrapf(err, "unable to dial %s", serverAddress))
	}

	defer conn.Close()

	client := dl.NewPredictorClient(conn)

	data, err := params.Body.Data.MarshalText()
	if err != nil {
		return NewError("Predictor", errors.Wrapf(err, "unable marshal data"))
	}

	resp, err := client.Predict(ctx, &dl.PredictRequest{
		ModelName:        modelName,
		ModelVersion:     modelVersion,
		FrameworkName:    frameworkName,
		FrameworkVersion: frameworkVersion,
		Limit:            params.Body.Limit,
		Data:             data,
	})

	if err != nil {
		return NewError("Predictor", errors.Wrap(err, "unable to predict model"))
	}

	res := new(webmodels.DlframeworkPredictResponse)
	if err := copier.Copy(res, resp); err != nil {
		return NewError("Predictor", errors.Wrap(err, "unable to copy predict response to webmodels"))
	}

	return predictor.NewPredictOK().WithPayload(res)
}
