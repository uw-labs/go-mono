package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/sirupsen/logrus"

	pkgctx "github.com/uw-labs/go-mono/pkg/context"
)

var (
	repoRoot     = flag.String("repo-root", ".", "The root of the repo, to find the git folder.")
	buildFile    = flag.String("build-file", "builds.txt", "The path to the build file to write release commands to.")
	moduleName   = flag.String("module-name", "github.com/uw-labs/go-mono", "The name of the local module.")
	baseRevision = flag.String("base", "", "The base revision to diff against when finding changes. Defaults to master.")
	headRevision = flag.String("head", "", "The head revision to diff with when finding changes. Defaults to HEAD.")
)

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	}

	err := run(logger, *repoRoot, *buildFile, *moduleName, *baseRevision, *headRevision)
	if err != nil {
		logger.WithError(err).Fatal()
	}
}

func run(logger *logrus.Logger, repoRoot, buildFilePath, moduleName, baseRevision, headRevision string) (err error) {
	buildFile, err := os.Create(buildFilePath)
	if err != nil {
		return fmt.Errorf("open build file: %w", err)
	}
	defer func() {
		cErr := buildFile.Close()
		if err == nil {
			err = cErr
		}
	}()

	ctx := pkgctx.WithSignalHandler(context.Background())

	revDeps, err := getReverseDependencies(ctx, logger, repoRoot, moduleName)
	if err != nil {
		return fmt.Errorf("get reverse dependencies: %w", err)
	}

	files, err := getChangedFiles(ctx, logger, repoRoot, baseRevision, headRevision)
	if err != nil {
		return fmt.Errorf("get changed files: %w", err)
	}

	logger.Infoln("Changed files:", files)

	folders := getFolders(files...)

	logger.Infoln("Changed folders:", folders)

	packages := getPackages(folders...)

	logger.Infoln("Changed packages:", packages)

	releases := map[string]string{}
	for _, pkg := range packages {
		for _, dependant := range revDeps[pkg] {
			if _, ok := releases[dependant]; ok {
				// Already added
				continue
			}

			df1 := filepath.Join(repoRoot, dependant, "deploy.yml")
			df2 := filepath.Join(repoRoot, dependant, "deploy.yaml")
			_, err1 := os.Stat(df1)
			_, err2 := os.Stat(df2)
			if os.IsNotExist(err1) && os.IsNotExist(err2) {
				// Binaries without a deploy file are not released
				continue
			}
			if err1 == nil {
				releases[dependant] = df1
				continue
			}
			if err2 == nil {
				releases[dependant] = df2
				continue
			}

			return fmt.Errorf("stat deployment file: %v, %w", err1, err2)
		}
	}

	for release, conf := range releases {
		if filepath.Dir(release) == "." {
			// Skip top level release
			continue
		}
		logger.Infoln("Release", release)
		_, err = buildFile.WriteString(conf + "\n")
		if err != nil {
			return fmt.Errorf("write releases to buildFile: %w", err)
		}
	}

	return nil
}

func getChangedFiles(ctx context.Context, logger *logrus.Logger, repoRoot, baseRevision, headRevision string) ([]string, error) {
	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("open local repository: %w", err)
	}

	baseRef, err := repo.Reference(plumbing.Master, true)
	if err != nil {
		return nil, fmt.Errorf("get master: %w", err)
	}
	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("get HEAD: %w", err)
	}

	base := baseRef.Hash()
	if baseRevision != "" {
		baseHash := plumbing.NewHash(baseRevision)
		err = repo.Storer.HasEncodedObject(baseHash)
		if err != nil {
			logger.WithError(err).Errorf("Failed to find base revision %q in the repository", baseRevision)
		} else {
			base = baseHash
		}
	}

	head := headRef.Hash()
	if headRevision != "" {
		headHash := plumbing.NewHash(headRevision)
		err = repo.Storer.HasEncodedObject(headHash)
		if err != nil {
			logger.WithError(err).Errorf("Failed to find head revision %q in the repository", headRevision)
		} else {
			head = headHash
		}
	}

	baseCommit, err := repo.CommitObject(base)
	if err != nil {
		return nil, fmt.Errorf("get base commit: %w", err)
	}

	headCommit, err := repo.CommitObject(head)
	if err != nil {
		return nil, fmt.Errorf("get head commit: %w", err)
	}

	logger.Infoln("Comparing from:" + baseCommit.Hash.String() + " to:" + headCommit.Hash.String())

	baseTree, err := baseCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("get base tree: %w", err)
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("get head tree: %w", err)
	}

	changes, err := baseTree.DiffContext(ctx, headTree)
	if err != nil {
		return nil, fmt.Errorf("get diff: %w", err)
	}

	var files []string
	for _, change := range changes {
		// Ignore deleted files
		action, err := change.Action()
		if err != nil {
			return nil, fmt.Errorf("get diff action: %w", err)
		}

		if action == merkletrie.Delete {
			continue
		}

		name := change.To.Name
		if change.From.Name != "" {
			name = change.From.Name
		}
		files = append(files, name)
	}

	return files, nil
}

func getFolders(files ...string) []string {
	set := map[string]struct{}{}
	var folders []string // nolint: prealloc
	for _, f := range files {
		dir := filepath.Dir(f)
		if _, ok := set[dir]; ok {
			continue
		}
		set[dir] = struct{}{}
		folders = append(folders, dir)
	}

	return folders
}

func getPackages(folders ...string) []string {
	var packages []string // nolint: prealloc
	for _, dir := range folders {
		switch dir {
		// Skip top level, not a package
		case ".":
			continue
		// Skip vendor, not a package
		case "vendor":
			continue
		// Skip CI directories
		case ".circleci", ".dependabot", ".github", ".github/workflows":
			continue
		}

		dir = strings.TrimPrefix(dir, "vendor/")

		packages = append(packages, dir)
	}

	return packages
}

type jsonPackage struct {
	ImportPath string
	Name       string
	Deps       []string
}

// getReverseDependencies gets all the executable (package main) dependencies
// of every package in the repo.
func getReverseDependencies(ctx context.Context, logger *logrus.Logger, repoRoot, moduleName string) (map[string][]string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("find go binary: %w", err)
	}
	cmd := exec.CommandContext(ctx, goBin, "list", "-json", "./...")
	cmd.Dir = repoRoot

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stdout: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stderr: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		drainClose(logger, stdOut, stdErr)
		return nil, fmt.Errorf("run go list: %w", err)
	}

	revDeps := map[string][]string{}
	dec := json.NewDecoder(stdOut)
	for dec.More() {
		var pkg jsonPackage
		err = dec.Decode(&pkg)
		if err != nil {
			drainClose(logger, stdOut, stdErr)
			return nil, fmt.Errorf("parse go list output: %w", err)
		}

		if pkg.Name != "main" {
			// Only care about main packages
			continue
		}

		// Trim module path from local packages, to match locally changed files
		pkgName := strings.TrimPrefix(pkg.ImportPath, moduleName)

		// Add the package as a dependency of itself, so that if only the package itself has
		// changed, we still build it
		revDeps[pkgName] = append(revDeps[pkgName], pkgName)

		for _, dep := range pkg.Deps {
			// Strip vendor prefixes from packages vendored by stdlib
			dep = strings.TrimPrefix(dep, "vendor/")

			if !strings.Contains(strings.Split(dep, "/")[0], ".") {
				// Skip packages without hostname in first element
				// (standard library/local imports)
				continue
			}

			// Trim module path from local packages, to match locally changed files
			dep = strings.TrimPrefix(dep, moduleName)

			revDeps[dep] = append(revDeps[dep], pkgName)
		}
	}

	err = cmd.Wait() // Closes "stdOut" and "stdErr"
	if err != nil {
		return nil, fmt.Errorf("run go list: %w", err)
	}

	return revDeps, nil
}

func drainClose(logger *logrus.Logger, readers ...io.ReadCloser) {
	for _, reader := range readers {
		bts, err := ioutil.ReadAll(reader)
		if err != nil {
			_ = reader.Close()
			continue
		}
		logger.Info(string(bts))
		_ = reader.Close()
	}
}
