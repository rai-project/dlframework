package main

import (
	"fmt"
	"os"

	"github.com/rai-project/dlframework/framework/cmd/client"
)

var (
	framework = client.Framework{
		FrameworkName:    "MxNet",
		FrameworkVersion: "0.11.0",
	}
	model = client.Model{
		ModelName:    "SqueezeNet",
		ModelVersion: "1.0",
	}
)

func main() {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}
	rootCmd, err := client.NewRootCommand(framework, model, imgURLs)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
