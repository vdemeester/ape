package main // import "go.sbr.pm/ape/cmd/ape"

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.sbr.pm/ape/walk"
)

var verbose bool

func walkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "walk [OPTIONS] COMMAND",
		Short: "walk the directory and its subdirectory and run a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return nil
			}
			return errors.Wrapf(walk.Walk(context.Background(), cwd, args, verbose), "Failed to walk")
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	return cmd
}
