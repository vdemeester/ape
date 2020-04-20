package main // import "go.sbr.pm/ape/cmd/ape"

import (
	"os"

	"github.com/spf13/cobra"
)

func apeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ape",
		Short: "vcs mirror update",
	}
	cmd.AddCommand(cloneCmd())
	cmd.AddCommand(upCmd())
	cmd.AddCommand(walkCmd())
	return cmd
}

func main() {
	cmd := apeCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
