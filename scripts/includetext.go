package main

import (
	"io/ioutil"
	"os"
	"strings"
)

// Reads all .json files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fs, _ := ioutil.ReadDir(".")
	out, _ := os.Create("swagger.pb.go")
	out.Write([]byte("package dlframework \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			name := strings.Replace(strings.TrimPrefix(f.Name(), "service."), ".", "_", -1)
			out.Write([]byte(strings.TrimSuffix(name, "_json") + " = `"))
			f, _ := os.Open(f.Name())
			defer f.Close()

			bts, _ := ioutil.ReadAll(f)
			str := strings.Replace(string(bts), "`", "'", -1)
			out.Write([]byte(str + "`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
