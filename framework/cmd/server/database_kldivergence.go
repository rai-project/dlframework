package server

import (
	"math"

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
	sourceEvaluationID   string
	targetEvaluationID   string
	divergenceTollerance float64
)

type featurePair struct {
	sourceInputID  string
	targetInputID  string
	sourceFeatures dlframework.Features
	targetFeatures dlframework.Features
}

var databaseKLDivergenceCmd = &cobra.Command{
	Use:   "kldivergence",
	Short: "Perform Kullback-Leibler divergence on two evaluation ids",
	Long:  `for example : go run mxnet.go database kldivergence --database_address=minsky1-1.csl.illinois.edu --database_name=carml --source=5a01fc48ca60cc797e63603c --target=5a0203f8ca60ccd42aa2a706`,
	PreRun: func(c *cobra.Command, args []string) {
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
	},
	RunE: func(c *cobra.Command, args []string) error {

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

		var bsonSourceEvaluationID, bsonTargetEvaluationID bson.ObjectId

		if bson.IsObjectIdHex(sourceEvaluationID) {
			bsonSourceEvaluationID = bson.ObjectIdHex(sourceEvaluationID)
		} else {
			bsonSourceEvaluationID = bson.ObjectId(sourceEvaluationID)
		}

		var sourceEvaluation evaluation.Evaluation
		err = evaluationCollection.FindOne(udb.Cond{"_id": bsonSourceEvaluationID}, &sourceEvaluation)
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

		if len(targetEvaluation.InputPredictionIDs) == 0 {
			return errors.Errorf("empty target evaluation with id = %v", bsonTargetEvaluationID.String())
		}

		if len(sourceEvaluation.InputPredictionIDs) != len(targetEvaluation.InputPredictionIDs) {
			return errors.Errorf("input prediction length mismatch %v != %v", len(sourceEvaluation.InputPredictionIDs), len(targetEvaluation.InputPredictionIDs))
		}

		featureChan := make(chan *featurePair, 1000)

		for ii := range sourceEvaluation.InputPredictionIDs {
			go func(ii int) {
				sourcePredictionID := sourceEvaluation.InputPredictionIDs[ii]
				targetPredictionID := targetEvaluation.InputPredictionIDs[ii]

				var sourcePrediction evaluation.InputPrediction
				err := inputPredictionCollection.FindOne(evaluation.InputPrediction{ID: sourcePredictionID}, &sourcePrediction)
				if err != nil {
					log.WithError(err).Errorf("cannot find source prediction with id = %v", sourcePredictionID)
					return
				}

				var targetPrediction evaluation.InputPrediction
				err = inputPredictionCollection.FindOne(evaluation.InputPrediction{ID: targetPredictionID}, &targetPrediction)
				if err != nil {
					log.WithError(err).Errorf("cannot find target prediction with id = %v", targetPredictionID)
					return
				}

				sourceFeatures := sourcePrediction.Features
				targetFeatures := targetPrediction.Features

				featureChan <- &featurePair{
					sourceInputID:  sourcePrediction.InputID,
					targetInputID:  targetPrediction.InputID,
					sourceFeatures: sourceFeatures,
					targetFeatures: targetFeatures,
				}
			}(ii)
		}

		for pair := range featureChan {
			divergence, err := pair.sourceFeatures.KullbackLeiblerDivergence(pair.targetFeatures)
			if err != nil {
				return errors.Wrapf(err, "cannot perform KullbackLeiblerDivergence")
			}
			if math.Abs(divergence) >= divergenceTollerance {
				println("source_input_id= ", pair.sourceInputID, "target_input_id= ", pair.targetInputID, " divergence=", divergence)
			}
		}

		return nil
	},
}

func init() {
	databaseKLDivergenceCmd.PersistentFlags().StringVar(&sourceEvaluationID, "source", "", "source id for the evaluation")
	databaseKLDivergenceCmd.PersistentFlags().StringVar(&targetEvaluationID, "target", "", "target id for the evaluation")
	databaseKLDivergenceCmd.PersistentFlags().Float64Var(&divergenceTollerance, "tollerance", 0.01, "tolerance to use while printing divergence information")
}
