package server

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	shutdown "github.com/klauspost/shutdown2"
	raicmd "github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/cmd"
	evalcmd "github.com/rai-project/evaluation/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//dllayer "github.com/rai-project/dllayer/cmd"

	_ "github.com/rai-project/logger/hooks"
)

var (
	local   bool
	profile bool
)

type FrameworkRegisterFunction func()

// represents the base command when called without any subcommands
func NewRootCommand(frameworkRegisterFunc FrameworkRegisterFunction, framework0 dlframework.FrameworkManifest) (*cobra.Command, error) {
	frameworkName := framework0.GetName()
	rootCmd := &cobra.Command{
		Use:   strings.ToLower(frameworkName) + "-agent",
		Short: "Run the MLModelScope " + frameworkName + " agent",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			frameworkRegisterFunc()
			framework = framework0
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

	c.AddCommand(serveCmd)
	c.AddCommand(predictCmd)
	c.AddCommand(downloadCmd)
	addContainerCmd(c)
	c.AddCommand(infoCmd)
	c.AddCommand(evalcmd.EvaluationCmd)
	c.AddCommand(traceCmd)

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
	})
}
