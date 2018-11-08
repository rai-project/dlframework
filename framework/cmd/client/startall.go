package client

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	shutdown "github.com/klauspost/shutdown2"
	"github.com/spf13/cobra"
)

var (
	base   = "src/github.com/rai-project"
	agents = []string{}
)

var startallCmd = &cobra.Command{
	Use:     "startall",
	Short:   "starts the " + strings.Join(agents, " ") + " agents",
	Aliases: []string{"startall", "start"},
	Long:    `starts all the MLModelScope agents`,
	RunE: func(c *cobra.Command, args []string) error {
		var wg sync.WaitGroup
		quitAll := make(chan struct{})
		shutdown.SecondFn(func() {
			defer close(quitAll)
		})
		for _, framework := range agents {
			wg.Add(1)
			go func(framework string) {
				defer wg.Done()
				main := filepath.Join(os.Getenv("GOPATH"), base, framework, framework+"-agent", "main.go")
				args := []string{
					"run",
					main,
					"-l",
					"-d",
					"-v",
				}
				frameworkExited := make(chan error)
				var cmd *exec.Cmd
				runFramework := func() {
					cmd = exec.Command("go", args...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					frameworkExited <- cmd.Run()
				}
				go runFramework()
				for {
					select {
					case <-frameworkExited:
						go runFramework()
					case <-quitAll:
						if cmd != nil {
							cmd.Process.Signal(syscall.SIGTERM)
							cmd.Process.Kill()
						}
						return
					}
				}
			}(framework)
		}
		wg.Wait()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(startallCmd)
	agents = []string{
		"mxnet",
		//	"tensorflow",
		"caffe",
		"caffe2",
	}
	if runtime.GOOS != "linux" {
		return
	}
	if runtime.GOARCH == "ppc64le" {
		return
	}
	agents = append(agents, "tensorrt")
	if runtime.GOARCH == "arm64" {
		return
	}
	agents = append(agents, "cntk")
}
