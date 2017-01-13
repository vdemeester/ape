package main

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ape",
	Short: "vcs mirror update",
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
