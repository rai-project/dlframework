package client

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/k0kubun/pp"

	"github.com/spf13/cobra"
)

var (
	base = "src/github.com/rai-project"
	all  = []string{"mxnet", "tensorflow"}
)

//os.Getenv("GOPATH")
var startallCmd = &cobra.Command{
	Use:     "startallCmd",
	Short:   "startallCmd",
	Aliases: []string{"startall", "start"},
	Long:    `startallCmd`,
	RunE: func(c *cobra.Command, args []string) error {
		for _, framework := range all {
			go func() {
				main := filepath.Join(os.Getenv("GOPATH"), base, framework, framework+"-agent", "main.go")
				args := []string{
					"run",
					main,
					"-l",
					"-d",
					"-v",
				}
				cmd := exec.Command("go", args...)
				buf, err := cmd.CombinedOutput()
				if err != nil {
					log.WithError(err).Error("Failed to run go " + strings.Join(args, " "))
				}
				log.Infof(string(buf))
				pp.Println(string(buf))
			}()
		}

		pp.Println("launched all agents")
		select {}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(startallCmd)
}
