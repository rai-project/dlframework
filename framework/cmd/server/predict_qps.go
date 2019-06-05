package server

import (
	"fmt"
	"math"
	"time"

	"github.com/spf13/cobra"
)

var predictQPSCmd = &cobra.Command{
	Use:     "max-qps",
	Short:   "Find the maximun qps of the system using the specified model, framework and workload",
	Aliases: []string{"qps", "findmaxqps", "maxqps"},
	RunE: func(c *cobra.Command, args []string) error {
		qpsLowerBound := 0.0
		qpsUpperBound := math.MaxFloat64

		iters := int64(0)
		relativeQpsTolerance := 0.01

		for (qpsUpperBound-qpsLowerBound)/qpsLowerBound > relativeQpsTolerance && iters < math.MaxInt64 {
			iters++
			targetQps := 0.0
			if qpsLowerBound == 0 && qpsUpperBound == math.MaxFloat64 {
				targetQps = qps
			} else if qpsUpperBound == math.MaxFloat64 {
				targetQps = 2 * qpsLowerBound
			} else {
				targetQps = (qpsLowerBound + qpsUpperBound) / 2
			}

			log.WithField("targetQps", targetQps).Debug("creating a new trace")

			trace, latency, err := computeLatency(targetQps)
			if err != nil {
				return err
			}
			traceQps := trace.QPS()
			if qpsLowerBound < traceQps && traceQps < qpsUpperBound {
				measuredLatency := latency

				fmt.Printf("qps = %v, latency = %v\n",
					traceQps,
					measuredLatency,
				)
				log.WithField("qps", traceQps).
					WithField("% latency", measuredLatency).
					Info("replayed trace")
				if measuredLatency > 100*time.Millisecond {
					qpsUpperBound = math.Min(qpsUpperBound, traceQps)
				} else {
					qpsLowerBound = math.Max(traceQps, qpsLowerBound)
				}
			}

			// fmt.Printf("qpsLowerBound = %v qpsUpperBound =%v traceQps =%v\n", qpsLowerBound, qpsUpperBound, traceQps)

			log.WithField("qpsUpperBound", qpsUpperBound).
				WithField("qpsLowerBound", qpsLowerBound).
				Debug("generated new trace")
		}

		fmt.Printf("Max QPS subject to %v ms %v latency bound = %v",
			latencyBound,
			100*latencyBoundPercentile,
			math.Min(qpsUpperBound, qpsLowerBound),
		)

		fmt.Printf("qps = %v\n", math.Min(qpsUpperBound, qpsLowerBound))
		return nil
	},
}

func init() {
	predictQPSCmd.PersistentFlags().Float64Var(&qps, "initial_qps", 8, "the initial QPS")
	predictQPSCmd.PersistentFlags().Float64Var(&latencyBoundPercentile, "percentile", 95, "the minimum percent of queries meeting the latency bound")
	predictQPSCmd.PersistentFlags().Int64Var(&minDuration, "min_duration", 1000, "the minimum duration of the trace in ms")
	predictQPSCmd.PersistentFlags().IntVar(&minQueries, "min_queries", 4096, "the minimum number of queries")
	predictQPSCmd.PersistentFlags().Int64Var(&latencyBound, "latency_bound", 100, "the target latency bound in ms")
}
