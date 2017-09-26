package client

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/k0kubun/pp"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/registryquery"
	rgrpc "github.com/rai-project/grpc"
	"github.com/spf13/cobra"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger/model"

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

		spanClosed = true
		span.Finish()

		// toHex := func(t jaeger.TraceID) string {
		// 	return fmt.Sprintf("%v", t.Low)
		// }

		time.Sleep(10 * time.Second)

		traceIDModel := span.Context().(jaeger.SpanContext).TraceID()
		query := fmt.Sprintf("http://localhost:16686/api/traces/%v", traceIDModel.String())
		resp, err := grequests.Get(query, nil)
		pp.Println(query, "    ", resp.String())

		_ = res

		// pp.Println(res)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(urlsCmd)
}

// TraceID is a serializable form of model.TraceID
type TraceID [16]byte

func TraceIDFromDomain(high, low uint64) TraceID {
	dbTraceID := TraceID{}
	binary.BigEndian.PutUint64(dbTraceID[:8], high)
	binary.BigEndian.PutUint64(dbTraceID[8:], low)
	return dbTraceID
}

// ToDomain converts trace ID from db-serializable form to domain TradeID
func (dbTraceID TraceID) ToDomain() model.TraceID {
	traceIDHigh := binary.BigEndian.Uint64(dbTraceID[:8])
	traceIDLow := binary.BigEndian.Uint64(dbTraceID[8:])
	return model.TraceID{High: traceIDHigh, Low: traceIDLow}
}

// String returns hex string representation of the trace ID.
func (dbTraceID TraceID) String() string {
	return string(dbTraceID[:])
}
