package client

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
	Long:    `starts all the CarML agents`,
	RunE: func(c *cobra.Command, args []string) error {
		var wg sync.WaitGroup
		for _, framework := range all {
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
				cmd := exec.Command("go", args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}(framework)
		}
		wg.Wait()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(startallCmd)
	agents = []string{"mxnet", "tensorflow", "caffe", "caffe2"}
	if runtime.GOOS != "linux" {
		return
	}
	if runtime.GOARCH == "ppc64le" {
		return
	}
	args = append(agents, "tensorrt")
	if runtime.GOARCH == "arm64" {
		return
	}
	args = append(agents, "cntk")
}
