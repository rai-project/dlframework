package client

import (
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/spf13/cobra"
)

type Framework struct {
	FrameworkName    string
	FrameworkVersion string
}

type Model struct {
	ModelName    string
	ModelVersion string
}

// represents the base command when called without any subcommands
func NewRootCommand(framework Framework, model Model, serverAddress string, data []string) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "carml client",
		Short: "Runs the carml client",
		RunE: func(cmd *cobra.Command, args []string) error {
			frameworkName := strings.ToLower(framework.FrameworkName)
			frameworkVersion := strings.ToLower(framework.FrameworkVersion)
			modelName := strings.ToLower(model.ModelName)
			modelVersion := strings.ToLower(model.ModelVersion)

			ctx := context.Background()

			conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithInsecure())
			if err != nil {
				return errors.Wrapf(err, "unable to dial %s", serverAddress)
			}
			defer conn.Close()
			client := dlframework.NewPredictClient(conn)

			predictor, err := client.Open(ctx, &dlframework.PredictorOpenRequest{
				ModelName:        modelName,
				ModelVersion:     modelVersion,
				FrameworkName:    frameworkName,
				FrameworkVersion: frameworkVersion,
			})
			if err != nil {
				return errors.Wrap(err, "unable to open the predictor")
			}
			defer client.Close(ctx, predictor)

			var urls []*dlframework.URLsRequest_URL
			for i, url := range data {
				urls = append(urls, &dlframework.URLsRequest_URL{
					Id:   string(i),
					Data: url,
				})
			}
			urlReq := dlframework.URLsRequest{
				Predictor: predictor,
				Urls:      urls,
			}
			res, err := client.URLs(ctx, &urlReq)
			if err != nil {
				return errors.Wrap(err, "unable to get response from urls request")
			}

			pp.Println(res)

			return nil
		},
	}
	return rootCmd, nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
