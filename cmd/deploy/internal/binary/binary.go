package binary

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Request is the input to the Build command
type Request struct {
	RepoRoot string
	Name     string
	MainPath string
}

// Build builds a CGO-disabled Go binary using a local version of
// "go" and returns the path where the binary lives.
func Build(ctx context.Context, logger *logrus.Logger, req *Request) (string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("find go binary: %w", err)
	}

	tempDir, err := ioutil.TempDir("", "build")
	if err != nil {
		return "", fmt.Errorf("create temp directory: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		goBin,
		"build",
		"-mod=vendor",
		"-o",
		filepath.Join(tempDir, "app"),
		filepath.Join(req.RepoRoot, req.MainPath),
	)

	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Infoln("Build output:\n", string(out))
		return "", err
	}
	return filepath.Join(tempDir, "app"), nil
}
