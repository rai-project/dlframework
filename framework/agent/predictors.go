package agent

import (
	"errors"

	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/predict"
	"golang.org/x/sync/syncmap"
)

var predictorServers syncmap.Map

func GetPredictor(framework dl.FrameworkManifest) (predict.Predictor, error) {
	key, err := framework.CanonicalName()
	if err != nil {
		return nil, err
	}
	val, ok := predictorServers.Load(key)
	if !ok {
		log.WithField("framework", key).
			Warn("cannot find registered predictor server")
		return nil, errors.New("cannot find registered predictor server")
	}
	predictor, ok := val.(predict.Predictor)
	if !ok {
		log.WithField("framework", key).
			Warn("invalid registered predictor server")
		return nil, errors.New("invalid predictor")
	}
	return predictor, nil
}

func AddPredictor(framework dl.FrameworkManifest, predictor predict.Predictor) error {
	key, err := framework.CanonicalName()
	if err != nil {
		return err
	}
	predictorServers.Store(key, predictor)
	return nil
}

func Predictors() []string {
	names := []string{}
	predictorServers.Range(func(key, _ interface{}) bool {
		if name, ok := key.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return names
}
