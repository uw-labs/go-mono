package podrick

import (
	"context"
	"fmt"
	"time"

	backoff "github.com/cenkalti/backoff/v3"
	"logur.dev/logur"
)

// StartContainer starts a container using the configured runtime.
// By default, a runtime is chosen automatically from those registered.
func StartContainer(ctx context.Context, repo, tag, port string, opts ...Option) (_ Container, err error) {
	conf := config{
		ContainerConfig: ContainerConfig{
			Repo: repo,
			Tag:  tag,
			Port: port,
		},
		logger:  logur.NewNoopLogger(),
		runtime: &autoRuntime{},
	}
	for _, o := range opts {
		o(&conf)
	}

	err = conf.runtime.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to runtime: %w", err)
	}
	defer func() {
		if err != nil {
			cErr := conf.runtime.Close(context.Background())
			if cErr != nil {
				conf.logger.Error("failed to close runtime", map[string]interface{}{
					"error": cErr.Error(),
				})
			}
		}
	}()

	ctr, err := conf.runtime.StartContainer(ctx, &conf.ContainerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}
	defer func() {
		if err != nil {
			cErr := ctr.Close(context.Background())
			if cErr != nil {
				conf.logger.Error("failed to close container", map[string]interface{}{
					"error": cErr.Error(),
				})
			}
		}
	}()

	err = ctr.StreamLogs(ctx, logur.NewWriter(conf.logger))
	if err != nil {
		return nil, fmt.Errorf("failed to stream container logs: %w", err)
	}

	if conf.liveCheck != nil {
		bk := backoff.NewExponentialBackOff()
		bk.MaxElapsedTime = 10 * time.Second
		cbk := backoff.WithContext(bk, ctx)
		err = backoff.RetryNotify(
			func() error {
				return conf.liveCheck(ctr.Address())
			},
			cbk,
			func(err error, next time.Duration) {
				conf.logger.Error("Liveness check failed", map[string]interface{}{
					"retry_in": next.Truncate(time.Millisecond).String(),
					"error":    err.Error(),
				})
			},
		)
		if err != nil {
			return nil, fmt.Errorf("liveness check failed: %w", err)
		}
	}

	return ctr, nil
}
