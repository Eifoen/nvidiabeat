package main

import (
	"os"

	"github.com/Eifoen/nvidiabeat/cmd"

	_ "github.com/Eifoen/nvidiabeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
