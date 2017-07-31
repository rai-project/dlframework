package http

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/rai-project/grpc"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predictor"
	context "golang.org/x/net/context"
)

func getBody(s, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
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

	conn, err := grpc.Dial(dl.PredictorServiceDescription, serverAddress)
	if err != nil {
		return NewError("Predictor", errors.Wrapf(err, "unable to dial %s", serverAddress))
	}

	defer conn.Close()

	client := dl.NewPredictorClient(conn)

	ctx := context.Background()

	data, err := params.Body.Data.MarshalText()
	if err != nil {
		return NewError("Predictor", errors.Wrapf(err, "unable marshal data"))
	}

	resp, err := client.Predict(ctx, &dl.PredictRequest{
		ModelName:        modelName,
		ModelVersion:     modelVersion,
		FrameworkName:    frameworkName,
		FrameworkVersion: frameworkVersion,
		Limit:            10,
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
