package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vdemeester/ape/config"
	"github.com/vdemeester/ape/mirror"
)

var upCmd = &cobra.Command{
	Use:   "up path",
	Short: "update mirrors",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("up requires an argument, the path where repositories are")
		}
		cfg, err := os.Open(filepath.Join(os.Getenv("HOME"), ".config/ape.conf"))
		if err != nil {
			return errors.Wrap(err, "Failed to open config file")
		}
		repositories, err := config.Parse(cfg)
		if err != nil {
			return err
		}
		if err := mirror.Mirror(args[0], repositories); err != nil {
			return errors.Wrap(err, "Some repositories failed to get mirrored")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(upCmd)
}
