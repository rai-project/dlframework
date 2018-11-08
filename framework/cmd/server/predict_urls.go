package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/dlframework/steps"
	"github.com/rai-project/evaluation"
	_ "github.com/rai-project/monitoring/monitors"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/ulule/deepcopier"
	"gopkg.in/mgo.v2/bson"
)

var (
	urlsFilePath string
)

var predictUrlsCmd = &cobra.Command{
	Use:     "url",
	Short:   "Evaluates the urls using the specified model and framework",
	Aliases: []string{"urls"},
	PreRunE: func(c *cobra.Command, args []string) error {
		traceLevel = tracer.LevelFromName(traceLevelName)
		if useGPU && !nvidiasmi.HasGPU {
			return errors.New("unable to find gpu on the system")
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		span, ctx := tracer.StartSpanFromContext(context.Background(), traceLevel, "urls")
		defer func() {
			if span != nil {
				span.Finish()
			}
		}()

		opts := []database.Option{}
		if len(databaseEndpoints) != 0 {
			opts = append(opts, database.Endpoints(databaseEndpoints))
		}
		db, err := mongodb.NewDatabase(databaseName, opts...)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to create new database %s at %v",
				databaseName, databaseEndpoints,
			)
		}
		defer db.Close()

		predictorFramework, err := agent.GetPredictor(framework)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to get predictor for %s. make sure you have "+
					"imported the framework's predictor package",
				framework.MustCanonicalName(),
			)
		}

		model, err := framework.FindModel(modelName + ":" + modelVersion)
		if err != nil {
			return err
		}

		var dc map[string]int32
		if useGPU {
			if !nvidiasmi.HasGPU {
				return errors.New("not gpu found")
			}
			dc = map[string]int32{"GPU": 0}
		} else {
			dc = map[string]int32{"CPU": 0}
		}

		execOpts := &dl.ExecutionOptions{
			TraceLevel: dl.ExecutionOptions_TraceLevel(
				dl.ExecutionOptions_TraceLevel_value[traceLevel.String()],
			),
			DeviceCount: dc,
		}
		predOpts := &dl.PredictionOptions{
			FeatureLimit:     10,
			BatchSize:        int32(batchSize),
			ExecutionOptions: execOpts,
		}

		predictor, err := predictorFramework.Load(ctx, *model, options.PredictorOptions(predOpts))
		if err != nil {
			return err
		}

		var imagePredictor common.ImagePredictor

		err = deepcopier.Copy(predictor).To(&imagePredictor)
		if err != nil {
			return errors.Errorf("failed to copy to an image predictor for %v", model.MustCanonicalName())
		}
		// dims, err = imagePredictor.GetImageDimensions()
		// if err != nil {
		// 	return err
		// }
		// if len(dims) != 3 {
		// 	return errors.Errorf("expecting a 3 element vector for dimensions %v", dims)
		// }
		// width, height := dims[1], dims[2]
		// if width != height {
		// 	return errors.Errorf("expecting a square image dimensions width = %v, height = %v", width, height)
		// }

		inputPredictionIds := []bson.ObjectId{}

		hostName, _ := os.Hostname()
		metadata := map[string]string{}
		if useGPU {
			if bts, err := json.Marshal(nvidiasmi.Info); err == nil {
				metadata["nvidia_smi"] = string(bts)
			}
		}

		// Dummy userID and runID hardcoded
		// TODO read userID from manifest file
		// calculate runID from table
		userID := "admin"
		runID := 1

		evaluationEntry := evaluation.Evaluation{
			ID:                  bson.NewObjectId(),
			UserID:              userID,
			RunID:               runID,
			CreatedAt:           time.Now(),
			Framework:           *model.GetFramework(),
			Model:               *model,
			DatasetCategory:     "",
			DatasetName:         "",
			Public:              false,
			Hostname:            hostName,
			UsingGPU:            useGPU,
			BatchSize:           batchSize,
			TraceLevel:          traceLevel.String(),
			MachineArchitecture: runtime.GOARCH,
			Metadata:            metadata,
		}

		evaluationTable, err := evaluation.NewEvaluationCollection(db)
		if err != nil {
			return err
		}
		defer evaluationTable.Close()

		modelAccuracyTable, err := evaluation.NewModelAccuracyCollection(db)
		if err != nil {
			return err
		}
		defer modelAccuracyTable.Close()

		performanceTable, err := evaluation.NewPerformanceCollection(db)
		if err != nil {
			return err
		}
		defer performanceTable.Close()

		inputPredictionsTable, err := evaluation.NewInputPredictionCollection(db)
		if err != nil {
			return err
		}
		defer inputPredictionsTable.Close()

		preprocessOptions, err := predictor.GetPreprocessOptions(nil) // disable tracing
		if err != nil {
			return err
		}
		_ = preprocessOptions

		cntTop1 := 0
		cntTop5 := 0

		var urls []string
		urlsFilePath, err := filepath.Abs(urlsFilePath)
		if err != nil {
			return errors.Wrapf(err, "cannot get absolute path of %s", urlsFilePath)
		}
		f, err := os.Open(urlsFilePath)
		if err != nil {
			return errors.Wrapf(err, "cannot read %s", urlsFilePath)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			urls = append(urls, line)
		}
		// fill the batch with the same image
		if len(urls) == 1 {
			for ii := 0; ii < batchSize; ii++ {
				urls = append(urls, urls[0])
			}
		}
		// cleanNames()

		outputs := make(chan interface{}, DefaultChannelBuffer)
		partlabels := map[string]string{}

		log.WithField("urls_file_path", urlsFilePath).
			WithField("urls_length", len(urls)).
			WithField("using_gpu", useGPU).
			Info("starting inference on urls")

		inferenceProgress := newProgress("infering", len(urls))
		for _, url := range urls {
			input := make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				id := uuid.NewV4()
				lbl := steps.NewIDWrapper(id, url)
				partlabels[lbl.GetID()] = "" // no label for the input url
				input <- lbl
			}()

			output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewReadURL()).
				Then(steps.NewReadImage(preprocessOptions)).
				Then(steps.NewPreprocessImage(preprocessOptions)).
				Run(input)

			var images []interface{}
			for out := range output {
				images = append(images, out)
			}

			parts := agent.Partition(images, batchSize)

			input = make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for _, part := range parts {
					input <- part
				}
			}()

			output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewPredictImage(predictor)).
				Run(input)

			inferenceProgress.Increment()

			for o := range output {
				if err, ok := o.(error); ok && failOnFirstError {
					//inferenceProgress.FinishPrint("inference halted")
					inferenceProgress.Finish()

					log.WithError(err).Error("encountered an error while performing inference")
					os.Exit(-1)
				}
				outputs <- o
			}

		}
		span.Finish()
		defer func() {
			span = nil
		}()

		//inferenceProgress.FinishPrint("inference complete")
		inferenceProgress.Finish()

		close(outputs)

		if publishEvaluation == false {
			for range outputs {
			}
			return nil
		}

		databaseInsertProgress := newProgress("inserting prediction", batchSize)

		for out0 := range outputs {
			out, ok := out0.(steps.IDer)
			if !ok {
				return errors.Errorf("expecting steps.IDer, but got %v", out0)
			}
			_ = out
			id := out.GetID()
			label := partlabels[id]

			features := out.GetData().(dl.Features)
			if !ok {
				return errors.Errorf("expecting a dlframework.Features type, but got %v", out.GetData())
			}

			if publishPredictions == true {
				inputPrediction := evaluation.InputPrediction{
					ID:            bson.NewObjectId(),
					CreatedAt:     time.Now(),
					InputID:       id,
					ExpectedLabel: label,
					Features:      features,
				}

				err = inputPredictionsTable.Insert(inputPrediction)
				if err != nil {
					log.WithError(err).Errorf("failed to insert input prediction into database")
				}

				inputPredictionIds = append(inputPredictionIds, inputPrediction.ID)
			}
			databaseInsertProgress.Increment()

			features.Sort()

			label = strings.TrimSpace(strings.ToLower(label))
			if strings.TrimSpace(strings.ToLower(features[0].GetName())) == label {
				cntTop1++
			}
			for _, f := range features[:5] {
				if strings.TrimSpace(strings.ToLower(f.GetName())) == label {
					cntTop5++
				}
			}
		}

		//databaseInsertProgress.FinishPrint("inserting prediction complete")
		databaseInsertProgress.Finish()

		modelAccuracy := evaluation.ModelAccuracy{
			ID:        bson.NewObjectId(),
			CreatedAt: time.Now(),
			Top1:      float64(cntTop1) / float64(len(urls)),
			Top5:      float64(cntTop5) / float64(len(urls)),
		}
		if err := modelAccuracyTable.Insert(modelAccuracy); err != nil {
			log.WithError(err).Error("failed to publish model accuracy entry")
		}

		log.WithField("model", model.MustCanonicalName()).
			WithField("accuracy", spew.Sprint(modelAccuracy)).
			Info("finished publishing prediction result")

		traceID := span.Context().(jaeger.SpanContext).TraceID()
		traceIDVal := strconv.FormatUint(traceID.Low, 16)
		tracer.Close()
		query := fmt.Sprintf("http://%s/api/traces/%v", traceServerAddress, traceIDVal)
		resp, err := grequests.Get(query, nil)

		if err != nil {
			log.WithError(err).
				WithField("trace_id", traceIDVal).
				Error("failed to download span information")
		} else {
			log.WithField("trace_id", traceIDVal).
				Info("downloaded span information")
		}

		if err == nil {
			var trace evaluation.TraceInformation
			dec := json.NewDecoder(resp)
			if err := dec.Decode(&trace); err != nil {
				log.WithError(err).Error("failed to decode trace information")
			}
			performance := evaluation.Performance{
				ID:         bson.NewObjectId(),
				CreatedAt:  time.Now(),
				Trace:      trace,
				TraceLevel: traceLevel,
			}
			evaluationEntry.PerformanceID = performance.ID
			performanceTable.Insert(performance)
			log.Info("inserted span information")
		}
		evaluationEntry.ModelAccuracyID = modelAccuracy.ID
		evaluationEntry.InputPredictionIDs = inputPredictionIds

		log.Info("inserting evaluation information")
		if err := evaluationTable.Insert(evaluationEntry); err != nil {
			log.WithError(err).Error("failed to publish evaluation entry")
		}

		log.WithField("top1_accuracy", modelAccuracy.Top1).
			WithField("top5_accuracy", modelAccuracy.Top5).
			Info("done")

		return nil
	},
}

func init() {
	predictUrlsCmd.PersistentFlags().StringVar(&urlsFilePath, "urls_file_path", "../run/urls_file", "the path of the file containing the urls to perform the evaluations on.")
}
