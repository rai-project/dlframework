package http

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	webmodels "github.com/rai-project/dlframework/httpapi/models"
	"github.com/rai-project/dlframework/httpapi/restapi/operations/predict"
	"github.com/rai-project/dlframework/registryquery"
	"github.com/rai-project/grpc"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cast"
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

func fromPredictionOptions(opts *webmodels.DlframeworkPredictionOptions) *dl.PredictionOptions {
	if opts == nil {
		opts = &webmodels.DlframeworkPredictionOptions{}
	}

	execOpts := &dl.ExecutionOptions{}
	if opts.ExecutionOptions != nil {
		execOpts = &dl.ExecutionOptions{
			TraceLevel: dl.ExecutionOptions_TraceLevel(
				tracer.LevelFromName(string(opts.ExecutionOptions.TraceLevel)),
			),
			TimeoutInMs: cast.ToInt64(opts.ExecutionOptions.TimeoutInMs),
			DeviceCount: opts.ExecutionOptions.DeviceCount,
			// CPUOptions: opts.ExecutionOptions.CPUOptions,
			// GpuOptions *DlframeworkGPUOptions `json:"gpu_options,omitempty"`
		}
	} else {
		execOpts = &dl.ExecutionOptions{
			TraceLevel: dl.ExecutionOptions_TraceLevel(
				tracer.Config.Level,
			),
		}
	}

	if opts.RequestID == "" {
		opts.RequestID = uuid.NewV4()
	}
	if opts.FeatureLimit == 0 {
		opts.FeatureLimit = 10
	}

	predOpts := &dl.PredictionOptions{
		RequestID:        opts.RequestID,
		FeatureLimit:     opts.,
		BatchSize:        int(opts.BatchSize),
		FeatureLimit:     int(opts.FeatureLimit),
		ExecutionOptions: execOpts,
	}

	return predOpts
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

	var agent *webmodels.DlframeworkAgent
	if params.Body.Options == nil || params.Body.Options.Agent == "" {
		agent = agents[rand.Intn(len(agents))]
	} else {
		for _, a := range agents {
			if a.Host == params.Body.Options.Agent {
				agent = a
				break
			}
		}
		if agent == nil {
			return NewError("Predict/Open",
				errors.Errorf("unable to find agent %v which supports framework=%s:%s model=%s:%s",
					params.Body.Options.Agent, frameworkName, frameworkVersion, modelName, modelVersion,
				))
		}
	}

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
		Options:          fromPredictionOptions(params.Body.Options),
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

func toDlframeworkTraceID(traceId *dl.TraceID) *webmodels.DlframeworkTraceID {
	if traceId == nil {
		return nil
	}
	return &webmodels.DlframeworkTraceID{
		ID: traceId.Id,
	}
}

func toDlframeworkFeaturesResponse(responses []*dl.FeatureResponse) []*webmodels.DlframeworkFeatureResponse {
	resps := make([]*webmodels.DlframeworkFeatureResponse, len(responses))
	for ii, fr := range responses {
		features := make([]*webmodels.DlframeworkFeature, len(fr.Features))
		for jj, f := range fr.Features {
			features[jj] = &webmodels.DlframeworkFeature{
				Index:       cast.ToString(f.Index),
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
			ID:   image.ID,
			Data: image.Data,
			// Preprocessed: image.Preprocessed,
		}
	}

	ret, err := client.Images(ctx,
		&dl.ImagesRequest{
			Predictor: &dl.Predictor{
				ID: predictorID,
			},
			Images:  images,
			Options: fromPredictionOptions(params.Body.Options),
		},
	)

	if err != nil {
		return NewError("Predict/Images", err)
	}

	resps := toDlframeworkFeaturesResponse(ret.Responses)

	return predict.NewImagesOK().
		WithPayload(&webmodels.DlframeworkFeaturesResponse{
			ID:        predictorID,
			TraceID:   toDlframeworkTraceID(ret.TraceId),
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

	ret, err := client.URLs(ctx,
		&dl.URLsRequest{
			Predictor: &dl.Predictor{
				ID: predictorID,
			},
			Urls:    urls,
			Options: fromPredictionOptions(params.Body.Options),
		},
	)

	if err != nil {
		return NewError("Predict/URLs", err)
	}

	resps := toDlframeworkFeaturesResponse(ret.Responses)

	return predict.NewUrlsOK().
		WithPayload(&webmodels.DlframeworkFeaturesResponse{
			ID:        predictorID,
			TraceID:   toDlframeworkTraceID(ret.TraceId),
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
