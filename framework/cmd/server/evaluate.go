// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"bitbucket.org/c3sr/p3sr-pdf/cmd"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/rai-project/config"
)

var (
	models = []string{
		"BVLC-AlexNet",
		"BVLC-GoogleNet",
		"VGG16",
		"RestNet50",
	}

	frameworks = []string{
		"mxnet",
	}
	batchSizes = []int{
		1,
		8,
		64,
		256,
	}
	usingGPU                 = true
	sourcePath               = sourcepath.MustAbsoluteDir()
	log        *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
)

func main() {

	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
	})

	cmd.Init()

	for _, framework := range frameworks {
		for _, model := range models {
			progress := pb.New(len(batchSizes)).Prefix(model)
			progress.Start()
			for _, batchSize := range batchSizes {

				main := filepath.Join(sourcePath, "main.go")
				shellCmd := "run " +
					main +
					" dataset" +
					" --publish=true" +
					fmt.Sprintf(" --gpu=%v", usingGPU) + fmt.Sprintf(" -b %v", batchSize) +
					fmt.Sprintf(" --modelName=%v", modelName)
				args, err := shellwords.Parse(shellCmd)
				if err != nil {
					log.WithError(err).WithField("cmd", shellCmd).Error("failed to parse shell command")
					os.Exit(-1)
				}

				cmd := exec.Command("go", args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()

				progress.Increment()
			}
			progress.FinishPrint("finished evaluation of " + framework + "/" + model)
		}
	}
}

func init() {

}
