package agent

import (
	"sync"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/predictor"
	"golang.org/x/sync/syncmap"
)

var (
	predictorServers struct {
		syncmap.Map
		sync.Mutex
	}
)

func GetPredictors(framework dl.FrameworkManifest) ([]predictor.Predictor, error) {
	name, err := framework.CanonicalName()
	if err != nil {
		return nil, err
	}
	val, ok := predictorServers.Load(name)
	if !ok {
		log.WithField("framework", framework.MustCanonicalName()).
			Warn("cannot find registered predictor server")
		return nil, errors.New("cannot find registered predictor server")
	}
	predictors, ok := val.([]predictor.Predictor)
	if !ok {
		log.WithField("framework", framework.MustCanonicalName()).
			Warn("invalid registered predictor server")
		return nil, errors.New("invalid predictor")
	}
	return predictors, nil
}

func GetPredictor(framework dl.FrameworkManifest) (predictor.Predictor, error) {
	predictors, err := GetPredictors(framework)
	if err != nil {
		return nil, err
	}
	if len(predictors) != 1 {
		return nil, errors.Errorf("expecting only one predictor but got %v", len(predictors))
	}
	return predictors[0], nil
}

func AddPredictor(framework dl.FrameworkManifest, pred predictor.Predictor) error {
	name, err := framework.CanonicalName()
	if err != nil {
		return err
	}

	var predictors []predictor.Predictor
	predictorServers.Lock()
	defer predictorServers.Unlock()

	val, ok := predictorServers.Load(name)
	if !ok {
		predictors = []predictor.Predictor{pred}
	} else {
		predictors = append(val.([]predictor.Predictor), pred)
	}
	predictorServers.Store(name, predictors)
	return nil
}

func PredictorFrameworks() []dl.FrameworkManifest {
	frameworks := []dl.FrameworkManifest{}
	predictorServers.Range(func(_ interface{}, val interface{}) bool {
		if predictor, ok := val.(predictor.Predictor); ok {
			framework, _, _ := predictor.Info()
			frameworks = append(frameworks, framework)
		}
		return true
	})
	return frameworks
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
