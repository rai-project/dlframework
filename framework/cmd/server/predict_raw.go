package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/levigross/grequests"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	dlcmd "github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/dlframework/framework/options"
	common "github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/evaluation"
	machine "github.com/rai-project/machine/info"
	nvidiasmi "github.com/rai-project/nvidia-smi"
	"github.com/rai-project/tracer"
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"
	"gopkg.in/mgo.v2/bson"
)

var (
	numBatches       int
	numWarmUpBatches int
)

func runPredictRawCmd(c *cobra.Command, args []string) error {
	model, err := framework.FindModel(modelName + ":" + modelVersion)
	if err != nil {
		return err
	}
	log.WithField("model", modelName).Info("running predict urls")

	if publishToDatabase == true {
		opts := []database.Option{}
		if len(databaseEndpoints) != 0 {
			opts = append(opts, database.Endpoints(databaseEndpoints))
		}

		db, err = mongodb.NewDatabase(databaseName, opts...)
		if err != nil {
			return errors.Wrapf(err,
				"⚠️ failed to create new database %s at %v",
				databaseName, databaseEndpoints,
			)
		}
		defer db.Close()

		modelAccuracyTable, err = evaluation.NewModelAccuracyCollection(db)
		if err != nil {
			return err
		}
		defer modelAccuracyTable.Close()

		inputPredictionsTable, err = evaluation.NewInputPredictionCollection(db)
		if err != nil {
			return err
		}
		defer inputPredictionsTable.Close()

		evaluationTable, err = evaluation.NewEvaluationCollection(db)
		if err != nil {
			return err
		}
		defer evaluationTable.Close()

		performanceTable, err = evaluation.NewPerformanceCollection(db)
		if err != nil {
			return err
		}
		defer performanceTable.Close()

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
		log.WithField("gpu = ", nvidiasmi.Info.GPUS[gpuDeviceId].ProductName).Info("Running evalaution on GPU")
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
		GpuMetrics:       gpuMetrics,
		ExecutionOptions: execOpts,
	}

	rootSpan, ctx := tracer.StartSpanFromContext(
		context.Background(),
		tracer.APPLICATION_TRACE,
		"evaluation_predict_raw",
		opentracing.Tags{
			"framework_name":     framework.Name,
			"framework_version":  framework.Version,
			"model_name":         modelName,
			"model_version":      modelVersion,
			"batch_size":         batchSize,
			"use_gpu":            useGPU,
			"gpu_metrics":        gpuMetrics,
			"num_warmup_batches": numWarmUpBatches,
		},
	)
	if rootSpan == nil {
		panic("invalid span")
	}

	predictor, err := predictorHandle.Load(
		ctx,
		*model,
		options.PredictorOptions(predOpts),
		options.DisableFrameworkAutoTuning(disableFrameworkAutoTuning),
	)
	if err != nil {
		return err
	}

	inputPredictionIds := []bson.ObjectId{}

	log.WithField("model", modelName).
		WithField("using_gpu", useGPU).
		Info("starting inference using raw predictor")

	if numWarmUpBatches != 0 {
		warmUpSpan, warmUpSpanCtx := tracer.StartSpanFromContext(
			ctx,
			tracer.APPLICATION_TRACE,
			"warm_up",
			opentracing.Tags{
				"num_warmup_batches": numWarmUpBatches,
			},
		)
		tracer.SetLevel(tracer.NO_TRACE)
		for ii := 0; ii < numWarmUpBatches; ii++ {
			predictor.Predict(warmUpSpanCtx, nil, options.PredictorOptions(predOpts))
		}
		tracer.SetLevel(traceLevel)
		warmUpSpan.Finish()
	}

	hostName, _ := os.Hostname()
	hostIP := getHostIP()
	metadata := map[string]string{}
	if useGPU {
		if bts, err := json.Marshal(nvidiasmi.Info); err == nil {
			metadata["nvidia_smi"] = string(bts)
			rootSpan.SetTag("nvidia_smi", string(bts))
		}
	}

	inferenceProgress := dlcmd.NewProgress("inferring", numBatches)
	for ii := 0; ii < numBatches; ii++ {
		evaluateBatchSpan, evaluateBatchCtx := tracer.StartSpanFromContext(
			ctx,
			tracer.APPLICATION_TRACE,
			"evaluate_batch",
			opentracing.Tags{
				"batch_index": ii,
			},
		)
		for ii := 0; ii < numWarmUpBatches; ii++ {
			predictor.Predict(evaluateBatchCtx, nil, options.PredictorOptions(predOpts))
			inferenceProgress.Add(batchSize)
		}
		evaluateBatchSpan.Finish()
	}
	//inferenceProgress.FinishPrint("inference complete")
	inferenceProgress.Finish()
	rootSpan.Finish()
	tracer.ResetStd()

	traceID := rootSpan.Context().(jaeger.SpanContext).TraceID()
	traceIDVal := traceID.String()
	if runtime.GOARCH == "ppc64le" {
		traceIDVal = strconv.FormatUint(traceID.Low, 16)
	}
	tracerServerAddr := getTracerServerAddress(tracerAddress)
	pp.Println(fmt.Sprintf("the trace is at http://%s:16686/trace/%v", tracerServerAddr, traceIDVal))
	traceURL := fmt.Sprintf("http://%s:16686/api/traces/%v?raw=true", tracerServerAddr, traceIDVal)

	var device string
	if useGPU {
		device = "gpu"
	} else {
		device = "cpu"
		gpuDeviceId = -1
	}

	if publishToDatabase == false {
		outputDir := filepath.Join(baseDir, framework.Name, framework.Version, model.Name, model.Version, strconv.Itoa(batchSize), device, hostName)
		if !com.IsDir(outputDir) {
			os.MkdirAll(outputDir, os.ModePerm)
		}

		if useGPU {
			if bts, err := json.Marshal(nvidiasmi.Info); err == nil {
				ioutil.WriteFile(filepath.Join(outputDir, "nvidia_info.json"), bts, 0644)
			}
		}

		if machine.Info != nil && machine.Info.Hostname != "" {
			bts, err := json.Marshal(machine.Info)
			if err == nil {
				ioutil.WriteFile(filepath.Join(outputDir, "system_info.json"), bts, 0644)
			}
		}

		ts := strings.ToLower(tracer.LevelToName(traceLevel))
		traceFileName := "trace_" + ts + ".json"
		tracePath := filepath.Join(outputDir, traceFileName)
		if (publishToDatabase == false) && com.IsFile(tracePath) {
			log.WithField("path", tracePath).Info("trace file already exists")
			return nil
		}

		resp, err := grequests.Get(traceURL, nil)
		if err != nil {
			log.WithError(err).
				WithField("trace_id", traceIDVal).
				Error("failed to download span information")
		}
		log.WithField("model", modelName).WithField("trace_id", traceIDVal).WithField("traceURL", traceURL).Info("downloaded trace information")

		err = ioutil.WriteFile(tracePath, resp.Bytes(), 0644)
		if err != nil {
			return err
		}

		log.WithField("model", modelName).WithField("path", tracePath).Infof("publishToDatabase is false, writing the trace to a local file")

		pp.Println(fmt.Sprintf("the trace is at %v locally", tracePath))

		return nil
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
		HostIP:              hostIP,
		UsingGPU:            useGPU,
		BatchSize:           batchSize,
		GPUMetrics:          gpuMetrics,
		TraceLevel:          traceLevel.String(),
		MachineArchitecture: runtime.GOARCH,
		MachineInformation:  machine.Info,
		Metadata:            metadata,
	}

	if nvidiasmi.Info != nil {
		evaluationEntry.GPUDriverVersion = &nvidiasmi.Info.DriverVersion
		if useGPU {
			evaluationEntry.GPUDevice = &gpuDeviceId
			evaluationEntry.GPUInformation = &nvidiasmi.Info.GPUS[gpuDeviceId]
		}
	}

	modelAccuracy := evaluation.ModelAccuracy{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		Top1:      float64(0),
		Top5:      float64(0),
	}
	if err := modelAccuracyTable.Insert(modelAccuracy); err != nil {
		log.WithError(err).Error("failed to publish model accuracy entry")
	}

	log.WithField("model", modelName).Info("downloading trace information")

	performance := &evaluation.Performance{
		ID:         bson.NewObjectId(),
		CreatedAt:  time.Now(),
		Trace:      nil,
		TraceLevel: traceLevel,
		TraceURL:   traceURL,
	}

	if err = performance.CompressTrace(); err != nil {
		log.WithError(err).Error("failed to compress trace information")
	}

	if err := performanceTable.Insert(performance); err != nil {
		le := log.WithError(err)
		cause := errors.Cause(err)
		if cause != err {
			le = log.WithField("cause", cause.Error())
		}
		le.Error("failed to publish performance entry")
	}

	log.WithField("model", modelName).Info("inserted performance information")

	evaluationEntry.PerformanceID = performance.ID
	evaluationEntry.ModelAccuracyID = modelAccuracy.ID
	evaluationEntry.InputPredictionIDs = inputPredictionIds

	if err := evaluationTable.Insert(evaluationEntry); err != nil {
		le := log.WithError(err)
		cause := errors.Cause(err)
		if cause != err {
			le = log.WithField("cause", cause.Error())
		}
		le.Error("failed to publish evaluation entry")
	}

	log.WithField("model", modelName).Info("inserted evaluation information")

	return nil
}

var predictRawCmd = &cobra.Command{
	Use:     "raw",
	Short:   "Evaluate using the specified model and framework",
	Aliases: []string{""},
	RunE: func(c *cobra.Command, args []string) error {
		if modelName == "all" {
			for _, model := range framework.Models() {
				modelName = model.Name
				modelVersion = model.Version
				runPredictRawCmd(c, args)
			}
			return nil
		}
		return runPredictRawCmd(c, args)
	},
}

func init() {
	predictRawCmd.PersistentFlags().IntVar(&numBatches, "num_batches", 1, "the number of batches to evaluate.")
	predictRawCmd.PersistentFlags().IntVar(&numWarmUpBatches, "num_warmup_batches", 1, "the number of batches to process during the warmup period.")
}
