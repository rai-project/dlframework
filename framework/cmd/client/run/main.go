package main

import (
	"fmt"
	"os"

	"github.com/rai-project/dlframework/framework/cmd/client"
)

func main() {
	if err := client.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
