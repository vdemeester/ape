package main

import (
	"os"

	"github.com/spf13/cobra"
)

func apeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ape",
		Short: "vcs mirror update",
	}
	cmd.AddCommand(upCmd())
	return cmd
}

func main() {
	cmd := apeCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
