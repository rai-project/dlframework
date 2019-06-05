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
	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	"github.com/rai-project/dldataset"
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
	"github.com/rai-project/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
)

var (
	datasetCategory    string
	datasetName        string
	numFileParts       int
	numWarmUpFileParts int
	hostName, _        = os.Hostname()
	centerCrop         bool
)

func runPredictDatasetCmd(c *cobra.Command, args []string) error {
	model, err := framework.FindModel(modelName + ":" + modelVersion)
	if err != nil {
		return err
	}
	log.WithField("model", modelName).Info("running predict dataset")

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

	ctx := context.Background()

	predictor, err := predictorHandle.Load(
		ctx,
		*model,
		options.PredictorOptions(predOpts),
		options.DisableFrameworkAutoTuning(false),
	)
	if err != nil {
		return err
	}

	if datasetName == "ilsvrc2012_validation" {
		dims, err := model.GetInputDimensions()
		if err != nil {
			return err
		}
		if len(dims) < 3 {
			return errors.Errorf("expecting a 3 element vector for dimensions %v", dims)
		}
		width, height := dims[1], dims[2]
		if width != height {
			return errors.Errorf("expecting a square image dimensions width = %v, height = %v", width, height)
		}

		datasetName = fmt.Sprintf("%s_%v", datasetName, width)

		if centerCrop {
			datasetName += "_center_crop_875"
		}
	}

	log.WithField("dataset_category", datasetCategory).
		WithField("dataset_name", datasetName).
		Info("using the specified dataset")

	dataset, err := dldataset.Get(datasetCategory, datasetName)
	if err != nil {
		return err
	}
	defer dataset.Close()

	err = dataset.Download(ctx)
	if err != nil {
		return err
	}

	files, err := dataset.List(ctx)
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

	preprocessOptions, err := predictor.GetPreprocessOptions()
	if err != nil {
		return err
	}

	fileParts := dl.PartitionStringList(files, partitionListSize)

	outputs := make(chan interface{}, DefaultChannelBuffer)
	partlabels := map[string]string{}

	log.WithField("model", modelName).
		WithField("files_length", len(files)).
		WithField("file_parts_length", len(fileParts)).
		WithField("partition_list_size", partitionListSize).
		WithField("num_file_parts", numFileParts).
		WithField("using_gpu", useGPU).
		Info("starting inference on a dataset")

	if numFileParts == -1 {
		numFileParts = len(fileParts)
	}

	fileCnt := len(files)

	log.WithField("dataset_category", datasetCategory).
		WithField("dataset_name", datasetName).
		Info("using the specified dataset")

	inferenceProgress := dlcmd.NewProgress("inferring", fileCnt)

	for _, part := range fileParts[:numFileParts] {
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

		opts := []pipeline.Option{pipeline.ChannelBuffer(DefaultChannelBuffer)}
		output := pipeline.New(opts...).
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

	// inferenceProgress.FinishPrint("inference complete")
	inferenceProgress.Finish()

	close(outputs)

	if publishEvaluation == false {
		for range outputs {
		}
		log.WithField("model", modelName).Info("publishEvaluation is false")
		return nil
	}

	cnt := 0
	cntTop1 := 0
	cntTop5 := 0

	databaseInsertProgress := dlcmd.NewProgress("inserting prediction", batchSize)

	for out0 := range outputs {
		if cnt > fileCnt {
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

			err = inputPredictionsTable.Insert(inputPrediction)
			if err != nil {
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

	databaseInsertProgress.Finish()
	log.WithField("model", modelName).Info("finised inserting prediction")

	modelAccuracy := evaluation.ModelAccuracy{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		Top1:      float64(cntTop1) / float64(fileCnt),
		Top5:      float64(cntTop5) / float64(fileCnt),
	}

	pp.Println("Top1 = ", float64(cntTop1)/float64(fileCnt))
	pp.Println("Top5 = ", float64(cntTop5)/float64(fileCnt))

	if err := modelAccuracyTable.Insert(modelAccuracy); err != nil {
		log.WithError(err).Error("failed to publish model accuracy entry")
	}

	log.WithField("model", modelName).Info("downloading trace information")

	performance := evaluation.Performance{
		ID:         bson.NewObjectId(),
		CreatedAt:  time.Now(),
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

var predictDatasetCmd = &cobra.Command{
	Use:     "dataset",
	Short:   "Evaluate the dataset using the specified model and framework",
	Aliases: []string{"dataset"},
	RunE: func(c *cobra.Command, args []string) error {
		if modelName == "all" {
			for _, model := range framework.Models() {
				modelName = model.Name
				modelVersion = model.Version
				runPredictDatasetCmd(c, args)
			}
			return nil
		}
		return runPredictDatasetCmd(c, args)
	},
}

func init() {
	predictDatasetCmd.PersistentFlags().StringVar(&datasetCategory, "dataset_category", "vision", "the dataset category to use for prediction")
	predictDatasetCmd.PersistentFlags().StringVar(&datasetName, "dataset_name", "ilsvrc2012_validation", "the name of the dataset to perform the evaluations on. When using `ilsvrc2012_validation`, optimized versions of the dataset are used when the input network takes 224 or 227")
	predictDatasetCmd.PersistentFlags().BoolVar(&centerCrop, "center_crop", true, "determines whether the dataset used should be center cropped")
	predictDatasetCmd.PersistentFlags().IntVar(&numFileParts, "num_file_parts", -1, "the number of file parts to process. Setting file parts to a value other than -1 means that only the first num_file_parts * batch_size images are infered from the dataset. This is useful while performing performance evaluations, where only a few hundred evaluation samples are useful")
}
