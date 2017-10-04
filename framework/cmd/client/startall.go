package client

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	base = "src/github.com/rai-project"
	all  = []string{"mxnet", "tensorflow"}
)

func copyLogs(r io.Reader, logfn func(args ...interface{})) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			logfn(buf[0:n])
		}
		if err != nil {
			break
		}
	}
}

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
