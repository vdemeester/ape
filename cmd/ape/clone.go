package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func cloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "clone",
		Short:              "clone a repository",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			posArgs := []string{}
			for _, a := range args {
				if strings.HasPrefix(a, "-") {
					continue
				}
				posArgs = append(posArgs, a)
			}
			switch len(posArgs) {
			case 1:
				// only the clone, do some magic
				// 1. detect host in git url (github.com, gitlab) for http, ssh and git
				// 2. call git with the correct args
				return smartGitClone(posArgs[0], args)
			default:
				// shell everything out !
				return gitClone(args)
			}
		},
	}
	return cmd
}

func smartGitClone(url string, args []string) error {
	exp := regexp.MustCompile(`(?:git@|git://|https://)(.+)[:/]([^/]+)/(.+).git`)
	if exp.MatchString(url) {
		parts := exp.FindStringSubmatch(url)
		switch parts[1] {
		case "github.com", "gitlab.com":
			home, err := homedir.Dir()
			if err != nil {
				return err
			}
			path := filepath.Join(append([]string{home, "src"}, parts[1:]...)...)
			args = append(args, path)
		}
	}
	return gitClone(args)
}

func gitClone(args []string) error {
	args = append([]string{"clone"}, args...)
	gcmd := exec.Command("git", args...)
	gcmd.Stdin = os.Stdin
	gcmd.Stderr = os.Stderr
	gcmd.Stdout = os.Stdout
	gcmd.Env = os.Environ()
	return gcmd.Run()
}
