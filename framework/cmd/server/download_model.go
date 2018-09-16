package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/spf13/cobra"
)

func downloadModels(ctx context.Context) error {
	predictorFramework, err := agent.GetPredictor(framework)
	if err != nil {
		return errors.Wrapf(err,
			"⚠️ failed to get predictor for %s. make sure you have "+
				"imported the framework's predictor package",
			framework.MustCanonicalName(),
		)

	}
	models := framework.Models()
	pb := newProgress("download models", len(models))
	for _, model := range models {
		err := predictorFramework.Download(ctx, model)
		if err != nil {
			return errors.Wrapf(err, "failed to download %s model", model.MustCanonicalName())
		}
		pb.Increment()
	}
	pb.Finish()

	return nil
}

var downloadModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Download CarML models",
	RunE: func(c *cobra.Command, args []string) error {
		ctx := context.Background()
		return downloadModels(ctx)
	},
}
