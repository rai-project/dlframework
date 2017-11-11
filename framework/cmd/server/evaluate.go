// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"context"

	shellwords "github.com/junegunn/go-shellwords"
	"github.com/k0kubun/pp"
	"github.com/rai-project/dlframework/framework/cmd"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/rai-project/config"
	"github.com/rai-project/cpu/cpuid"
	_ "github.com/rai-project/logger/hooks"
	"github.com/sirupsen/logrus"
)

func parse(model string) (modelName, version string) {
	splt := strings.Split(model, "_")
	modelName, modelVersion := splt[0], splt[1]
	return modelName, modelVersion
}

var (
	models = []string{
		// "SqueezeNet_1.0",
		// "SqueezeNet_1.1",
		// "BVLC-AlexNet_1.0",
		//"BVLC-GoogleNet_1.0",
		// "VGG16_1.0",
		// "ResNet101_1.0",
		// "ResNet101_2.0",
		"WRN50_2.0",
		// "BVLC-Reference-RCNN-ILSVRC13_1.0",
		// "Inception_3.0",
		"Inception_4.0",
		// "ResNeXt50-32x4d_1.0",
	}

	frameworks = []string{
		// "mxnet",
		// "caffe2",
		"caffe",
	}
	batchSizes = []int{
		// 256,
		// 128,
		// 64,
		// 32,
		//16,
		//8,
		4,
		// 2,
		// 1,
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

	for _, usingGPU := range []bool{true} {
		var device string
		if usingGPU {
			device = "gpu"
		} else {
			device = "cpu"
		}
		for _, framework := range frameworks {
			mainFile := filepath.Join(sourcePath, framework+".go")
			compileArgs := []string{
				"build",
			}
			if !cpuid.SupportsAVX() {
				compileArgs = append(compileArgs, "-tags=noasm")
			}
			compileArgs = append(compileArgs, mainFile)
			fmt.Printf("Compiling using :: go %#v\n", compileArgs)
			cmd := exec.Command("go", compileArgs...)
			err := cmd.Run()
			if err != nil {
				log.WithError(err).
					WithField("framework", framework).
					Error("failed to compile " + mainFile)
				continue
			}

			for _, model := range models {
				modelName, modelVersion := parse(model)
				for _, batchSize := range batchSizes {
					pp.Println("Running", framework, "::", model, "on", device, "with batch size", batchSize)
					ctx, _ := context.WithTimeout(context.Background(), timeout)
					shellCmd := "dataset" +
						" --debug" +
						" --verbose" +
						" --publish=true" +
						fmt.Sprintf(" --gpu=%v", usingGPU) +
						fmt.Sprintf(" --batch_size=%v", batchSize) +
						fmt.Sprintf(" --model_name=%v", modelName) +
						// " --publish_predictions=true" +
						fmt.Sprintf(" --model_version=%v", modelVersion)
					shellCmd = shellCmd + " " + strings.Join(os.Args, " ")
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
				pp.Println("Finished running", framework, "::", model, "on", device)
			}
		}

	}
}

func init() {

}
