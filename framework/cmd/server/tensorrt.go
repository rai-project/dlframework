// +build ignore

package main

import (
	"fmt"
	"os"

	cmd "github.com/rai-project/dlframework/framework/cmd/server"
	"github.com/rai-project/tensorrt"
	_ "github.com/rai-project/tensorrt/predict"
	"github.com/rai-project/tracer"

	_ "github.com/rai-project/tracer/all"
)

func main() {

	rootCmd, err := cmd.NewRootCommand(tensorrt.FrameworkManifest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer tracer.Close()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
