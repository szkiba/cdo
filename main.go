package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/szkiba/cdo/internal/cmd"
	"mvdan.cc/sh/v3/interp"
)

func main() {
	root, err := cmd.New(os.Args[1:])
	cobra.CheckErr(err)

	err = root.Execute()

	if status, ok := interp.IsExitStatus(err); ok {
		os.Exit(int(status))
	}

	cobra.CheckErr(err)
}
