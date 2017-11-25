// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/rai-project/caffe2"
	_ "github.com/rai-project/caffe2/predict"
	cmd "github.com/rai-project/dlframework/framework/cmd/server"
	"github.com/rai-project/tracer"

	_ "github.com/rai-project/tracer/all"
)

func main() {

	rootCmd, err := cmd.NewRootCommand(caffe2.Register, caffe2.FrameworkManifest)
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
