package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	"github.com/mailru/easyjson"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/archive"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	dlcmd "github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/dlframework/steps"
	"github.com/rai-project/evaluation"
	machine "github.com/rai-project/machine/info"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"
	"gopkg.in/mgo.v2/bson"
)

var (
	urlsFilePath      string
	duplicateInput    int
	numUrlParts       int
	numWarmUpUrlParts int
)

func runPredictUrlsCmd(c *cobra.Command, args []string) error {
	model, err := framework.FindModel(modelName + ":" + modelVersion)
	if err != nil {
		return err
	}
	log.WithField("model", modelName).Info("running predict urls")

	var device string
	if useGPU {
		device = "gpu"
	} else {
		device = "cpu"
	}
	baseDir := filepath.Join("experiments", framework.Name, framework.Version, model.Name, model.Version, strconv.Itoa(batchSize), device, hostName)
	if !com.IsDir(baseDir) {
		os.MkdirAll(baseDir, os.ModePerm)
	}
	ts := strings.ToLower(tracer.LevelToName(traceLevel))
	tracerFileName := "trace_" + ts + ".json"
	tracerFilePath := filepath.Join(baseDir, tracerFileName)
	if (publishEvaluation == false) && com.IsFile(tracerFilePath) {
		log.WithField("path", tracerFilePath).Info("trace file already exists")
		return nil
	}

	if useGPU {
		if bts, err := json.Marshal(nvidiasmi.Info); err == nil {
			ioutil.WriteFile(filepath.Join(baseDir, "nvidia_info.json"), bts, 0644)
		}
	}

	if machine.Info != nil && machine.Info.Hostname != "" {
		bts, err := json.Marshal(machine.Info)
		if err == nil {
			ioutil.WriteFile(filepath.Join(baseDir, "system_info.json"), bts, 0644)
		}
	}

	rootSpan, ctx := tracer.StartSpanFromContext(
		context.Background(),
		tracer.APPLICATION_TRACE,
		"evaluation_predict_urls",
		opentracing.Tags{
			"framework_name":     framework.Name,
			"framework_version":  framework.Version,
			"model_name":         modelName,
			"model_version":      modelVersion,
			"use_gpu":            useGPU,
			"batch_size":         batchSize,
			"num_warmup_batches": numWarmUpUrlParts,
		},
	)
	if rootSpan == nil {
		panic("invalid span")
	}

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

	predictor, err := predictorHandle.Load(
		ctx,
		*model,
		options.PredictorOptions(predOpts),
		options.DisableFrameworkAutoTuning(false),
	)
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

	log.WithField("urls_file_path", urlsFilePath).
		Debug("using the specified urls file path")

	if len(urls) == 0 {
		log.WithError(err).Error("the urls file has no url")
		os.Exit(-1)
	}

	tmp := urls
	if duplicateInput < batchSize {
		duplicateInput = batchSize
	}
	for ii := 1; ii < duplicateInput; ii++ {
		urls = append(urls, tmp...)
	}

	inputPredictionIds := []bson.ObjectId{}

	hostName, _ := os.Hostname()
	metadata := map[string]string{}
	if useGPU {
		if bts, err := json.Marshal(nvidiasmi.Info); err == nil {
			metadata["nvidia_smi"] = string(bts)
			rootSpan.SetTag("nvidia_smi", string(bts))
		}
	}

	// Dummy userID and runID hardcoded
	// TODO read userID from manifest file
	// calculate runID from table
	userID := "evaluator"
	runID := uuid.NewV4()

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

	preprocessOptions, err := predictor.GetPreprocessOptions()
	if err != nil {
		return err
	}

	urlParts := dl.PartitionStringList(urls, partitionListSize)

	oldTraceLevel := tracer.GetLevel()
	tracer.SetLevel(tracer.NO_TRACE)

	outputs := make(chan interface{}, DefaultChannelBuffer)
	partlabels := map[string]string{}

	log.WithField("model", modelName).
		WithField("urls_length", len(urls)).
		WithField("url_parts_length", len(urlParts)).
		WithField("partition_list_size", partitionListSize).
		WithField("num_url_part", numUrlParts).
		WithField("num_warmup_url_parts", numWarmUpUrlParts).
		WithField("using_gpu", useGPU).
		Info("starting inference on urls")

	if numWarmUpUrlParts != 0 {
		for _, part := range urlParts[0:numWarmUpUrlParts] {
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

			output := pipeline.New(pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewReadURL()).
				Then(steps.NewReadImage(preprocessOptions)).
				Then(steps.NewPreprocessImage(preprocessOptions)).
				Run(input)

			var images []interface{}
			for out := range output {
				images = append(images, out)
			}

			imageParts := dl.Partition(images, batchSize)

			input = make(chan interface{}, DefaultChannelBuffer)
			go func() {
				defer close(input)
				for _, p := range imageParts {
					input <- p
				}
			}()

			output = pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
				Then(steps.NewPredict(predictor)).
				Run(input)

			for o := range output {
				if err, ok := o.(error); ok && failOnFirstError {
					log.WithError(err).Error("encountered an error while performing warmup inference")
					os.Exit(-1)
				}
				outputs <- o
			}
		}
	}

	close(outputs)
	for range outputs {
	}

	tracer.SetLevel(oldTraceLevel)

	outputs = make(chan interface{}, DefaultChannelBuffer)

	if numUrlParts == -1 {
		numUrlParts = len(urlParts)
	}

	rootSpan.SetTag("num_batches", numUrlParts)

	urlCnt := len(urls)

	inferenceProgress := dlcmd.NewProgress("inferring", urlCnt)

	for ii, part := range urlParts[:numUrlParts] {
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

		evaluateBatchSpan, evaluateBatchCtx := tracer.StartSpanFromContext(
			ctx,
			tracer.APPLICATION_TRACE,
			"evaluate_batch",
			opentracing.Tags{
				"batch_index": ii,
				"batch_size":  batchSize,
			},
		)

		opts := []pipeline.Option{pipeline.ChannelBuffer(DefaultChannelBuffer)}
		if tracePreprocess == true {
			opts = append(opts, pipeline.Context(evaluateBatchCtx))
		}
		output := pipeline.New(opts...).
			Then(steps.NewReadURL()).
			Then(steps.NewReadImage(preprocessOptions)).
			Then(steps.NewPreprocessImage(preprocessOptions)).
			Run(input)

		var images []interface{}
		for out := range output {
			images = append(images, out)
		}

		imageParts := dl.Partition(images, batchSize)

		input = make(chan interface{}, DefaultChannelBuffer)
		go func() {
			defer close(input)
			for _, p := range imageParts {
				input <- p
			}
		}()

		output = pipeline.New(pipeline.Context(evaluateBatchCtx), pipeline.ChannelBuffer(DefaultChannelBuffer)).
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
		evaluateBatchSpan.Finish()
	}

	//inferenceProgress.FinishPrint("inference complete")
	inferenceProgress.Finish()

	rootSpan.Finish()
	tracer.ResetStd()

	close(outputs)

	traceID := rootSpan.Context().(jaeger.SpanContext).TraceID()
	traceIDVal := traceID.String()
	if runtime.GOARCH == "ppc64le" {
		traceIDVal = strconv.FormatUint(traceID.Low, 16)
	}
	pp.Println(fmt.Sprintf("http://%s:16686/trace/%v", getTracerHostAddress(tracerAddress), traceIDVal))

	query := fmt.Sprintf("http://%s:16686/api/traces/%v?raw=true", getTracerHostAddress(tracerAddress), traceIDVal)
	resp, err := grequests.Get(query, nil)
	if err != nil {
		log.WithError(err).
			WithField("trace_id", traceIDVal).
			Error("failed to download span information")
	}
	log.WithField("model", modelName).WithField("trace_id", traceIDVal).WithField("query", query).Info("downloaded trace information")

	if publishEvaluation == false {
		for range outputs {
		}

		err := ioutil.WriteFile(tracerFilePath, resp.Bytes(), 0644)
		if err != nil {
			return err
		}

		if false {
			archiver, err := archive.Zip(tracerFilePath, archive.BZip2Format())
			if err != nil {
				return err
			}

			archiveFilePath := strings.TrimSuffix(tracerFilePath, ".json") + ".tar.bz2"
			f, err := os.Create(archiveFilePath)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, archiver)
		}
		log.WithField("model", modelName).WithField("path", tracerFilePath).Info("publishEvaluation is false, wrote the trace to a local file")

		return nil
	}

	cnt := 0
	cntTop1 := 0
	cntTop5 := 0

	databaseInsertProgress := dlcmd.NewProgress("inserting prediction", batchSize)

	for out0 := range outputs {
		if cnt > urlCnt {
			break
		}
		out, ok := out0.(steps.IDer)
		if !ok {
			return errors.Errorf("expecting steps.IDer, but got %v", out0)
		}
		id := out.GetID()
		label := partlabels[id]

		features := out.GetData().(dl.Features)
		if !ok {
			return errors.Errorf("expecting a dlframework.Features type, but got %v", out.GetData())
		}

		if publishPredictions == true {
			log.WithField("model", modelName).Info("inserting predictions into inputPredictionsTable")

			inputPrediction := evaluation.InputPrediction{
				ID:            bson.NewObjectId(),
				CreatedAt:     time.Now(),
				InputID:       id,
				ExpectedLabel: label,
				Features:      features,
			}

			if err := inputPredictionsTable.Insert(inputPrediction); err != nil {
				log.WithError(err).Errorf("failed to insert input prediction into database")
			}
			inputPredictionIds = append(inputPredictionIds, inputPrediction.ID)
		}

		databaseInsertProgress.Increment()

		features.Sort()

		if features[0].Type == dl.FeatureType_CLASSIFICATION {
			label = strings.TrimSpace(strings.ToLower(label))
			if strings.TrimSpace(strings.ToLower(features[0].Feature.(*dl.Feature_Classification).Classification.GetLabel())) == label {
				cntTop1++
			}
			for _, f := range features[:5] {
				if strings.TrimSpace(strings.ToLower(f.Feature.(*dl.Feature_Classification).Classification.GetLabel())) == label {
					cntTop5++
				}
			}
		} else {
			panic("expecting a Classification type")
		}
		cnt++
	}

	//databaseInsertProgress.FinishPrint("inserting prediction complete")
	databaseInsertProgress.Finish()
	log.WithField("model", modelName).Info("finised inserting prediction")

	modelAccuracy := evaluation.ModelAccuracy{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		Top1:      float64(cntTop1) / float64(urlCnt),
		Top5:      float64(cntTop5) / float64(urlCnt),
	}
	if err := modelAccuracyTable.Insert(modelAccuracy); err != nil {
		log.WithError(err).Error("failed to publish model accuracy entry")
	}

	log.WithField("model", modelName).Info("downloading trace information")

	var trace evaluation.TraceInformation
	err = easyjson.UnmarshalFromReader(resp, &trace)
	if err != nil {
		log.WithError(err).Error("failed to decode trace information")
	}
	performance := evaluation.Performance{
		ID:         bson.NewObjectId(),
		CreatedAt:  time.Now(),
		Trace:      trace,
		TraceLevel: traceLevel,
	}
	performanceTable.Insert(performance)

	log.WithField("model", modelName).Info("inserted performance information")

	evaluationEntry.PerformanceID = performance.ID
	evaluationEntry.ModelAccuracyID = modelAccuracy.ID
	evaluationEntry.InputPredictionIDs = inputPredictionIds

	if err := evaluationTable.Insert(evaluationEntry); err != nil {
		log.WithError(err).Error("failed to publish evaluation entry")
	}

	log.WithField("model", model.MustCanonicalName()).
		WithField("accuracy", spew.Sprint(modelAccuracy)).
		Info("inserted evaluation information")

	return nil
}

var predictUrlsCmd = &cobra.Command{
	Use:     "urls",
	Short:   "Evaluate the urls using the specified model and framework",
	Aliases: []string{"url"},
	RunE: func(c *cobra.Command, args []string) error {
		if modelName == "all" {
			for _, model := range framework.Models() {
				modelName = model.Name
				modelVersion = model.Version
				runPredictUrlsCmd(c, args)
			}
			return nil
		}
		return runPredictUrlsCmd(c, args)
	},
}

func init() {
	sourcePath := sourcepath.MustAbsoluteDir()
	defaultURLsPath := filepath.Join(sourcePath, "..", "_fixtures", "urlsfile")
	if !com.IsFile(defaultURLsPath) {
		defaultURLsPath = ""
	}
	defaultDuplicateInput := 1
	predictUrlsCmd.PersistentFlags().IntVar(&duplicateInput, "duplicate_input", defaultDuplicateInput, "duplicate the input urls in urls_file")
	predictUrlsCmd.PersistentFlags().StringVar(&urlsFilePath, "urls_file_path", defaultURLsPath, "the path of the file containing the urls to perform the evaluations on.")
	predictUrlsCmd.PersistentFlags().IntVar(&numUrlParts, "num_url_parts", -1, "the number of url parts to process. Setting url parts to a value other than -1 means that only the first num_url_parts * partition_list_size images are infered from the dataset. This is useful while performing performance evaluations, where only a few hundred evaluation samples are useful")
	predictUrlsCmd.PersistentFlags().IntVar(&numWarmUpUrlParts, "num_warmup_url_parts", 1, "the number of url parts to process during the warmup period. This is ignored if num_file_parts=-1")
}
