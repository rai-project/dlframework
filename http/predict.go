package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/snappy"

	goimage "image"
	"image/color"
	"image/jpeg"

	"github.com/cenkalti/backoff"
	"github.com/go-openapi/runtime/middleware"
	"github.com/k0kubun/pp"
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

const compressRawImage = false

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
			TimeoutInMs: cast.ToUint64(opts.ExecutionOptions.TimeoutInMs),
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
		BatchSize:        cast.ToInt32(opts.BatchSize),
		FeatureLimit:     cast.ToInt32(opts.FeatureLimit),
		ExecutionOptions: execOpts,
	}

	return predOpts
}

var findMaxTries uint64 = 10

func (p *PredictHandler) findAgent(params predict.OpenParams) (agents []*webmodels.DlframeworkAgent, err error) {

	frameworkName := dl.CleanString(getBody(params.Body.FrameworkName, "*"))
	frameworkVersion := dl.CleanString(getBody(params.Body.FrameworkVersion, "*"))
	modelName := dl.CleanString(getBody(params.Body.ModelName, "*"))
	modelVersion := dl.CleanString(getBody(params.Body.ModelVersion, "*"))

	find := func() error {
		agents, err = registryquery.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
		if err != nil {
			return NewError("Predict/FindAgent", err)
		}

		if len(agents) == 0 {
			return NewError("Predict/FindAgent",
				errors.Errorf("unable to find agents for framework=%s:%s model=%s:%s",
					frameworkName, frameworkVersion, modelName, modelVersion,
				))
		}
		return nil
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Minute
	err = backoff.Retry(find, backoff.WithMaxRetries(b, findMaxTries))

	return
}

func (p *PredictHandler) Open(params predict.OpenParams) middleware.Responder {
	frameworkName := dl.CleanString(getBody(params.Body.FrameworkName, "*"))
	frameworkVersion := dl.CleanString(getBody(params.Body.FrameworkVersion, "*"))
	modelName := dl.CleanString(getBody(params.Body.ModelName, "*"))
	modelVersion := dl.CleanString(getBody(params.Body.ModelVersion, "*"))

	agents, err := p.findAgent(params)
	if err != nil {
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

func (p *PredictHandler) closeClient(ctx context.Context, id string, force bool) error {
	client, err := p.getClient(id)
	if err != nil {
		return err
	}
	if _, err := client.Close(ctx, &dl.PredictorCloseRequest{Predictor: &dl.Predictor{ID: id}, Force: force}); err != nil {
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
	id := params.Body.Predictor.ID
	force := params.Body.Force

	if err := p.closeClient(ctx, id, force); err != nil {
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

func compress(data interface{}) (strfmt.Base64, error) {
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	out := snappy.Encode(nil, bts)
	// return strfmt.Base64(base64.StdEncoding.EncodeToString(out)), nil
	return strfmt.Base64(out), nil
}

func toJPEGFromFloat32Slice(fs []float32, width, height, channels int32) (strfmt.Base64, error) {
	offset := 0
	img := goimage.NewRGBA(goimage.Rect(0, 0, int(width), int(height)))
	for h := 0; h < int(height); h++ {
		for w := 0; w < int(width); w++ {
			R := uint8(fs[offset+0])
			G := uint8(fs[offset+1])
			B := uint8(fs[offset+2])
			img.Set(w, h, color.RGBA{R, G, B, 255})
			offset += 3
		}
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return strfmt.Base64(buf.Bytes()), nil
}

func toJPEGFromInt32Slice(fs []int32, width, height, channels int32) (strfmt.Base64, error) {
	offset := 0
	img := goimage.NewRGBA(goimage.Rect(0, 0, int(width), int(height)))
	for h := 0; h < int(height); h++ {
		for w := 0; w < int(width); w++ {
			R := uint8(fs[offset+0])
			G := uint8(fs[offset+1])
			B := uint8(fs[offset+2])
			img.Set(w, h, color.RGBA{R, G, B, 255})
			offset += 3
		}
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return strfmt.Base64(buf.Bytes()), nil
}

func toDlframeworkFeaturesResponse(responses []*dl.FeatureResponse) []*webmodels.DlframeworkFeatureResponse {
	resps := make([]*webmodels.DlframeworkFeatureResponse, len(responses))
	for ii, fr := range responses {
		features := make([]*webmodels.DlframeworkFeature, len(fr.Features))
		for jj, f := range fr.Features {
			features[jj] = &webmodels.DlframeworkFeature{
				ID:          f.ID,
				Metadata:    f.Metadata,
				Probability: f.Probability,
				Type:        webmodels.DlframeworkFeatureType(dl.FeatureType_name[int32(f.Type)]),
			}
			switch feature := f.Feature.(type) {
			case *dl.Feature_Classification:
				features[jj].Classification = &webmodels.DlframeworkClassification{
					Index: feature.Classification.Index,
					Label: feature.Classification.Label,
				}
			case *dl.Feature_BoundingBox:
				features[jj].BoundingBox = &webmodels.DlframeworkBoundingBox{
					Index: feature.BoundingBox.Index,
					Label: feature.BoundingBox.Label,
					Xmax:  feature.BoundingBox.Xmax,
					Xmin:  feature.BoundingBox.Xmin,
					Ymax:  feature.BoundingBox.Ymax,
					Ymin:  feature.BoundingBox.Ymin,
				}
			case *dl.Feature_SemanticSegment:
				features[jj].SemanticSegment = &webmodels.DlframeworkSemanticSegment{
					Height:  feature.SemanticSegment.Height,
					Width:   feature.SemanticSegment.Width,
					IntMask: feature.SemanticSegment.IntMask,
				}
			case *dl.Feature_InstanceSegment:
				features[jj].InstanceSegment = &webmodels.DlframeworkInstanceSegment{
					Index:     feature.InstanceSegment.Index,
					Label:     feature.InstanceSegment.Label,
					Xmax:      feature.InstanceSegment.Xmax,
					Xmin:      feature.InstanceSegment.Xmin,
					Ymax:      feature.InstanceSegment.Ymax,
					Ymin:      feature.InstanceSegment.Ymin,
					MaskType:  feature.InstanceSegment.MaskType,
					Height:    feature.InstanceSegment.Height,
					Width:     feature.InstanceSegment.Width,
					IntMask:   feature.InstanceSegment.IntMask,
					FloatMask: feature.InstanceSegment.FloatMask,
				}
			case *dl.Feature_Image:
				pp.Println("Feature_Image")
				features[jj].Image = &webmodels.DlframeworkImage{
					Data: feature.Image.Data,
				}
			case *dl.Feature_RawImage:
				pp.Println("Feature_RawImage")
				features[jj].RawImage = &webmodels.DlframeworkRawImage{
					Channels: feature.RawImage.Channels,
					Height:   feature.RawImage.Height,
					Width:    feature.RawImage.Width,
				}

				if feature.RawImage.FloatList != nil && feature.RawImage.CharList != nil {
					panic("cannot have both float and char list values")
				}

				if feature.RawImage.FloatList != nil && compressRawImage {
					features[jj].RawImage.DataType = "float32"
					compressed, err := compress(feature.RawImage.FloatList)
					if err != nil {
						panic("failed to compress float list data")
					}
					features[jj].RawImage.CompressedData = compressed
				}

				if feature.RawImage.FloatList != nil && !compressRawImage {
					features[jj].RawImage.DataType = "float32"
					jpg, err := toJPEGFromFloat32Slice(feature.RawImage.FloatList, feature.RawImage.Width, feature.RawImage.Height, feature.RawImage.Channels)
					if err != nil {
						panic("failed to compress float list data")
					}
					features[jj].RawImage.JpegData = jpg
				}

				if feature.RawImage.CharList != nil && !compressRawImage {
					features[jj].RawImage.DataType = "uint8"
					jpg, err := toJPEGFromInt32Slice(feature.RawImage.CharList, feature.RawImage.Width, feature.RawImage.Height, feature.RawImage.Channels)
					if err != nil {
						panic("failed to compress char list data")
					}
					features[jj].RawImage.JpegData = jpg
				}

			case *dl.Feature_Text:
				features[jj].Text = &webmodels.DlframeworkText{
					Data: feature.Text.Data,
				}
			case *dl.Feature_Region:
				features[jj].Region = &webmodels.DlframeworkRegion{
					Data:   feature.Region.Data,
					Format: feature.Region.Format,
				}
			case *dl.Feature_Audio:
				features[jj].Audio = &webmodels.DlframeworkAudio{
					Data:   feature.Audio.Data,
					Format: feature.Audio.Format,
				}
			case *dl.Feature_Geolocation:
				features[jj].Geolocation = &webmodels.DlframeworkGeoLocation{
					Index:     feature.Geolocation.Index,
					Latitude:  feature.Geolocation.Latitude,
					Longitude: feature.Geolocation.Longitude,
				}
			case *dl.Feature_Raw:
				features[jj].Raw = &webmodels.DlframeworkRaw{
					Data:   feature.Raw.Data,
					Format: feature.Raw.Format,
				}
			default:
				pp.Println("unhandeled feature type")
				pp.Println(feature)
			}

			// Index:       cast.ToInt32(f.Index),
			// Name:        f.Name,
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

	images := make([]*dl.Image, len(params.Body.Images))
	for ii, image := range params.Body.Images {
		images[ii] = &dl.Image{
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
