// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/fatih/color"
	shellwords "github.com/junegunn/go-shellwords"
	"github.com/rai-project/config"
	"github.com/rai-project/cpu/cpuid"
	dl "github.com/rai-project/dlframework"
	dlcmd "github.com/rai-project/dlframework/framework/cmd"
	"github.com/sirupsen/logrus"

	_ "github.com/rai-project/logger/hooks"
)

var (
	frameworks = []string{
		// "mxnet",
		// "caffe2",
		"tensorflow",
		// "caffe",
		//"tensorrt",
		// "cntk",
	}

	models = []string{
		// "SqueezeNet_1.0",
		// "SqueezeNet_1.1",
		// "BVLC-AlexNet_1.0",
		// "BVLC-Reference-CaffeNet_1.0",
		// "BVLC-GoogLeNet_1.0",
		// "ResNet101_1.0",
		// "ResNet101_2.0",
		// "WRN50_2.0",
		// "BVLC-Reference-RCNN-ILSVRC13_1.0",
		// "Inception_3.0",
		// "Inception_4.0",
		// "ResNeXt50-32x4d_1.0",
		// "VGG16_1.0",
		// "VGG19_1.0",
		// "ResNet50_1.0",
		// "SphereFace_1.0",
		// "ShuffleNet_Caffe2_1.0",
	}

	batchSizes = []int{
		// 512,
		// // 448,
		// // 384,
		// // 320,
		// 256,
		// // 196,
		// 128,
		// // 96,
		// 64,
		// // 48,
		// 32,
		// 16,
		// 8,
		// 4,
		// 2,
		1,
	}

	useGPU = []bool{
		true,
		false,
	}

	timeout                       = 300 * time.Minute
	databaseAddress               = "localhost"
	traceLevel                    = "NO_TRACE"
	sourcePath                    = sourcepath.MustAbsoluteDir()
	log             *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd/run")
	debug                         = false
)

func main() {
	config.AfterInit(func() {
		log = logrus.New().WithField("pkg", "dlframework/framework/cmd/run")
	})

	dlcmd.Init()
	for i := 0; i < 1; i++ {
		for _, usingGPU := range useGPU {
			var device string
			if usingGPU {
				device = "gpu"
			} else {
				device = "cpu"
			}
			for _, framework := range frameworks {
				compileArgs := []string{
					"build",
				}
				if runtime.GOARCH == "amd64" && !cpuid.SupportsAVX() {
					compileArgs = append(compileArgs, "-tags=noasm")
				}
				if !usingGPU {
					compileArgs = append(compileArgs, "-tags=nogpu")
				}
				if debug {
					compileArgs = append(compileArgs, "-tags=debug")
				}
				cmd := exec.Command("go", compileArgs...)
				agentPath := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/rai-project/", framework, framework+"-agent")
				cmd.Dir = agentPath
				fmt.Printf("Compiling using :: go %#v in %v\n", compileArgs, cmd.Dir)
				err := cmd.Run()
				if err != nil {
					log.WithError(err).
						WithField("framework", framework).
						Error("failed to compile agent")
					continue
				}

				for _, model := range models {
					modelName, modelVersion := dlcmd.ParseModelName(model)
					for _, batchSize := range batchSizes {
						color.Red("⇛ Running %v :: %v on %v with batch size %v", framework, model, device, batchSize)
						ctx, _ := context.WithTimeout(context.Background(), timeout)
						shellCmd := "predict dataset" +
							" --debug" +
							" --verbose" +
							" --publish=false" +
							" --publish_predictions=false" +
							" --fail_on_error=true" +
							" --warmup_num_file_parts=0" +
							" --num_file_parts=-1" +
							fmt.Sprintf(" --use_gpu%v", usingGPU) +
							fmt.Sprintf(" --batch_size=%v", batchSize) +
							fmt.Sprintf(" --model_name=%v", modelName) +
							fmt.Sprintf(" --model_version=%v", modelVersion) +
							fmt.Sprintf(" --database_name=%v", dl.CleanString(modelName+"_"+modelVersion+"_predictions")) +
							fmt.Sprintf(" --database_address=%v", databaseAddress) +
							fmt.Sprintf(" --trace_level=%v", traceLevel)
						if len(os.Args) < 3 {
							panic("Need to set database_adress, tracer_address and trace_level")
						}
						shellCmd = shellCmd + " " + strings.Join(os.Args, " ")
						args, err := shellwords.Parse(shellCmd)
						if err != nil {
							log.WithError(err).WithField("cmd", shellCmd).Error("failed to parse shell command")
							continue
						}
						fmt.Println("Running " + shellCmd)
						var agentCmd = agentPath + framework + "-agent"
						cmd := exec.Command(agentCmd, args...)
						cmd.Dir = agentPath
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
					color.Red("⇛ Finished running %v :: %v on %v", framework, model, device)
				}
			}
		}
	}
}

func init() {

}
