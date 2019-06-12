package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/k0kubun/pp"
	"github.com/rai-project/tracer/convert/flame"
)

var flameGraphCmd = &cobra.Command{
	Use:   "flamegraph",
	Short: "Create a framegraph out of a trace",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running " + cmd.Name())
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tr, err := flame.ConvertTraceFile(tracePath)
		if err != nil {
			pp.Println(err)
			os.Exit(1)
		}
		outputJSONFile := filepath.Join(outputPath, "flamegraph.json")
		err = ioutil.WriteFile(outputJSONFile, tr, 0600)
		if err != nil {
			pp.Println(err)
			os.Exit(1)
		}

		outputHTMLFile := filepath.Join(outputPath, "flamegraph.html")
		wr := &bytes.Buffer{}
		flame.GenerateHtml(wr, "flame", string(tr))
		err = ioutil.WriteFile(outputHTMLFile, wr.Bytes(), 0600)
		if err != nil {
			pp.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
