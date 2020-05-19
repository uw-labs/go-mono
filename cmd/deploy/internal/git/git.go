package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

// Metadata contains git metadata
type Metadata struct {
	GitSHA    string
	GitBranch string
	BuildTime time.Time
}

// GetMetadata reads the git metadata from the repo root.
func GetMetadata(repoRoot string) (*Metadata, error) {
	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("open local repository: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("get HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("get base commit: %w", err)
	}

	branch := headRef.Name().Short()
	// Remove commas
	branch = strings.ReplaceAll(branch, ":", "-")
	// Remove slashes
	branch = strings.ReplaceAll(branch, "/", "-")

	md := &Metadata{
		GitSHA:    headCommit.Hash.String(),
		GitBranch: branch,
		BuildTime: time.Now(),
	}

	return md, nil
}
