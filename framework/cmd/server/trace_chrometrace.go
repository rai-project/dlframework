package server

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer/convert/chrome"
	"github.com/spf13/cobra"
)

var chromeTraceCmd = &cobra.Command{
	Use:   "chrome",
	Short: "Convert a trace to chrome trace format",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running " + cmd.Name())
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tr, err := chrome.ConvertFile(tracePath)
		if err != nil {
			pp.Println(err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(outputPath, tr, 0600)
		if err != nil {
			pp.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
