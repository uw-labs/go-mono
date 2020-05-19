package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/uw-labs/podrick"
	"golang.org/x/sync/errgroup"
)

func uploadFiles(ctx context.Context, client *docker.Client, cID string, files ...podrick.File) error {
	r, w := io.Pipe()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		defer func() {
			cErr := w.Close()
			if err == nil {
				err = cErr
			}
		}()

		archive := tar.NewWriter(w)
		for _, f := range files {
			path := filepath.Clean(f.Path)
			if !filepath.IsAbs(path) {
				return fmt.Errorf("file paths must be absolute: %q", f.Path)
			}
			err = archive.WriteHeader(&tar.Header{
				Name: f.Path,
				Mode: int64(f.Mode),
				Size: int64(f.Size),
			})
			if err != nil {
				return fmt.Errorf("failed to write file header: %w", err)
			}
			_, err = io.Copy(archive, f.Content)
			if err != nil {
				return fmt.Errorf("failed to write file contents: %w", err)
			}
		}

		err = archive.Close()
		if err != nil {
			return fmt.Errorf("failed to write tar footer: %w", err)
		}

		return nil
	})

	eg.Go(func() (err error) {
		defer func() {
			cErr := r.Close()
			if err == nil {
				err = cErr
			}
		}()

		err = client.CopyToContainer(ctx, cID, "/", r, types.CopyToContainerOptions{})
		if err != nil {
			return fmt.Errorf("failed to copy files to container: %w", err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	return nil
}
