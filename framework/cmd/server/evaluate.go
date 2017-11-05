// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"context"

	"github.com/rai-project/dlframework/framework/cmd"
	shellwords "github.com/junegunn/go-shellwords"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/sirupsen/logrus"
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
		256,
		64,
		8,
		1,
	}
	timeout                  = 4*time.Hour
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
		mainFile := filepath.Join(sourcePath, framework+".go")
		cmd := exec.Command("go", "build", mainFile)
		err := cmd.Run()
		if err != nil {
			log.WithError(err).
				WithField("framework", framework).
				Error("failed to compile " + mainFile)
			continue
		}

		for _, model := range models {
			progress := pb.New(len(batchSizes)).Prefix(model)
			progress.Start()
			for _, batchSize := range batchSizes {
				ctx, _ := context.WithTimeout(context.Background(), timeout)
				shellCmd := "dataset" +
					" --publish=true" +
					fmt.Sprintf(" --gpu=%v", usingGPU) + fmt.Sprintf(" -b %v", batchSize) +
					fmt.Sprintf(" --modelName=%v", model)
				args, err := shellwords.Parse(shellCmd)
				if err != nil {
					log.WithError(err).WithField("cmd", shellCmd).Error("failed to parse shell command")
					os.Exit(-1)
				}

				cmd := exec.Command(filepath.Join(sourcePath, framework), args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Start()
				if err != nil {
					progress.Increment()
					log.WithError(err).WithField("cmd", shellCmd).Error("failed to run command")
					continue
				}

				done := make(chan error)
				go func() { done <- cmd.Wait() }()

				select {
				case err := <-done:
					if err != nil {
						log.WithError(err).WithField("cmd", shellCmd).Error("failed to wait for command")
					}
				case <-ctx.Done():
					cmd.Process.Kill()
					log.WithError(ctx.Err()).WithField("cmd", shellCmd).Error("command timeout")
				}

				progress.Increment()
			}
			progress.FinishPrint("finished evaluation of " + framework + "/" + model)
		}
	}
}

func init() {

}
