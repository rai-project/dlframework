package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/VividCortex/robustly"
	"github.com/cockroachdb/cmux"
	"github.com/facebookgo/freeport"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	monitors "github.com/rai-project/monitoring/monitors"
	"github.com/rai-project/utils"
	echologger "github.com/rai-project/web/logger"
	"github.com/spf13/cobra"
)

var (
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
		log.Debugf("➡️  "+framework.Name+" service is listening on %s", address)
		predictorServer.Serve(lis)
		done <- true
	}()

	return done, nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the agent serving process",
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
			return errors.Errorf("⚠️ %s has panniced %d times ... giving up", framework.Name+"-agent", e)
		}
		return nil
	},
}
