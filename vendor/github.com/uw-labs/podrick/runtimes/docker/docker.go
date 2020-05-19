package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"logur.dev/logur"

	"github.com/uw-labs/podrick"
)

func init() {
	podrick.RegisterAutoRuntime(&Runtime{})
}

// Runtime implements the Runtime interface with
// a Docker runtime backend.
//
// Supported environment variables:
// DOCKER_HOST to set the url to the docker server.
// DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
// DOCKER_CERT_PATH to load the TLS certificates from.
// DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
type Runtime struct {
	Logger podrick.Logger

	client *docker.Client
}

// Connect connects to the Docker API.
func (r *Runtime) Connect(ctx context.Context) error {
	if r.Logger == nil {
		r.Logger = logur.NewNoopLogger()
	}

	var err error
	r.client, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to connect to docker: %w", err)
	}
	_, err = r.client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping docker: %w", err)
	}

	return nil
}

// Close releases the resources of the Runtime.
func (Runtime) Close(context.Context) error {
	return nil
}

// StartContainer starts a container with Docker as the backing runtime.
func (r *Runtime) StartContainer(ctx context.Context, conf *podrick.ContainerConfig) (podrick.Container, error) {
	ctr := &container{
		runtime: r,
	}
	_, _, err := r.client.ImageInspectWithRaw(ctx, conf.Repo+":"+conf.Tag)
	if err != nil {
		bd, err := r.client.ImagePull(ctx, conf.Repo+":"+conf.Tag, types.ImagePullOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to pull image: %w", err)
		}
		_, err = io.Copy(logur.NewWriter(r.Logger), bd)
		if err != nil {
			return nil, fmt.Errorf("failed to stream image: %w", err)
		}
		err = bd.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close pull body: %w", err)
		}
	}

	cc, hc, nc := createConfig(conf)
	resp, err := r.client.ContainerCreate(ctx, cc, hc, nc, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	if len(conf.Files) > 0 {
		err = uploadFiles(ctx, r.client, resp.ID, conf.Files...)
		if err != nil {
			return nil, fmt.Errorf("failed to upload files to container: %w", err)
		}
	}

	if err := r.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	ctr.close = func(ctx context.Context) error {
		return r.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
	}

	ctr.container, err = r.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	if ctr.container.NetworkSettings == nil {
		return nil, fmt.Errorf("failed to get container network")
	}

	ctr.portToaddress = make(map[string]string)
	for addr, hostPorts := range ctr.container.NetworkSettings.Ports {
		for _, port := range hostPorts {
			// Will use the last one, don't care for now
			ctr.portToaddress[addr.Port()] = net.JoinHostPort(port.HostIP, port.HostPort)
		}
	}

	if ctr.portToaddress[conf.Port] == "" {
		return nil, fmt.Errorf("failed to get container address")
	}

	ctr.address = ctr.portToaddress[conf.Port]
	return ctr, nil
}

type container struct {
	address       string
	portToaddress map[string]string
	close         func(context.Context) error

	container types.ContainerJSON
	runtime   *Runtime
}

func (c *container) Address() string {
	return c.address
}

func (c container) AddressForPort(port string) (string, error) {
	hostPort, ok := c.portToaddress[port]
	if !ok {
		return "", fmt.Errorf("no address found for port %q", port)
	}
	return hostPort, nil
}

func (c *container) Close(ctx context.Context) error {
	return c.close(ctx)
}

func (c *container) StreamLogs(_ context.Context, w io.Writer) error {
	// Decoupled context from input context, since it controls logging lifetime.
	ctx, cancel := context.WithCancel(context.Background())
	body, err := c.runtime.client.ContainerLogs(ctx, c.container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		cancel()
		return fmt.Errorf("failed to connect to container log output: %w", err)
	}

	wg := &sync.WaitGroup{}

	cls := c.close
	c.close = func(ctx context.Context) error {
		cancel()
		wg.Wait() // Wait for goroutine to exit
		cErr := body.Close()
		if cErr != nil {
			c.runtime.Logger.Error("failed to close container logs", map[string]interface{}{
				"error": cErr.Error(),
			})
		}
		return cls(ctx)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(w, body)
		if err != nil {
			c.runtime.Logger.Error("failed to copy container logs", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}
