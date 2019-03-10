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

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/Unknwon/com"
	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/pp"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/dlframework/steps"
	"github.com/rai-project/evaluation"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"
	"gopkg.in/mgo.v2/bson"

	_ "github.com/rai-project/monitoring/monitors"
)

var (
	urlsFilePath   string
	duplicateInput int
	numUrlParts    int
)

var predictUrlsCmd = &cobra.Command{
	Use:     "urls",
	Short:   "Evaluates the urls using the specified model and framework",
	Aliases: []string{"url"},
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

		model, err := framework.FindModel(modelName + ":" + modelVersion)
		if err != nil {
			return err
		}

		predictors, err := agent.GetPredictors(framework)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to get predictor for %s. make sure you have "+
					"imported the framework's predictor package",
				framework.MustCanonicalName(),
			)
		}

		var predictorHandle common.Predictor
		for _, pred := range predictors {
			predModality, err := pred.Modality()
			if err != nil {
				continue
			}
			modelModality, err := model.Modality()
			if err != nil {
				continue
			}
			if predModality == modelModality {
				predictorHandle = pred
				break
			}
		}
		if predictorHandle == nil {
			return errors.New("unable to find predictor for requested modality")
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

		predictor, err := predictorHandle.Load(ctx, *model, options.PredictorOptions(predOpts))
		if err != nil {
			return err
		}

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

		tmp := urls
		for ii := 1; ii < duplicateInput; ii++ {
			urls = append(urls, tmp...)
		}

		log.WithField("urls_file_path", urlsFilePath).
			Debug("using the specified urls file path")

		if len(urls) == 0 {
			log.WithError(err).Error("the urls file has no url")
			os.Exit(-1)
		}

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

		urlParts := partitionList(urls, partitionListSize)

		cntTop1 := 0
		cntTop5 := 0

		outputs := make(chan interface{}, DefaultChannelBuffer)
		partlabels := map[string]string{}

		log.WithField("url_parts_length", len(urlParts)).
			WithField("url_parts_element_length", len(urlParts[0])).
			WithField("urls_length", len(urls)).
			WithField("using_gpu", useGPU).
			Info("starting inference on urls")

		if numUrlParts == -1 {
			numUrlParts = len(urlParts)
		}

		inferenceProgress := newProgress("infering", len(urls))

		for _, part := range urlParts[0:numUrlParts] {
			input := make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for _, url := range part {
					id := uuid.NewV4()
					lbl := steps.NewIDWrapper(id, url)
					partlabels[lbl.GetID()] = "" // no label for the input url
					input <- lbl
				}
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
				Then(steps.NewPredict(predictor)).
				Run(input)

			inferenceProgress.Add(batchSize)

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

			pp.Println(features[:3])
			return nil

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

			if features[0].Type != dl.FeatureType_CLASSIFICATION {
				panic("expecting a Classification type")
			}

			label = strings.TrimSpace(strings.ToLower(label))
			if strings.TrimSpace(strings.ToLower(features[0].Feature.(*dl.Feature_Classification).Classification.GetLabel())) == label {
				cntTop1++
			}
			for _, f := range features[:5] {
				if strings.TrimSpace(strings.ToLower(f.Feature.(*dl.Feature_Classification).Classification.GetLabel())) == label {
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
	sourcePath := sourcepath.MustAbsoluteDir()
	defaultURLsPath := filepath.Join(sourcePath, "..", "client", "run", "urlsfile")
	if !com.IsFile(defaultURLsPath) {
		defaultURLsPath = ""
	}
	defaultDuplicateInput := 1
	predictUrlsCmd.PersistentFlags().IntVar(&duplicateInput, "duplicate_input", defaultDuplicateInput, "duplicate the input urls ine urls_file")
	predictUrlsCmd.PersistentFlags().StringVar(&urlsFilePath, "urls_file_path", defaultURLsPath, "the path of the file containing the urls to perform the evaluations on.")
	predictDatasetCmd.PersistentFlags().IntVar(&numUrlParts, "num_url_parts", -1, "the number of url parts to process. Setting url parts to a value other than -1 means that only the first num_url_parts * partition_list_size images are infered from the dataset. This is useful while performing performance evaluations, where only a few hundred evaluation samples are useful")
}
