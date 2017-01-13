package config

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	ConfigFile = "ape.conf"
)

// Repository represents a vcs repository that is supposed to get mirrored
type Repository struct {
	Remote   string
	Upstream string
}

func Parse(r io.Reader) ([]Repository, error) {
	repositories := []Repository{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		ln := strings.TrimSpace(s.Text())
		if strings.HasPrefix(ln, "#") || ln == "" {
			continue
		}
		cidx := strings.Index(ln, "#")
		if cidx > 0 {
			ln = ln[:cidx]
		}
		ln = strings.TrimSpace(ln)
		parts := strings.Fields(ln)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config format: %s", ln)
		}
		repository := Repository{
			Remote:   parts[0],
			Upstream: parts[1],
		}
		repositories = append(repositories, repository)
	}
	return repositories, nil
}
