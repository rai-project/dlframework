package http

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predict"
	"github.com/rai-project/dlframework/registryquery"
	"github.com/rai-project/grpc"
	"golang.org/x/sync/syncmap"
	gogrpc "google.golang.org/grpc"
)

type PredictHandler struct {
	clients     syncmap.Map
	connections syncmap.Map
}

func getBody(s, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

func (p *PredictHandler) Open(params predict.OpenParams) middleware.Responder {

	frameworkName := strings.ToLower(getBody(params.Body.FrameworkName, "*"))
	frameworkVersion := strings.ToLower(getBody(params.Body.FrameworkVersion, "*"))
	modelName := strings.ToLower(getBody(params.Body.ModelName, "*"))
	modelVersion := strings.ToLower(getBody(params.Body.ModelVersion, "*"))

	agents, err := registryquery.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
	if err != nil {
		return NewError("Predict/Open", err)
	}

	if len(agents) == 0 {
		return NewError("Predict/Open",
			errors.Errorf("unable to find agents for framework=%s:%s model=%s:%s",
				frameworkName, frameworkVersion, modelName, modelVersion,
			))
	}

	agent := agents[rand.Intn(len(agents))]
	serverAddress := fmt.Sprintf("%s:%s", agent.Host, agent.Port)

	ctx := params.HTTPRequest.Context()
	conn, err := grpc.DialContext(ctx, dl.PredictServiceDescription, serverAddress)
	if err != nil {
		return NewError("Predict/Open", errors.Wrapf(err, "unable to dial %s", serverAddress))
	}

	client := dl.NewPredictClient(conn)

	predictor, err := client.Open(ctx, &dl.PredictorOpenRequest{
		ModelName:        modelName,
		ModelVersion:     modelVersion,
		FrameworkName:    frameworkName,
		FrameworkVersion: frameworkVersion,
	})

	if err != nil {
		defer conn.Close()
		return NewError("Predict/Open", errors.Wrap(err, "unable to open model"))
	}

	p.clients.Store(predictor.ID, client)
	p.connections.Store(predictor.ID, conn)

	return predict.NewOpenOK().WithPayload(&webmodels.DlframeworkPredictor{
		ID: predictor.ID,
	})
}

func (p *PredictHandler) getClient(id string) (dl.PredictClient, error) {
	val, ok := p.clients.Load(id)
	if !ok {
		return nil, errors.New("unable to get client predictor value")
	}
	client, ok := val.(dl.PredictClient)
	if !ok {
		return nil, errors.New("unable to get client predictor connection")
	}
	return client, nil
}

func (p *PredictHandler) getConnection(id string) (*gogrpc.ClientConn, error) {
	val, ok := p.connections.Load(id)
	if !ok {
		return nil, errors.New("unable to get connection predictor value")
	}
	conn, ok := val.(*gogrpc.ClientConn)
	if !ok {
		return nil, errors.New("unable to get connection predictor connection")
	}
	return conn, nil
}

func (p *PredictHandler) closeClient(ctx context.Context, id string) error {
	client, err := p.getClient(id)
	if err != nil {
		return err
	}
	if _, err := client.Close(ctx, &dl.Predictor{ID: id}); err != nil {
		return err
	}
	return nil
}

func (p *PredictHandler) closeConnection(ctx context.Context, id string) error {
	conn, err := p.getConnection(id)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (p *PredictHandler) Close(params predict.CloseParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	id := params.Body.ID

	if err := p.closeClient(ctx, id); err != nil {
		defer p.closeConnection(ctx, id)
		return NewError("Predict/Close", errors.Wrap(err, "failed to close predictor client"))
	}

	if err := p.closeConnection(ctx, id); err != nil {
		return NewError("Predict/Close", errors.Wrap(err, "failed to close grpc connection"))
	}

	var resp webmodels.DlframeworkPredictorCloseResponse
	return predict.NewCloseOK().WithPayload(resp)
}

func (p *PredictHandler) Reset(params predict.ResetParams) middleware.Responder {
	predictorId := params.Body.Predictor.ID

	client, err := p.getClient(predictorId)
	if err != nil {
		return NewError("Predict/Reset", err)
	}

	ctx := params.HTTPRequest.Context()
	id := params.Body.ID

	resp, err := client.Reset(ctx,
		&dl.ResetRequest{
			ID: id,
			Predictor: &dl.Predictor{
				ID: predictorId,
			},
		},
	)
	if err != nil {
		return NewError("Predict/Reset", err)
	}

	return predict.NewResetOK().
		WithPayload(&webmodels.DlframeworkResetResponse{
			Predictor: &webmodels.DlframeworkPredictor{ID: resp.Predictor.ID},
		})
}

func toDlframeworkFeaturesResponse(responses []*dl.FeatureResponse) []*webmodels.DlframeworkFeatureResponse {
	resps := make([]*webmodels.DlframeworkFeatureResponse, len(responses))
	for ii, fr := range responses {
		features := make([]*webmodels.DlframeworkFeature, len(fr.Features))
		for jj, f := range fr.Features {
			features[jj] = &webmodels.DlframeworkFeature{
				Index:       f.Index,
				Metadata:    f.Metadata,
				Name:        f.Name,
				Probability: f.Probability,
			}
		}
		resps[ii] = &webmodels.DlframeworkFeatureResponse{
			Features:  features,
			ID:        fr.ID,
			InputID:   fr.InputID,
			Metadata:  fr.Metadata,
			RequestID: fr.RequestID,
		}
	}
	return resps
}

func (p *PredictHandler) Images(params predict.ImagesParams) middleware.Responder {
	predictor := params.Body.Predictor
	if predictor == nil {
		return NewError("Predict/Images", errors.New("invalid nil predictor"))
	}
	predictorID := predictor.ID

	client, err := p.getClient(predictorID)
	if err != nil {
		return NewError("Predict/Images", err)
	}

	ctx := params.HTTPRequest.Context()

	images := make([]*dl.ImagesRequest_Image, len(params.Body.Images))
	for ii, image := range params.Body.Images {
		images[ii] = &dl.ImagesRequest_Image{
			ID:           image.ID,
			Data:         image.Data,
			Preprocessed: image.Preprocessed,
		}
	}

	options := params.Body.Options
	if options == nil {
		options = &webmodels.DlframeworkPredictionOptions{}
	}

	requestID := params.HTTPRequest.Header.Get(echo.HeaderXRequestID)
	if options.RequestID != "" {
		requestID = options.RequestID
	}

	ret, err := client.Images(ctx,
		&dl.ImagesRequest{
			Predictor: &dl.Predictor{
				ID: predictorID,
			},
			Images: images,
			Options: &dl.PredictionOptions{
				RequestID:    requestID,
				FeatureLimit: options.FeatureLimit,
			},
		},
	)

	if err != nil {
		return NewError("Predict/Images", err)
	}

	resps := toDlframeworkFeaturesResponse(ret.Responses)

	return predict.NewImagesOK().
		WithPayload(&webmodels.DlframeworkFeaturesResponse{
			ID:        predictorID,
			Responses: resps,
		})
}

func (p *PredictHandler) URLs(params predict.UrlsParams) middleware.Responder {
	predictor := params.Body.Predictor
	if predictor == nil {
		return NewError("Predict/URLs", errors.New("invalid nil predictor"))
	}
	predictorID := predictor.ID

	client, err := p.getClient(predictorID)
	if err != nil {
		return NewError("Predict/URLs", err)
	}

	ctx := params.HTTPRequest.Context()

	urls := make([]*dl.URLsRequest_URL, len(params.Body.Urls))
	for ii, url := range params.Body.Urls {
		urls[ii] = &dl.URLsRequest_URL{
			ID:   url.ID,
			Data: url.Data,
		}
	}

	options := params.Body.Options
	if options == nil {
		options = &webmodels.DlframeworkPredictionOptions{}
	}

	requestID := params.HTTPRequest.Header.Get(echo.HeaderXRequestID)
	if options.RequestID != "" {
		requestID = options.RequestID
	}

	ret, err := client.URLs(ctx,
		&dl.URLsRequest{
			Predictor: &dl.Predictor{
				ID: predictorID,
			},
			Urls: urls,
			Options: &dl.PredictionOptions{
				RequestID:    requestID,
				FeatureLimit: options.FeatureLimit,
			},
		},
	)

	if err != nil {
		return NewError("Predict/URLs", err)
	}

	resps := toDlframeworkFeaturesResponse(ret.Responses)

	return predict.NewUrlsOK().
		WithPayload(&webmodels.DlframeworkFeaturesResponse{
			ID:        predictorID,
			Responses: resps,
		})
}

func (p *PredictHandler) Dataset(params predict.DatasetParams) middleware.Responder {
	return middleware.NotImplemented("operation predict.Dataset has not yet been implemented")
}

// func PredictorPredictHandler(params predictor.PredictParams) middleware.Responder {

// 	frameworkName := strings.ToLower(getBody(params.Body.FrameworkName, "*"))
// 	frameworkVersion := strings.ToLower(getBody(params.Body.FrameworkVersion, "*"))
// 	modelName := strings.ToLower(getBody(params.Body.ModelName, "*"))
// 	modelVersion := strings.ToLower(getBody(params.Body.ModelVersion, "*"))

// 	agents, err := Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
// 	if err != nil {
// 		return NewError("Predictor", err)
// 	}

// 	if len(agents) == 0 {
// 		return NewError("Predictor",
// 			errors.Errorf("unable to find agents for framework=%s:%s model=%s:%s",
// 				frameworkName, frameworkVersion, modelName, modelVersion,
// 			))
// 	}

// 	agent := agents[rand.Intn(len(agents))]
// 	serverAddress := fmt.Sprintf("%s:%s", agent.Host, agent.Port)

// 	ctx := params.HTTPRequest.Context()
// 	conn, err := grpc.DialContext(ctx, dl.PredictorServiceDescription, serverAddress)
// 	if err != nil {
// 		return NewError("Predictor", errors.Wrapf(err, "unable to dial %s", serverAddress))
// 	}

// 	defer conn.Close()

// 	client := dl.NewPredictorClient(conn)

// 	data, err := params.Body.Data.MarshalText()
// 	if err != nil {
// 		return NewError("Predictor", errors.Wrapf(err, "unable marshal data"))
// 	}

// 	resp, err := client.Predict(ctx, &dl.PredictRequest{
// 		ModelName:        modelName,
// 		ModelVersion:     modelVersion,
// 		FrameworkName:    frameworkName,
// 		FrameworkVersion: frameworkVersion,
// 		Limit:            params.Body.Limit,
// 		Data:             data,
// 	})

// 	if err != nil {
// 		return NewError("Predictor", errors.Wrap(err, "unable to predict model"))
// 	}

// 	res := new(webmodels.DlframeworkPredictResponse)
// 	if err := copier.Copy(res, resp); err != nil {
// 		return NewError("Predictor", errors.Wrap(err, "unable to copy predict response to webmodels"))
// 	}

// 	return predictor.NewPredictOK().WithPayload(res)
// }
