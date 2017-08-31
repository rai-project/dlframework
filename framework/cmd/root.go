package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Unknwon/com"
	"github.com/facebookgo/freeport"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/agent"
	"github.com/rai-project/logger"
	"github.com/rai-project/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isDebug   bool
	isVerbose bool
	local     bool
	appSecret string
	cfgFile   string

	log *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd")
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
		RunE: func(cmd *cobra.Command, args []string) error {
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
				h, err = getHost()
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

			agnt, err := agent.New(predictor, agent.WithHost(host), agent.WithPort(portInt))
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

func setupFlags(cmd *cobra.Command) {

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(raicmd.VersionCmd)
	cmd.AddCommand(raicmd.LicenseCmd)
	cmd.AddCommand(raicmd.EnvCmd)
	cmd.AddCommand(raicmd.GendocCmd)
	cmd.AddCommand(raicmd.CompletionCmd)
	cmd.AddCommand(raicmd.BuildTimeCmd)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.carml_config.yaml)")
	cmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	cmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")
	cmd.PersistentFlags().StringVarP(&appSecret, "secret", "s", "", "The application secret.")
	cmd.PersistentFlags().BoolVarP(&local, "local", "l", false, "Listen on local address.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("app.secret", cmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", cmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", cmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	log.Level = logrus.DebugLevel
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "dlframework/framework/cmd")
	})

	color.NoColor = false
	opts := []config.Option{
		config.AppName("carml"),
		config.ColorMode(true),
	}
	if c, err := homedir.Expand(cfgFile); err == nil {
		cfgFile = c
	}
	if config.IsValidRemotePrefix(cfgFile) {
		opts = append(opts, config.ConfigRemotePath(cfgFile))
	} else if com.IsFile(cfgFile) {
		if c, err := filepath.Abs(cfgFile); err == nil {
			cfgFile = c
		}
		opts = append(opts, config.ConfigFileAbsolutePath(cfgFile))
	}

	if appSecret != "" {
		opts = append(opts, config.AppSecret(appSecret))
	}
	config.Init(opts...)
}
