package client

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/k0kubun/pp"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/registryquery"
	rgrpc "github.com/rai-project/grpc"
	"github.com/rai-project/tracer"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"

	"google.golang.org/grpc"
)

var urlsCmd = &cobra.Command{
	Use:     "urlsCmd",
	Short:   "urlsCmd",
	Aliases: []string{"urls", "url"},
	Long:    `urlsCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("urlsfile path needs to be provided")
		}
		urlsFile, _ := filepath.Abs(args[0])

		if len(args) > 1 {
			batchSize, _ = strconv.Atoi(args[1])
		}

		var outputDir string
		if len(args) > 2 {
			outputDir = args[2]
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				return errors.Wrap(err, "error creating output dir")

			}
		}

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

		// fill the batch with the same image
		if len(data) == 1 {
			for ii := 0; ii < batchSize; ii++ {
				data = append(data, data[0])
			}
		}

		cleanNames()
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
		span, ctx := tracer.StartSpanFromContext(ctx, "urls")
		spanClosed := false
		defer func() {
			if !spanClosed {
				span.Finish()
			}
		}()

		conn, err := rgrpc.DialContext(
			ctx,
			dlframework.PredictServiceDescription,
			serverAddress,
			grpc.WithInsecure(),
		)
		if err != nil {
			return errors.Wrapf(err, "unable to dial %s", serverAddress)
		}
		defer conn.Close()

		client := dlframework.NewPredictClient(conn)

		ExecutionOptions := &dlframework.ExecutionOptions{
			DeviceCount: map[string]int32{"gpu": 0},
		}

		predictor, err := client.Open(ctx, &dlframework.PredictorOpenRequest{
			ModelName:        modelName,
			ModelVersion:     modelVersion,
			FrameworkName:    frameworkName,
			FrameworkVersion: frameworkVersion,
			Options: &dlframework.PredictionOptions{
				BatchSize:        uint32(batchSize),
				ExecutionOptions: ExecutionOptions,
			},
		})
		if err != nil {
			return errors.Wrap(err, "unable to open the predictor")
		}

		defer client.Close(ctx, predictor)

		var urls []*dlframework.URLsRequest_URL
		for i, url := range data {
			urls = append(urls, &dlframework.URLsRequest_URL{
				ID:   string(i),
				Data: url,
			})
		}

		urlsReq := dlframework.URLsRequest{
			Predictor: predictor,
			Urls:      urls,
			Options: &dlframework.PredictionOptions{
				BatchSize: uint32(batchSize),
			},
		}

		res, err := client.URLs(ctx, &urlsReq)
		if err != nil {
			return errors.Wrap(err, "unable to get response from urls request")
		}
		_ = res

		// tag res in tracing
		fr := dlframework.Features(res.Responses[0].GetFeatures())
		fr.Sort()
		for ii, f := range fr[:5] {
			span.SetTag("predictions_"+strconv.Itoa(ii), struct {
				Pos   int
				Index int64
				Name  string
				Prob  float32
			}{
				Pos:   ii,
				Index: f.GetIndex(),
				Name:  f.GetName(),
				Prob:  f.GetProbability(),
			})
		}

		spanClosed = true
		span.Finish()

		pp.Println(args)
		pp.Println(modelName, modelVersion)

		time.Sleep(10 * time.Second)

		traceID := span.Context().(jaeger.SpanContext).TraceID()
		query := fmt.Sprintf("http://localhost:16686/api/traces/%v", strconv.FormatUint(traceID.Low, 16))
		resp, err := grequests.Get(query, nil)

		outputfile := frameworkName + "_" + frameworkVersion + "_" + modelName + "_" + modelVersion + "_" + strconv.Itoa(batchSize)
		err = ioutil.WriteFile(filepath.Join(outputDir, outputfile), []byte(resp.String()), 0644)
		if err != nil {
			return errors.Wrap(err, "unable to write to file")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(urlsCmd)
}
