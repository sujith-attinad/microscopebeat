package main

import (
	"os"

	"github.com/sujith-attinad/microscopebeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
