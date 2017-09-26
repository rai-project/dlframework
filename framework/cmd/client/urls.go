package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/registryquery"
	rgrpc "github.com/rai-project/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var urlsCmd = &cobra.Command{
	Use:   "urlsCmd",
	Short: "urlsCmd",
	Long:  `urlsCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("urlsfile path needs to be provided")
		}
		urlsFile, _ := filepath.Abs(args[0])

		agents, err := registryquery.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
		if err != nil {
			return err
		}
		if len(agents) == 0 {
			return errors.Errorf("no agent found for %v:%v/%v:%v", frameworkName, frameworkVersion, modelName, modelVersion)
		}

		agent := agents[0]
		serverAddress := fmt.Sprintf("%s:%s", agent.Host, agent.Port)

		ctx := context.Background()

		conn, err := rgrpc.DialContext(ctx, dlframework.PredictServiceDescription, serverAddress, grpc.WithInsecure())
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
			Options: &dlframework.PredictionOptions{
				BatchSize: uint32(batchSize),
			},
		})
		if err != nil {
			return errors.Wrap(err, "unable to open the predictor")
		}

		defer client.Close(ctx, predictor)

		var data []string
		f, err := os.Open(urlsFile)
		if err != nil {
			return errors.Wrapf(err, "cannot read %s", urlsFile)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			data = append(data, line)
		}

		var urls []*dlframework.URLsRequest_URL

		for i, url := range data {
			urls = append(urls, &dlframework.URLsRequest_URL{
				ID:   string(i),
				Data: url,
			})
		}
		urlReq := dlframework.URLsRequest{
			Predictor: predictor,
			Urls:      urls,
			Options: &dlframework.PredictionOptions{
				BatchSize: uint32(batchSize),
			},
		}

		res, err := client.URLs(ctx, &urlReq)
		if err != nil {
			return errors.Wrap(err, "unable to get response from urls request")
		}

		_ = res
		pp.Println(res)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(urlsCmd)
}
