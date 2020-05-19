package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/uw-labs/go-mono/cmd/deploy/internal/docker/static"
)

//go:generate go-bindata -pkg static -prefix static -nometadata -ignore bindata -o ./static/bindata.go ./static

// The streamed response from the Docker build/push
// contain one or more of the following.
//
// There really is no documentation for this.
// Take a look:
//		https://docs.docker.com/engine/api/v1.40/#operation/ImageBuild
//		https://docs.docker.com/engine/api/v1.40/#operation/ImagePush
// Nothing.
type dockerResp struct {
	Stream string `json:"stream,omitempty"`
	Aux    struct {
		ID     string `json:"ID,omitempty"`
		Tag    string `json:"Tag,omitempty"`
		Digest string `json:"Digest,omitempty"`
		Size   int    `json:"Size,omitempty"`
	} `json:"aux,omitempty"`
	Status         string `json:"status,omitempty"`
	ID             string `json:"id,omitempty"`
	ProgressDetail struct {
		Current float64
		Total   float64
	} `json:"progressDetail,omitempty"`
	Progress string `json:"progress,omitempty"`
}

// Request is the input to BuildAndPushImage
type Request struct {
	RepoRoot         string
	BinaryPath       string
	Registry         string
	RegistryUser     string
	RegistryPassword string
	Name             string
	GitSHA           string
	Tag              string
}

// BuildAndPushImage builds a docker image using the default Dockerfile and pushes
// it to the partner registry. It returns the registry SHA256 digest of the pushed image.
func BuildAndPushImage(ctx context.Context, logger *logrus.Logger, req *Request) (digest string, err error) {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("connect to docker: %w", err)
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	w := tar.NewWriter(gw)

	dockerfile := static.MustAsset("Dockerfile")
	err = w.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Mode: 0o400,
		Size: int64(len(dockerfile)),
	})
	if err != nil {
		return "", fmt.Errorf("create tar header for Dockerfile: %w", err)
	}

	_, err = w.Write(dockerfile)
	if err != nil {
		return "", fmt.Errorf("write Dockerfile to tar: %w", err)
	}

	f, err := os.Open(req.BinaryPath)
	if err != nil {
		return "", fmt.Errorf("open binary for reading: %w", err)
	}
	defer func() {
		cErr := f.Close()
		if err == nil {
			err = cErr
		}
	}()

	finfo, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("stat binary: %w", err)
	}

	err = w.WriteHeader(&tar.Header{
		Name: "app",
		Mode: 0o500,
		Size: finfo.Size(),
	})

	_, err = io.Copy(w, f)
	if err != nil {
		return "", fmt.Errorf("copy binary: %w", err)
	}

	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("close tar writer: %w", err)
	}

	err = gw.Close()
	if err != nil {
		return "", fmt.Errorf("close gzip writer: %w", err)
	}

	pr, pw := io.Pipe()
	tags := []string{
		req.Registry + "/" + req.Name + ":" + req.GitSHA,
		req.Registry + "/" + req.Name + ":" + req.Tag,
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer pw.Close()
		resp, err := client.ImageBuild(egCtx, &buf, types.ImageBuildOptions{
			Labels: map[string]string{
				"revision": req.GitSHA,
			},
			Tags: tags,
		})
		if err != nil {
			return fmt.Errorf("build docker image: %w", err)
		}
		defer resp.Body.Close()
		// Stream response
		_, err = io.Copy(pw, resp.Body)
		if err != nil {
			return fmt.Errorf("read docker build response: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		defer pr.Close()
		dec := json.NewDecoder(pr)
		for dec.More() {
			var msg dockerResp
			err = dec.Decode(&msg)
			if err != nil {
				return fmt.Errorf("parse piped docker build response: %w", err)
			}

			if msg.Stream != "" {
				logger.Println(strings.TrimSpace(msg.Stream))
			}
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return "", err
	}

	auth := types.AuthConfig{
		Username: req.RegistryUser,
		Password: req.RegistryPassword,
	}
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("marshal docker auth: %w", err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	for _, image := range tags {
		eg, egCtx = errgroup.WithContext(ctx)

		pr, pw = io.Pipe()

		eg.Go(func() error {
			defer pw.Close()
			resp, err := client.ImagePush(egCtx, image, types.ImagePushOptions{
				RegistryAuth: authStr,
			})
			if err != nil {
				return fmt.Errorf("failed to push docker image (%s): %w", image, err)
			}
			defer resp.Close()
			// Stream response
			_, err = io.Copy(pw, resp)
			if err != nil {
				return fmt.Errorf("read docker push response: %w", err)
			}
			return nil
		})
		eg.Go(func() error {
			defer pr.Close()
			dec := json.NewDecoder(pr)
			for dec.More() {
				var msg dockerResp
				err = dec.Decode(&msg)
				if err != nil {
					return fmt.Errorf("parse piped docker push response: %w", err)
				}

				if msg.Status != "" {
					logger.Println(strings.TrimSpace(msg.Status))
				}

				if msg.Aux.Digest != "" {
					digest = fmt.Sprintf("%s/%s@%s", req.Registry, req.Name, msg.Aux.Digest)
				}
			}

			return err
		})

		err = eg.Wait()
		if err != nil {
			return "", err
		}
	}

	return digest, nil
}
