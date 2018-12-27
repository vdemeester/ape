package mirror

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/vdemeester/ape/config"
)

func Mirror(workingDirectory string, repositories []config.Repository) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(repositories))
	limit := make(chan struct{}, 8)
	for _, d := range repositories {
		wg.Add(1)
		go func(d config.Repository) {
			limit <- struct{}{}
			errCh <- mirror(workingDirectory, d)
			wg.Done()
			<-limit
		}(d)
	}
	wg.Wait()
	close(errCh)
	var errs []string
	for err := range errCh {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("Errors on clone:\n%s", strings.Join(errs, "\n"))
}

func mirror(workingDirectory string, repository config.Repository) error {
	repositoryName, exists := lookInWorkingDirectory(workingDirectory, repository.Remote)
	// fmt.Printf("🐒 %s\n", repositoryName)
	if !exists {
		if err := cloneRepository(workingDirectory, repository.Remote, repositoryName); err != nil {
			return err
		}
	}
	if err := addRemoteIfNotPresent(workingDirectory, repositoryName, repository.Upstream); err != nil {
		return err
	}
	if err := fetchAndRebase(workingDirectory, repositoryName); err != nil {
		return err
	}
	return pushRepository(workingDirectory, repositoryName)
}

func pushRepository(workingDirectory string, name string) error {
	fmt.Fprintf(os.Stderr, "🐵 %s:  push to origin\n", name)
	dir := filepath.Join(workingDirectory, name)
	cmd := exec.Command("git", "push", "-f", "origin", "master")
	cmd.Dir = dir
	return errors.Wrapf(cmd.Run(), "error pushing to origin in %s", dir)
}

func fetchAndRebase(workingDirectory string, name string) error {
	fmt.Fprintf(os.Stderr, "🙊 %s: fetch and rebase\n", name)
	dir := filepath.Join(workingDirectory, name)
	cmd := exec.Command("git", "fetch", "-p", "--all")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error fetching all remotes in %s", dir)
	}
	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error checking out master in %s", dir)
	}
	cmd = exec.Command("git", "rebase", "--autostash", "upstream/master")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error rebasing upstream/master to master in %s", dir)
	}
	return nil
}

func addRemoteIfNotPresent(workingDirectory string, name string, upstreamRemote string) error {
	dir := filepath.Join(workingDirectory, name)
	cmd := exec.Command("git", "remote")
	cmd.Dir = dir
	content, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "error looking at remotes in %s", dir)
	}
	if strings.Contains(string(content), "upstream") {
		return nil
	}

	fmt.Fprintf(os.Stderr, "🙉 %s: add upstream\n", upstreamRemote)
	cmd = exec.Command("git", "remote", "add", "upstream", upstreamRemote)
	cmd.Dir = dir
	return cmd.Run()
}

func cloneRepository(workingDirectory string, remote string, name string) error {
	fmt.Fprintf(os.Stderr, "🐵 %s: cloning\n", remote)
	cmd := exec.Command("git", "clone", remote, name)
	cmd.Dir = workingDirectory
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error cloning %s in %s", remote, workingDirectory)
	}
	return nil
}

func lookInWorkingDirectory(workingDirectory string, repository string) (string, bool) {
	extension := filepath.Ext(repository)
	base := filepath.Base(repository)
	name := base[0 : len(base)-len(extension)]
	_, err := os.Stat(filepath.Join(workingDirectory, name))
	return name, err == nil
}
