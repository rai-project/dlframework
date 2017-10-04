package agent

import (
	"errors"

	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/predict"
	"golang.org/x/sync/syncmap"
)

var predictorServers syncmap.Map

func GetPredictor(framework dl.FrameworkManifest) (predict.Predictor, error) {
	val, ok := predictorServers.Load(framework)
	if !ok {
		log.WithField("framework", framework.MustCanonicalName()).
			Warn("cannot find registered predictor server")
		return nil, errors.New("cannot find registered predictor server")
	}
	predictor, ok := val.(predict.Predictor)
	if !ok {
		log.WithField("framework", framework.MustCanonicalName()).
			Warn("invalid registered predictor server")
		return nil, errors.New("invalid predictor")
	}
	return predictor, nil
}

func AddPredictor(framework dl.FrameworkManifest, predictor predict.Predictor) error {
	predictorServers.Store(framework, predictor)
	return nil
}

func PredictorFrameworks() []dl.FrameworkManifest {
	frameworks := []dl.FrameworkManifest{}
	predictorServers.Range(func(key, _ interface{}) bool {
		if framework, ok := key.(dl.FrameworkManifest); ok {
			frameworks = append(frameworks, framework)
		}
		return true
	})
	return frameworks
}

func Predictors() []string {
	names := []string{}
	predictorServers.Range(func(key, _ interface{}) bool {
		if framework, ok := key.(dl.FrameworkManifest); ok {
			names = append(names, framework.MustCanonicalName())
		}
		return true
	})
	return names
}
