package podrick

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// Runtime supports starting containers.
type Runtime interface {
	Close(context.Context) error
	Connect(context.Context) error
	StartContainer(context.Context, *ContainerConfig) (Container, error)
}

// Container represents a running container.
type Container interface {
	// Context releases resources associated with the container.
	Close(context.Context) error
	// Address returns the IP and port of the running container.
	Address() string
	// AddressForPort returns the address for the specified port,
	// or an error, if the port was not exposed.
	AddressForPort(string) (string, error)
	// StreamLogs asynchronously streams logs from the
	// running container to the writer. The writer must
	// be safe for concurrent use.
	// If the context is cancelled after logging has been set up,
	// it has no effect. Use Close to stop logging.
	// This function is called automatically on the runtimes
	// configured logger, so there is no need to explicitly call this.
	StreamLogs(context.Context, io.Writer) error
}

var autoRuntimes []Runtime

// RegisterAutoRuntime allows a runtime to register itself
// for auto-selection of a runtime, when one isn't explicitly specified.
func RegisterAutoRuntime(r Runtime) {
	autoRuntimes = append(autoRuntimes, r)
}

type autoRuntime struct {
	Runtime
}

// Connect establishes a connection with the underlying runtime.
func (r *autoRuntime) Connect(ctx context.Context) error {
	if len(autoRuntimes) == 0 {
		return errors.New("no container runtimes registered, import one or choose explicitly")
	}

	var errs []error
	for _, r.Runtime = range autoRuntimes {
		err := r.Runtime.Connect(ctx)
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}

	errStr := "failed to automatically choose runtime:\n"
	for _, err := range errs {
		errStr += fmt.Sprintf("\t%q", err.Error())
	}

	return errors.New(errStr)
}
