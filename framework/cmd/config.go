package cmd

import (
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd")
)

// Init reads in config file and ENV variables if set.
func Init(cfgFile, appSecret string) {

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
