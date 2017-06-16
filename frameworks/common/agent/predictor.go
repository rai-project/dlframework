package agent

import (
	dl "github.com/rai-project/dlframework"
	context "golang.org/x/net/context"
)

type Predictor struct {
	Base
}

func (p *Predictor) Predict(ctx context.Context, req *dl.PredictRequest) (*dl.PredictResponse, error) {
	return nil, nil
}

func (p *Predictor) PublishInRegistery() error {
	return p.Base.PublishInRegistery("predictor")
}
