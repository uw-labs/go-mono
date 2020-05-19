package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Deployment describes a parsed deploy.yaml file
type Deployment struct {
	Main string `yaml:"main"`
	Name string `yaml:"name"`
}

// Parse parses the deploy.yaml file at the path
func Parse(repoRoot, path string) (_ *Deployment, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open the deploy file: %w", err)
	}
	defer func() {
		cErr := f.Close()
		if err == nil {
			err = cErr
		}
	}()

	name := filepath.Dir(path)

	dc := Deployment{
		Main: name,
	}
	err = yaml.NewDecoder(f).Decode(&dc)
	if err != nil {
		return nil, fmt.Errorf("parse the deploy file: %w", err)
	}

	return &dc, nil
}
