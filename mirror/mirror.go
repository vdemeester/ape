package mirror

import (
	"crypto/sha256"
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
	if !exists {
		if err := cloneRepository(workingDirectory, repository.Remote, repositoryName); err != nil {
			return err
		}
	}
	if err := addRemoteIfNotPresent(workingDirectory, repositoryName, repository.Upstream); err != nil {
		return err
	}
	if err := fetchAndRebase(workingDirectory, repositoryName, repository.Remote); err != nil {
		return err
	}
	return pushRepository(workingDirectory, repositoryName, repository.Remote)
}

func pushRepository(workingDirectory, name, remote string) error {
	fmt.Fprintf(os.Stderr, "🐵 %s: push to origin\n", remote)
	dir := filepath.Join(workingDirectory, name)
	cmd := exec.Command("git", "push", "-f", "origin", "master")
	cmd.Dir = dir
	return errors.Wrapf(cmd.Run(), "error pushing to origin in %s", dir)
}

func fetchAndRebase(workingDirectory, name, remote string) error {
	fmt.Fprintf(os.Stderr, "🙊 %s: fetch and rebase\n", remote)
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

func addRemoteIfNotPresent(workingDirectory, name, upstreamRemote string) error {
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

func cloneRepository(workingDirectory, remote, name string) error {
	fmt.Fprintf(os.Stderr, "🐵 %s: cloning in %s\n", remote, name)
	cmd := exec.Command("git", "clone", remote, name)
	cmd.Dir = workingDirectory
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error cloning %s in %s", remote, workingDirectory)
	}
	return nil
}

func lookInWorkingDirectory(workingDirectory string, repository string) (string, bool) {
	sum := sha256.Sum256([]byte(repository))
	encodedName := fmt.Sprintf("%x", sum)
	_, err := os.Stat(filepath.Join(workingDirectory, encodedName))
	return encodedName, err == nil
}
