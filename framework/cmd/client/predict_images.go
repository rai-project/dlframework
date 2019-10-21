package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/registryquery"
	rgrpc "github.com/rai-project/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var imagesCmd = &cobra.Command{
	Use:     "images",
	Short:   "Request MLModelScope agents to predict images",
	Aliases: []string{"image"},
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("images dir path needs to be provided")
		}
		imgDir, _ := filepath.Abs(args[0])

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
				BatchSize:    int32(batchSize),
				FeatureLimit: int32(featureLimit),
			},
		})
		if err != nil {
			return errors.Wrap(err, "unable to open the predictor")
		}

		defer client.Close(ctx, &dlframework.PredictorCloseRequest{
			Predictor: predictor,
		})

		var data [][]byte
		err = filepath.Walk(imgDir, func(path string, info os.FileInfo, err error) error {
			if path == imgDir {
				return nil
			}

			reader, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(reader)
			data = append(data, buf.Bytes())

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "failed to walk image dir")
		}

		var imgs []*dlframework.Image
		for i, v := range data {
			imgs = append(imgs, &dlframework.Image{
				ID:   string(i),
				Data: v,
				// Preprocessed: false,
			})
		}

		imgsReq := dlframework.ImagesRequest{
			Predictor: predictor,
			Images:    imgs,
			Options: &dlframework.PredictionOptions{
				BatchSize:    int32(batchSize),
				FeatureLimit: int32(featureLimit),
			},
		}

		res, err := client.Images(ctx, &imgsReq)
		if err != nil {
			return errors.Wrap(err, "unable to get response from images request")
		}

		_ = res
		// pp.Println(res)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(imagesCmd)
}
