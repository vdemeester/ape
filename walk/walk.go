package walk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Walk(ctx context.Context, cwd string, args []string, verbose bool) error {
	g, _ := errgroup.WithContext(context.Background())
	limit := make(chan struct{}, 4)
	if err := filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			dirname := filepath.Join(path, ".git")
			if _, err := os.Stat(dirname); err == nil {
				g.Go(func() error {
					limit <- struct{}{}
					relPath, _ := filepath.Rel(cwd, path)
					err := execute(ctx, relPath, args, verbose)
					<-limit
					return err
				})
				return filepath.SkipDir
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return g.Wait()
}

func execute(ctx context.Context, path string, command []string, verbose bool) error {
	fmt.Printf("🐵 %s: %s …\n", path, strings.Join(command, " "))
	var buf io.Writer = bytes.NewBuffer([]byte{})
	if verbose {
		buf = os.Stderr
	}
	c := exec.CommandContext(ctx, command[0], command[1:]...)
	c.Dir = path
	c.Stderr = buf
	c.Stdout = buf
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		if !verbose {
			fmt.Fprint(os.Stderr, buf)
		}
		fmt.Fprintf(os.Stderr, "🙊 %s : %s … %s\n", path, strings.Join(command, " "), err)
		return err
	}
	return nil
}
