package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/VividCortex/robustly"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/pkg/errors"
	_ "github.com/rai-project/caffe/predict"
	_ "github.com/rai-project/tensorrt/predict"
	//_ "github.com/rai-project/caffe2/predict"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/dlframework/framework/cmd/server"
	_ "github.com/rai-project/mxnet/predict"
	_ "github.com/rai-project/tensorflow/predict"
	"github.com/rai-project/tracer"
	"github.com/spf13/cobra"
)

// represents the base command when called without any subcommands
func NewRootCommand() (*cobra.Command, error) {
	frameworks := agent.PredictorFrameworks()
	if len(frameworks) == 0 {
		return nil, errors.New("no frameworks found")
	}
	frameworkNames := make([]string, len(frameworks))
	for ii, framework := range frameworks {
		frameworkNames[ii] = framework.MustCanonicalName()
	}
	rootCmd := &cobra.Command{
		Use:   "all-agents",
		Short: "Runs the carml " + strings.Join(frameworkNames, ", ") + " agent",
		RunE: func(c *cobra.Command, args []string) error {
			e := robustly.Run(
				func() {
					anyDone := make(chan bool)
					for _, framework := range frameworks {
						go func() {
							done, err := server.RunRootE(c, framework, args)
							if err != nil {
								panic("⚠️  " + err.Error())
							}
							v := <-done
							anyDone <- v
						}()
					}
					<-anyDone
				},
				server.DefaultRunOptions,
			)
			if e != 0 {
				return errors.Errorf("⚠️  %s has panniced %d times ... giving up", strings.Join(frameworkNames, ", ")+"-agent", e)
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
	cmd.Init()
	rootCmd, err := NewRootCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.SetTimeout(time.Second * 1)
	shutdown.SecondFn(func() {
		fmt.Println("🛑  shutting down!!")
		tracer.Close()
	})
}
