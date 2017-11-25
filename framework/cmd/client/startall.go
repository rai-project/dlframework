package client

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	base = "src/github.com/rai-project"
	all  = []string{"mxnet", "tensorflow", "caffe", "caffe2", "tensorrt"}
)

var startallCmd = &cobra.Command{
	Use:     "startallCmd",
	Short:   "startallCmd",
	Aliases: []string{"startall", "start"},
	Long:    `startallCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		for _, framework := range all {
			go func(framework string) {
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
		select {}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(startallCmd)
}
