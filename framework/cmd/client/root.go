package client

import (
	"strings"

	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework/framework/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	log              *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/client")
	frameworkName    string
	frameworkVersion string
	modelName        string
	modelVersion     string
	batchSize        int
)

var RootCmd = &cobra.Command{
	Use:          "carml client",
	Short:        "Runs the carml client",
	SilenceUsage: true,
}

func init() {
	cobra.OnInitialize(initConfig)
	setup(RootCmd)
	RootCmd.PersistentFlags().StringVar(&frameworkName, "frameworkName", "MxNet", "frameworkName")
	RootCmd.PersistentFlags().StringVar(&frameworkVersion, "frameworkVersion", "0.11.0", "frameworkVersion")
	RootCmd.PersistentFlags().StringVar(&modelName, "modelName", "CaffeNet", "modelName")
	RootCmd.PersistentFlags().StringVar(&modelVersion, "modelVersion", "1.0", "modelVersion")
	RootCmd.PersistentFlags().IntVar(&batchSize, "batchSize", 16, "batch size")
	cleanNames()
}

func setup(c *cobra.Command) {
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("app.secret", c.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", c.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", c.PersistentFlags().Lookup("verbose"))
}

func cleanNames() {
	frameworkName = strings.ToLower(frameworkName)
	frameworkVersion = strings.ToLower(frameworkVersion)
	modelName = strings.ToLower(modelName)
	modelVersion = strings.ToLower(modelVersion)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cmd.Init()
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/client")
	})
}
