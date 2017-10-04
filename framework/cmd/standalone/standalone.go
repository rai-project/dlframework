package main

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/VividCortex/robustly"
	"github.com/k0kubun/pp"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/pkg/errors"
	_ "github.com/rai-project/caffe/predict"
	_ "github.com/rai-project/caffe2/predict"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/cmd/server"
	_ "github.com/rai-project/mxnet/predict"
	_ "github.com/rai-project/tensorflow/predict"
	"github.com/rai-project/tracer"
	"github.com/spf13/cobra"
)

// represents the base command when called without any subcommands
func NewRootCommand() (*cobra.Command, error) {
	frameworks := agent.PredictorFrameworks()
	rootCmd := &cobra.Command{
		Use:   "all-agents",
		Short: "Runs the carml " + frameworks + " agent",
		RunE: func(c *cobra.Command, args []string) error {
			e := robustly.Run(
				func() {
					anyDone := make(chan struct{})
					for _, framework := range frameworks {
						done, err := RunRootE(c, framework, args)
						if err != nil {
							panic("⚠️ " + err.Error())
						}
						go anyDone <- done
					}
					<-anyDone
				},
				DefaultRunOptions,
			)
			if e != 0 {
				return errors.Errorf("⚠️ %s has panniced %d times ... giving up", frameworkName+"-agent", e)
			}
			return nil
		},
	}
	var once sync.Once
	once.Do(func() {
		server.SetupFlags(rootCmd)
	})
	return rootCmd, nil
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.SetTimeout(time.Second * 1)
	shutdown.SecondFn(func() {
		pp.Println("🛑 shutting down!!")
		tracer.Close()
	})
}
