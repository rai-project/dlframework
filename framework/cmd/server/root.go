package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/VividCortex/robustly"
	"github.com/cockroachdb/cmux"
	"github.com/facebookgo/freeport"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/cmd"
	monitors "github.com/rai-project/monitoring/monitors"
	"github.com/rai-project/utils"
	echologger "github.com/rai-project/web/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//dllayer "github.com/rai-project/dllayer/cmd"

	_ "github.com/rai-project/logger/hooks"
)

var (
	local             bool
	profile           bool
	DefaultRunOptions = &robustly.RunOptions{
		RateLimit:  1,                   // the rate limit in crashes per second
		Timeout:    time.Second,         // the timeout (after which Run will stop trying)
		PrintStack: true,                // whether to print the panic stacktrace or not
		RetryDelay: 0 * time.Nanosecond, // inject a delay before retrying the run
	}
)

func freePort() (string, error) {
	port, err := freeport.Get()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(port), nil
}

func getHost() (string, error) {
	if local {
		return utils.GetLocalIp()
	}
	address, err := utils.GetExternalIp()
	if err != nil {
		return "", err
	}
	return address, nil
}

func RunRootE(c *cobra.Command, framework dlframework.FrameworkManifest, args []string) (<-chan bool, error) {

	done := make(chan bool)

	frameworkName := framework.GetName()
	port, found := os.LookupEnv("PORT")
	if !found {
		p, err := freePort()
		if err != nil {
			return done, err
		}
		port = p
	}

	host, found := os.LookupEnv("HOST")
	if !found {
		h, err := getHost()
		if err != nil {
			return done, err
		}
		host = h
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return done, errors.Wrapf(err, "⚠️ the port %s is not a valid integer", port)
	}

	predictors, err := agent.GetPredictors(framework)
	if err != nil {
		return done, errors.Wrapf(err,
			"⚠️ failed to get predictors for %s. make sure you have "+
				"imported the framework's predictor package",
			frameworkName,
		)
	}

	externalHost := host
	if e, ok := os.LookupEnv("EXTERNAL_HOST"); ok {
		externalHost = e
	}

	externalPort := port
	if p, ok := os.LookupEnv("EXTERNAL_PORT"); ok {
		externalPort = p
	}

	agnt, err := agent.New(predictors, agent.WithHost(externalHost), agent.WithPortString(externalPort))
	if err != nil {
		return done, err
	}

	registeryServer, err := agnt.RegisterManifests()
	if err != nil {
		return done, err
	}

	predictorServer, err := agnt.RegisterPredictor()
	if err != nil {
		return done, err
	}

	address := fmt.Sprintf("%s:%d", host, portInt)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return done, errors.Wrapf(err, "failed to listen on ip %s", address)
	}

	ctx := context.Background()

	if profile {
		// create the cmux object that will multiplex 2 protocols on same port
		m := cmux.New(lis)
		// match gRPC requests, otherwise regular HTTP requests
		grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		httpL := m.Match(cmux.Any())

		e := echo.New()
		e.Logger = &echologger.EchoLogger{Entry: log}
		monitors.AddRoutes(e)

		log.Debugf("➡️  "+frameworkName+" service is listening on %s", address)

		go func() {
			defer registeryServer.GracefulStop()
			registeryServer.Serve(grpcL)
			done <- true
		}()
		go func() {
			defer predictorServer.GracefulStop()
			predictorServer.Serve(grpcL)
			done <- true
		}()
		go func() {
			defer e.Shutdown(ctx)
			e.Listener = httpL
			err := e.Start(address)
			if err != nil {
				log.WithError(err).Error("failed to start echo server")
				return
			}
			done <- true
		}()

		log.Println("listening and serving (multiplexed) on", address)
		err = m.Serve()
		if err != nil {
			return nil, err
		}
		return done, nil
	}

	go func() {
		defer registeryServer.GracefulStop()
		registeryServer.Serve(lis)
		done <- true
	}()

	go func() {
		defer predictorServer.GracefulStop()
		log.Debugf("➡️  "+frameworkName+" service is listening on %s", address)
		predictorServer.Serve(lis)
		done <- true
	}()

	return done, nil
}

type FrameworkRegisterFunction func()

// represents the base command when called without any subcommands
func NewRootCommand(frameworkRegisterFunc FrameworkRegisterFunction, framework0 dlframework.FrameworkManifest) (*cobra.Command, error) {
	frameworkName := framework0.GetName()
	rootCmd := &cobra.Command{
		Use:   frameworkName + "-agent",
		Short: "Runs the carml " + frameworkName + " agent",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			frameworkRegisterFunc()
			framework = framework0
			//dllayer.Framework = framework
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			e := robustly.Run(
				func() {
					done, err := RunRootE(c, framework, args)
					if err != nil {
						panic("⚠️ " + err.Error())
					}
					<-done
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
		SetupFlags(rootCmd)
	})

	return rootCmd, nil
}

func SetupFlags(c *cobra.Command) {

	cobra.OnInitialize(initConfig)

	c.AddCommand(raicmd.VersionCmd)
	c.AddCommand(raicmd.LicenseCmd)
	c.AddCommand(raicmd.EnvCmd)
	c.AddCommand(raicmd.GendocCmd)
	c.AddCommand(raicmd.CompletionCmd)
	c.AddCommand(raicmd.BuildTimeCmd)

	c.AddCommand(predictCmd)
	c.AddCommand(downloadCmd)
	c.AddCommand(containerCmd)
	c.AddCommand(infoCmd)

	c.PersistentFlags().StringVar(&cmd.CfgFile, "config", "", "config file (default is $HOME/.carml_config.yaml)")
	c.PersistentFlags().BoolVarP(&cmd.IsVerbose, "verbose", "v", false, "Toggle verbose mode.")
	c.PersistentFlags().BoolVarP(&cmd.IsDebug, "debug", "d", false, "Toggle debug mode.")
	c.PersistentFlags().StringVarP(&cmd.AppSecret, "secret", "s", "", "The application secret.")
	c.PersistentFlags().BoolVarP(&local, "local", "l", false, "Listen on local address.")
	c.PersistentFlags().BoolVar(&profile, "profile", false, "Enable profile mode.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("app.secret", c.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", c.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", c.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/server")
	})

	cmd.Init()
	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.SetTimeout(time.Second * 1)
	shutdown.FirstFn(func() {})
	shutdown.SecondFn(func() {
		fmt.Println("🛑  shutting down!!")
		// tracer.Close()
	})
}
