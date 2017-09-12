package server

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/facebookgo/freeport"
	"github.com/pkg/errors"
	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	local bool
	log   *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/server")
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

// represents the base command when called without any subcommands
func NewRootCommand(framework dlframework.FrameworkManifest) (*cobra.Command, error) {
	frameworkName := framework.GetName()
	rootCmd := &cobra.Command{
		Use:   frameworkName + "-agent",
		Short: "Runs the carml " + frameworkName + " agent",
		RunE: func(c *cobra.Command, args []string) error {
			port, found := os.LookupEnv("PORT")
			if !found {
				p, err := freePort()
				if err != nil {
					return err
				}
				port = p
			}

			host, found := os.LookupEnv("HOST")
			if !found {
				h, err := getHost()
				if err != nil {
					return err
				}
				host = h
			}

			portInt, err := strconv.Atoi(port)
			if err != nil {
				return errors.Wrapf(err, "the port %s is not a valid integer", port)
			}

			predictor, err := agent.GetPredictor(framework)
			if err != nil {
				return errors.Wrapf(err,
					"failed to get predictor for %s. make sure you have "+
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

			agnt, err := agent.New(predictor, agent.WithHost(externalHost), agent.WithPortString(externalPort))
			if err != nil {
				return err
			}

			registeryServer, err := agnt.RegisterManifests()
			if err != nil {
				return err
			}

			predictorServer, err := agnt.RegisterPredictor()
			if err != nil {
				return err
			}

			address := fmt.Sprintf("%s:%d", host, portInt)

			lis, err := net.Listen("tcp", address)
			if err != nil {
				return errors.Wrapf(err, "failed to listen on ip %s", address)
			}

			log.Debugf(frameworkName+" service is listening on %s", address)

			defer registeryServer.GracefulStop()
			defer predictorServer.GracefulStop()

			go registeryServer.Serve(lis)
			predictorServer.Serve(lis)
			return nil
		},
	}
	setupFlags(rootCmd)
	return rootCmd, nil
}

func setupFlags(c *cobra.Command) {

	cobra.OnInitialize(initConfig)

	c.AddCommand(raicmd.VersionCmd)
	c.AddCommand(raicmd.LicenseCmd)
	c.AddCommand(raicmd.EnvCmd)
	c.AddCommand(raicmd.GendocCmd)
	c.AddCommand(raicmd.CompletionCmd)
	c.AddCommand(raicmd.BuildTimeCmd)

	c.PersistentFlags().StringVar(&cmd.CfgFile, "config", "", "config file (default is $HOME/.carml_config.yaml)")
	c.PersistentFlags().BoolVarP(&cmd.IsVerbose, "verbose", "v", false, "Toggle verbose mode.")
	c.PersistentFlags().BoolVarP(&cmd.IsDebug, "debug", "d", false, "Toggle debug mode.")
	c.PersistentFlags().StringVarP(&cmd.AppSecret, "secret", "s", "", "The application secret.")
	c.PersistentFlags().BoolVarP(&local, "local", "l", false, "Listen on local address.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("app.secret", c.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", c.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", c.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cmd.Init()
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/server")
	})
}
