package podrick

import (
	"io"
	"os"
)

// ContainerConfig is used by runtimes to start
// containers.
type ContainerConfig struct {
	Repo string
	Tag  string
	Port string

	// Optional
	Env        []string
	Entrypoint *string
	Cmd        []string
	Ulimits    []Ulimit
	Files      []File
	ExtraPorts []string
}

// Ulimit describes a container ulimit.
type Ulimit struct {
	Name string
	Soft int64
	Hard int64
}

// File describes a file in a container.
// All fields are mandatory.
type File struct {
	Content io.Reader
	Path    string
	Size    int
	Mode    os.FileMode
}
