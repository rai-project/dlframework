package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	"github.com/rai-project/dldataset"
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
	"gopkg.in/cheggaaa/pb.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	datasetCategory      string
	datasetName          string
	modelName            string
	modelVersion         string
	batchSize            int
	numFileParts         int
	numWarmupFileParts   int
	partitionDatasetSize int
	publishEvaluation    bool
	publishPredictions   bool
	useGPU               bool
	failOnFirstError     bool
	traceLevelName       string
	traceLevel           tracer.Level = tracer.APPLICATION_TRACE
	databaseAddress      string
	databaseName         string
	databaseEndpoints    []string
	traceServerAddress   string
	DefaultChannelBuffer = 100000
)

func newProgress(prefix string, count int) *pb.ProgressBar {
	// get the new original progress bar.
	//bar := pb.New(count).Prefix(prefix)
	// TODO: set prefix of bar
	bar := pb.New(count)
	//bar.Set("prefix", prefix)

	// Refresh rate for progress bar is set to 100 milliseconds.
	bar.SetRefreshRate(time.Millisecond * 100)

	bar.SetTemplateString(string(pb.Full))
	bar.Start()
	return bar
}

var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "Evaluates the dataset using the specified model and framework",
	PreRunE: func(c *cobra.Command, args []string) error {
		if partitionDatasetSize == 0 {
			partitionDatasetSize = batchSize
		}
		traceLevel = tracer.LevelFromName(traceLevelName)

		if databaseName == "" {
			databaseName = config.App.Name
		}
		if databaseAddress != "" {
			databaseEndpoints = []string{databaseAddress}
		}
		if useGPU && !nvidiasmi.HasGPU {
			return errors.New("unable to find gpu on the system")
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		span, ctx := tracer.StartSpanFromContext(context.Background(), traceLevel, "dataset")
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

		if datasetName == "ilsvrc2012_validation" {
			var imagePredictor common.ImagePredictor

			err := deepcopier.Copy(predictor).To(&imagePredictor)
			if err != nil {
				return errors.Errorf("failed to copy to an image predictor for %v", model.MustCanonicalName())
			}
			dims, err := imagePredictor.GetImageDimensions()
			if err != nil {
				return err
			}
			if len(dims) != 3 {
				return errors.Errorf("expecting a 3 element vector for dimensions %v", dims)
			}
			width, height := dims[1], dims[2]
			if width != height {
				return errors.Errorf("expecting a square image dimensions width = %v, height = %v", width, height)
			}

			datasetName = fmt.Sprintf("%s_%v", datasetName, width)
		}

		log.WithField("dataset_category", datasetCategory).
			WithField("dataset_name", datasetName).
			Debug("using specified dataset")

		dataset, err := dldataset.Get(datasetCategory, datasetName)
		if err != nil {
			return err
		}
		defer dataset.Close()

		err = dataset.Download(ctx)
		if err != nil {
			return err
		}

		fileList, err := dataset.List(ctx)
		if err != nil {
			return err
		}

		err = dataset.Load(ctx)
		if err != nil {
			return err
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
			DatasetCategory:     dataset.Category(),
			DatasetName:         dataset.Name(),
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

		fileNameParts := partitionDataset(fileList, partitionDatasetSize)

		cntTop1 := 0
		cntTop5 := 0

		outputs := make(chan interface{}, DefaultChannelBuffer)
		partlabels := map[string]string{}

		log.WithField("file_name_parts_length", len(fileNameParts)).
			WithField("file_name_parts_element_length", len(fileNameParts[0])).
			WithField("file_list_length", len(fileList)).
			WithField("using_gpu", useGPU).
			Info("starting inference on dataset")

		if numWarmupFileParts != 0 && numFileParts != -1 {
			panic("todo")
		}

		if numFileParts == -1 {
			numFileParts = len(fileNameParts)
		}
		inferenceProgress := newProgress("infering", numFileParts)
		for _, part := range fileNameParts[0:numFileParts] {
			input := make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for range part {
					lda, err := dataset.Next(ctx)
					if err != nil {
						continue
					}
					id := uuid.NewV4()
					lbl := steps.NewIDWrapper(id, lda)
					partlabels[lbl.GetID()] = lda.Label()
					input <- lbl
				}
			}()

			output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
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

		databaseInsertProgress := newProgress("inserting prediction", numFileParts*partitionDatasetSize)

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
			Top1:      float64(cntTop1) / float64(len(fileList)),
			Top5:      float64(cntTop5) / float64(len(fileList)),
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

func partitionDataset(in []string, partitionSize int) (out [][]string) {
	cnt := (len(in)-1)/partitionSize + 1
	for i := 0; i < cnt; i++ {
		start := i * partitionSize
		end := (i + 1) * partitionSize
		if end > len(in) {
			end = len(in)
		}
		part := in[start:end]
		out = append(out, part)
	}

	return out
}

func init() {
	datasetCmd.PersistentFlags().StringVar(&datasetCategory, "dataset_category", "vision", "the dataset category to use for prediction")
	datasetCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "the name of the database to publish the evaluation results to")
	datasetCmd.PersistentFlags().StringVar(&modelName, "model_name", "BVLC-AlexNet", "the name of the model to use for prediction")
	datasetCmd.PersistentFlags().StringVar(&modelVersion, "model_version", "1.0", "the version of the model to use for prediction")
	datasetCmd.PersistentFlags().IntVarP(&partitionDatasetSize, "partition_dataset_size", "p", 0, "the chunk size to partition the input dataset. By default this is the same as the batch size")
	datasetCmd.PersistentFlags().IntVarP(&batchSize, "batch_size", "b", 64, "the batch size to use while performing inference")
	datasetCmd.PersistentFlags().IntVar(&numWarmupFileParts, "warmup_num_file_parts", 10, "the number of file parts to process during the warmup period. this is ignored if num_file_parts=-1")
	datasetCmd.PersistentFlags().IntVar(&numFileParts, "num_file_parts", -1, "the number of file parts to process. Setting file parts to a value other than -1 means that only the first num_file_parts * batch_size images are infered from the dataset. This is useful while performing performance evaluations, where only a few hundred evaluation samples are useful")
	datasetCmd.PersistentFlags().BoolVar(&failOnFirstError, "fail_on_error", false, "turning on causes the process to terminate/exit upon first inference error. This is useful since some inferences will result in an error because they run out of memory")
	datasetCmd.PersistentFlags().BoolVar(&publishEvaluation, "publish", true, "whether to publish the evaluation to database. Turning this off will not publish anything to the database. This is ideal for using carml within profiling tools or performing experiments where the terminal output is sufficient.")
	datasetCmd.PersistentFlags().BoolVar(&useGPU, "gpu", false, "whether to enable the gpu. An error is returned if the gpu is not available")
	datasetCmd.PersistentFlags().StringVar(&traceLevelName, "trace_level", traceLevel.String(), "the trace level to use while performing evaluations")
	datasetCmd.PersistentFlags().BoolVar(&publishPredictions, "publish_predictions", false, "whether to publish prediction results to database. This will store all the probability outputs for the evaluation in the database which would be a few gigabytes of data for one dataset")
	datasetCmd.PersistentFlags().StringVar(&datasetName, "dataset_name", "ilsvrc2012_validation", "the name of the dataset to perform the evaluations on. When using `ilsvrc2012_validation`, optimized versions of the dataset are used when the input network takes 224 or 22 ")
	datasetCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "the address of the mongo database to store the results. By default the address in the config `database.endpoints` is used")
	datasetCmd.PersistentFlags().StringVar(&traceServerAddress, "tracer_address", "34.207.139.117:16686", "the address of the jaeger or the zipking trace server")
}
