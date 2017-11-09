package server

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/pkg/errors"
	"github.com/rai-project/config"
	"github.com/rai-project/database"
	mongodb "github.com/rai-project/database/mongodb"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	udb "upper.io/db.v3"
)

var (
	sourceEvaluationID     string
	targetEvaluationID     string
	divergenceTollerance   float64
	divergenceReporterName string
)

type featurePair struct {
	sourceID       bson.ObjectId
	targetID       bson.ObjectId
	sourceInputID  string
	targetInputID  string
	sourceFeatures dlframework.Features
	targetFeatures dlframework.Features
}

func doComputeDivergence(
	db database.Database,
	evaluationCollection *evaluation.EvaluationCollection,
	inputPredictionCollection *evaluation.InputPredictionCollection,
	divergenceCollection *evaluation.DivergenceCollection,
	sourceEvaluationID string,
	targetEvaluationID string,
	reporterName string,
	divs ...func(pq *featurePair, reporter func(string, *featurePair, float64)),
) error {

	var bsonSourceEvaluationID, bsonTargetEvaluationID bson.ObjectId

	if bson.IsObjectIdHex(sourceEvaluationID) {
		bsonSourceEvaluationID = bson.ObjectIdHex(sourceEvaluationID)
	} else {
		bsonSourceEvaluationID = bson.ObjectId(sourceEvaluationID)
	}

	var sourceEvaluation evaluation.Evaluation
	err := evaluationCollection.FindOne(udb.Cond{"_id": bsonSourceEvaluationID}, &sourceEvaluation)
	if err != nil {
		return errors.Wrapf(err, "cannot find source evaluation with id = %v", bsonSourceEvaluationID.String())
	}

	if len(sourceEvaluation.InputPredictionIDs) == 0 {
		return errors.Errorf("empty source evaluation with id = %v", bsonSourceEvaluationID.String())
	}

	if bson.IsObjectIdHex(targetEvaluationID) {
		bsonTargetEvaluationID = bson.ObjectIdHex(targetEvaluationID)
	} else {
		bsonTargetEvaluationID = bson.ObjectId(targetEvaluationID)
	}

	var targetEvaluation evaluation.Evaluation
	err = evaluationCollection.FindOne(udb.Cond{"_id": bsonTargetEvaluationID}, &targetEvaluation)
	if err != nil {
		return errors.Wrapf(err, "cannot find target evaluation with id = %v", bsonTargetEvaluationID.String())
	}

	if sourceEvaluation.Model.Name != sourceEvaluation.Model.Name {
		return nil
	}

	if len(targetEvaluation.InputPredictionIDs) == 0 {
		return errors.Errorf("empty target evaluation with id = %v", bsonTargetEvaluationID.String())
	}

	if len(sourceEvaluation.InputPredictionIDs) != len(targetEvaluation.InputPredictionIDs) {
		return errors.Errorf("input prediction length mismatch %v != %v", len(sourceEvaluation.InputPredictionIDs), len(targetEvaluation.InputPredictionIDs))
	}

	numEvals := len(sourceEvaluation.InputPredictionIDs)

	progress := newProgress("checking prediction divergence step", numEvals)
	defer progress.FinishPrint("finished prediction divergence step")

	var wg sync.WaitGroup
	wg.Add(numEvals)

	numCPUs := runtime.NumCPU()

	pool, _ := tunny.CreatePool(2*numCPUs, func(o interface{}) interface{} {
		defer progress.Increment()
		defer wg.Done()
		ii := o.(int)
		sourcePredictionID := sourceEvaluation.InputPredictionIDs[ii]
		targetPredictionID := targetEvaluation.InputPredictionIDs[ii]

		var sourcePrediction evaluation.InputPrediction
		err := inputPredictionCollection.FindOne(sourcePredictionID, &sourcePrediction)
		if err != nil {
			log.WithError(err).Errorf("cannot find source prediction with id = %v", sourcePredictionID)
			return nil
		}

		var targetPrediction evaluation.InputPrediction
		err = inputPredictionCollection.FindOne(targetPredictionID, &targetPrediction)
		if err != nil {
			log.WithError(err).Errorf("cannot find target prediction with id = %v", targetPredictionID)
			return nil
		}

		sourceFeatures := sourcePrediction.Features
		targetFeatures := targetPrediction.Features

		reporter := func(name string, pair *featurePair, divergence float64) {}
		switch reporterName {
		case "print", "Print":
			reporter = func(name string, pair *featurePair, divergence float64) {
				println("source_input_id= ", pair.sourceInputID, "target_input_id= ", pair.targetInputID, name, " divergence=", divergence)
			}
		case "database":
			reporter = func(name string, pair *featurePair, divergence float64) {
				divergenceCollection.Insert(evaluation.Divergence{
					ID:                           bson.NewObjectId(),
					CreatedAt:                    time.Now(),
					Method:                       name,
					Value:                        divergence,
					SourcePredictionID:           pair.sourceID,
					TargetPredictionID:           pair.targetID,
					SourceInputPredictionInputID: pair.sourceInputID,
					TargetInputPredictionInputID: pair.targetInputID,
					SourceFeatures:               pair.sourceFeatures,
					TargetFeatures:               pair.targetFeatures,
				})
			}
		}

		for _, div := range divs {
			div(
				&featurePair{
					sourceID:       sourcePredictionID,
					targetID:       targetPredictionID,
					sourceInputID:  sourcePrediction.InputID,
					targetInputID:  targetPrediction.InputID,
					sourceFeatures: sourceFeatures,
					targetFeatures: targetFeatures,
				},
				reporter,
			)
		}

		return nil
	}).Open()
	defer pool.Close()

	for ii := range sourceEvaluation.InputPredictionIDs {
		pool.SendWorkAsync(ii, nil)
	}
	wg.Wait()

	return nil
}

func computeDivergence(c *cobra.Command, args []string, divs ...func(pq *featurePair, reporter func(string, *featurePair, float64))) error {
	opts := []database.Option{}
	if len(databaseEndpoints) != 0 {
		opts = append(opts, database.Endpoints(databaseEndpoints))
	}
	db, err := mongodb.NewDatabase(databaseName, opts...)
	defer db.Close()

	evaluationCollection, err := evaluation.NewEvaluationCollection(db)
	if err != nil {
		return err
	}

	inputPredictionCollection, err := evaluation.NewInputPredictionCollection(db)
	if err != nil {
		return err
	}

	divergenceCollection, err := evaluation.NewDivergenceCollection(db)
	if err != nil {
		return err
	}

	sources := []string{}
	if sourceEvaluationID == "all" {
		srcs := []evaluation.Evaluation{}
		err = evaluationCollection.FindAll(nil, &srcs)
		if err != nil {
			return err
		}
		for _, src := range srcs {
			id := string(src.ID)
			sources = append(sources, id)
		}
	} else {
		sources = append(sources, sourceEvaluationID)
	}

	targets := []string{}
	if targetEvaluationID == "all" {
		trgts := []evaluation.Evaluation{}
		err = evaluationCollection.FindAll(nil, &trgts)
		if err != nil {
			return err
		}
		for _, trgt := range trgts {
			id := string(trgt.ID)
			targets = append(targets, id)
		}
	} else {
		targets = append(targets, sourceEvaluationID)
	}

	for _, src := range sources {
		for _, trgt := range targets {
			if src == trgt {
				continue
			}
			doComputeDivergence(
				db,
				evaluationCollection,
				inputPredictionCollection,
				divergenceCollection,
				src,
				trgt,
				divergenceReporterName,
				divs...,
			)
		}
	}

	return nil
}

func divergencePreRun(c *cobra.Command, args []string) {
	if databaseName == "" {
		databaseName = config.App.Name
	}
	if sourceEvaluationID == "" && targetEvaluationID == "" && len(args) >= 2 {
		sourceEvaluationID = args[0]
		targetEvaluationID = args[1]
	}

	if databaseAddress != "" {
		databaseEndpoints = []string{databaseAddress}
	}
}

var (
	divergenceDispatch = map[string]func(pair *featurePair, reporter func(string, *featurePair, float64)){
		"Bhattacharyya": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.Bhattacharyya(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform Bhattacharyya")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("Bhattacharyya", pair, divergence)
			}
		},

		"Hellinger": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.Hellinger(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform Hellinger")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("Hellinger", pair, divergence)
			}
		},

		"Correlation": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.Correlation(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform Correlation")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("Correlation", pair, divergence)
			}
		},

		"JensenShannon": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.JensenShannon(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform JensenShannon")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("JensenShannon", pair, divergence)
			}
		},

		"Covariance": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.Covariance(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform Covariance")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("Covariance", pair, divergence)
			}
		},

		"KullbackLeibler": func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			divergence, err := pair.sourceFeatures.KullbackLeiblerDivergence(pair.targetFeatures)
			if err != nil {
				log.WithError(err).Error("cannot perform KullbackLeiblerDivergence")
				return
			}
			if math.Abs(divergence) >= divergenceTollerance && divergence != 0 {
				reporter("KullbackLeiblerDivergence", pair, divergence)
			}
		},
	}
)

var databaseKLDivergenceCmd = &cobra.Command{
	Use:     "kldivergence",
	Aliases: []string{"kl", "KullbackLeibler"},
	Short:   "Perform Kullback-Leibler divergence on two evaluation ids",
	Long:    `for example : go run mxnet.go database kldivergence --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["KullbackLeibler"])
	},
}

var databaseJSDivergenceCmd = &cobra.Command{
	Use:     "jensenshannon",
	Aliases: []string{"js", "JensenShannon"},
	Short:   "Perform JensenShannon divergence on two evaluation ids",
	Long:    `for example : go run mxnet.go database jensenshannon --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["JensenShannon"])
	},
}

var databaseCovDivergenceCmd = &cobra.Command{
	Use:     "covariance",
	Aliases: []string{"cov", "Covariance"},
	Short:   "Perform Covariance divergence on two evaluation ids",
	Long:    `for example : go run mxnet.go database covariance --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["Covariance"])
	},
}

var databaseCorrDivergenceCmd = &cobra.Command{
	Use:     "correlation",
	Aliases: []string{"cor", "corr", "Correlation"},
	Short:   "Perform Correlation divergence on two evaluation ids",
	Long:    `for example : go run mxnet.go database correlation --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["Correlation"])
	},
}

var databaseHellDivergenceCmd = &cobra.Command{
	Use:     "hellinger",
	Aliases: []string{"hel", "hell", "Hellinger"},
	Short:   "Perform Correlation divergence on two evaluation ids",
	Long:    `for example : go run mxnet.go database hellinger --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["Hellinger"])
	},
}

var databaseBhattDivergenceCmd = &cobra.Command{
	Use:     "bhattacharyya",
	Aliases: []string{"bhat", "bhatt", "Bhattacharyya"},
	Short:   "Perform Correlation bhattacharyya on two evaluation ids",
	Long:    `for example : go run mxnet.go database bhattacharyya --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun:  divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, divergenceDispatch["Bhattacharyya"])
	},
}

var databaseDivergenceCmd = &cobra.Command{
	Use:    "divergence",
	Short:  "Perform Kullback-Leibler, JensenShannon, Covariance, Correlation, Hellinger, and Bhattacharyya divergence on two evaluation ids",
	Long:   `for example : go run mxnet.go database divergence --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun: divergencePreRun,
	RunE: func(c *cobra.Command, args []string) error {
		return computeDivergence(c, args, func(pair *featurePair, reporter func(string, *featurePair, float64)) {
			for _, f := range divergenceDispatch {
				f(pair, reporter)
			}
		})
	},
}

var divergenceCmds = []*cobra.Command{
	databaseKLDivergenceCmd,
	databaseJSDivergenceCmd,
	databaseCovDivergenceCmd,
	databaseCorrDivergenceCmd,
	databaseHellDivergenceCmd,
	databaseBhattDivergenceCmd,
	databaseDivergenceCmd,
}

func init() {
	for _, cmd := range divergenceCmds {
		cmd.PersistentFlags().StringVar(&sourceEvaluationID, "source", "", "source id for the evaluation")
		cmd.PersistentFlags().StringVar(&targetEvaluationID, "target", "", "target id for the evaluation")
		cmd.PersistentFlags().Float64Var(&divergenceTollerance, "tollerance", 0.01, "tolerance to use while printing divergence information")
		cmd.PersistentFlags().StringVar(&divergenceReporterName, "reporter", "print", "method to use to report divergence")
	}
}
