package client

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/cmd"
	"github.com/rai-project/dlframework/http"
	rgrpc "github.com/rai-project/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Framework struct {
	FrameworkName    string
	FrameworkVersion string
}

type Model struct {
	ModelName    string
	ModelVersion string
}

var (
	log *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/client")
)

// represents the base command when called without any subcommands
func NewRootCommand(framework Framework, model Model, data []string) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:          "carml client",
		Short:        "Runs the carml client",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			frameworkName := strings.ToLower(framework.FrameworkName)
			frameworkVersion := strings.ToLower(framework.FrameworkVersion)
			modelName := strings.ToLower(model.ModelName)
			modelVersion := strings.ToLower(model.ModelVersion)

			agents, err := http.Models.Agents(frameworkName, frameworkVersion, modelName, modelVersion)
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
			})
			if err != nil {
				return errors.Wrap(err, "unable to open the predictor")
			}

			defer client.Close(ctx, predictor)

			var urls []*dlframework.URLsRequest_URL
			for i, url := range data {
				urls = append(urls, &dlframework.URLsRequest_URL{
					Id:   string(i),
					Data: url,
				})
			}
			urlReq := dlframework.URLsRequest{
				Predictor: predictor,
				Urls:      urls,
			}
			res, err := client.URLs(ctx, &urlReq)
			if err != nil {
				return errors.Wrap(err, "unable to get response from urls request")
			}

			pp.Println(res)

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
	c.PersistentFlags().BoolVarP(&cmd.Local, "local", "l", false, "Listen on local address.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("app.secret", c.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", c.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", c.PersistentFlags().Lookup("verbose"))
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cmd.Init()
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/client")
	})
}
