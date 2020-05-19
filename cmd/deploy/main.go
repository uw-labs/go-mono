package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/uw-labs/go-mono/cmd/deploy/internal/binary"
	"github.com/uw-labs/go-mono/cmd/deploy/internal/deploy"
	"github.com/uw-labs/go-mono/cmd/deploy/internal/docker"
	"github.com/uw-labs/go-mono/cmd/deploy/internal/git"
	pkgcontext "github.com/uw-labs/go-mono/pkg/context"
)

var (
	repoRoot       = flag.String("repo-root", ".", "The root of the repo, to find the git folder.")
	dockerUser     = flag.String("docker-user", "", "The docker user to use when authenticating against the registry.")
	dockerPassword = flag.String("docker-password", "", "The password to use when authenticating the user against the registry.")
	dockerRegistry = flag.String("docker-registry", "docker.pkg.github.com/uw-labs/go-mono", "The registry to push images to. Can include any subpaths.")
	deployFile     = flag.String("deploy-file", "", "The deploy file to read deployment configuration from.")
)

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	}

	if *deployFile == "" {
		logger.Fatal("deploy-file must be specified")
	}

	err := run(logger, *repoRoot, *dockerUser, *dockerPassword, *dockerRegistry, *deployFile)
	if err != nil {
		logger.WithError(err).Fatal()
	}
}

func run(logger *logrus.Logger, repoRoot, dockerUser, dockerPassword, dockerRegistry, deployFile string) error {
	ctx := pkgcontext.WithSignalHandler(context.Background())

	md, err := git.GetMetadata(repoRoot)
	if err != nil {
		return fmt.Errorf("get git metadata: %w", err)
	}

	conf, err := deploy.Parse(repoRoot, deployFile)
	if err != nil {
		return fmt.Errorf("parse deployment: %w", err)
	}

	logger.Infoln("Deploying", conf.Name)

	logger.Infoln("Building binary")
	binPath, err := binary.Build(ctx, logger, &binary.Request{
		Name:     conf.Name,
		RepoRoot: repoRoot,
		MainPath: conf.Main,
	})
	if err != nil {
		return fmt.Errorf("build binary: %w", err)
	}
	defer func() {
		err := os.RemoveAll(filepath.Dir(binPath))
		if err != nil {
			logger.WithError(err).Infof("remove binary directory (%s)", filepath.Dir(binPath))
		}
	}()

	logger.Infoln("Building Docker image")
	digest, err := docker.BuildAndPushImage(ctx, logger, &docker.Request{
		RepoRoot:         repoRoot,
		BinaryPath:       binPath,
		Registry:         dockerRegistry,
		RegistryUser:     dockerUser,
		RegistryPassword: dockerPassword,
		GitSHA:           md.GitSHA,
		Name:             conf.Name,
		Tag:              md.GitBranch,
	})
	if err != nil {
		return fmt.Errorf("build Docker image: %w", err)
	}

	logger.Infof("Published %s", digest)

	return nil
}
