// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"context"

	shellwords "github.com/junegunn/go-shellwords"
	"github.com/rai-project/dlframework/framework/cmd"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/rai-project/config"
	"github.com/rai-project/cpu/cpuid"
	_ "github.com/rai-project/logger/hooks"
	"github.com/sirupsen/logrus"
)

var (
	models = []string{
		"BVLC-AlexNet",
		"BVLC-GoogleNet",
		"VGG16",
		"ResNet101",
	}

	frameworks = []string{
		"mxnet",
		"caffe2",
		"caffe",
	}
	batchSizes = []int{
		//256,
		//64,
		//50,
		//32,
		16,
		8,
		//1,
	}
	timeout                  = 4 * time.Hour
	usingGPU                 = true
	sourcePath               = sourcepath.MustAbsoluteDir()
	log        *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
)

func main() {
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/evaluate")
	})

	cmd.Init()

	tags := ""
	if !cpuid.SupportsAVX() {
		tags = "-tags=noasm"
	}

	for _, usingGPU := range []bool{true, false} {
		var device string
		if usingGPU {
			device = "gpu"
		} else {
			device = "cpu"
		}
		for _, framework := range frameworks {
			mainFile := filepath.Join(sourcePath, framework+".go")
			fmt.Println("Compiling using :: ", "go", "build", tags, mainFile)
			cmd := exec.Command("go", "build", tags, mainFile)
			err := cmd.Run()
			if err != nil {
				log.WithError(err).
					WithField("framework", framework).
					Error("failed to compile " + mainFile)
				continue
			}

			for _, model := range models {
				for _, batchSize := range batchSizes {
					fmt.Println("Running", framework, "::", model, "on", device, "with batch size", batchSize)
					ctx, _ := context.WithTimeout(context.Background(), timeout)
					shellCmd := "dataset" +
						" --publish=true" +
						fmt.Sprintf(" --gpu=%v", usingGPU) + fmt.Sprintf(" -b %v", batchSize) +
						fmt.Sprintf(" --modelName=%v", model)
					args, err := shellwords.Parse(shellCmd)
					if err != nil {
						log.WithError(err).WithField("cmd", shellCmd).Error("failed to parse shell command")
						//os.Exit(-1)
						continue
					}

					cmd := exec.Command(filepath.Join(sourcePath, framework), args...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					err = cmd.Start()
					if err != nil {
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
				}
				fmt.Println("Finished running", framework, "::", model, "on", device)
			}
		}

	}
}

func init() {

}
