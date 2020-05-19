package podrick

// Option configures the configuration of the started container.
type Option func(*config)

// WithEnv configures the environment of the container.
func WithEnv(in []string) Option {
	return func(c *config) {
		c.ContainerConfig.Env = in
	}
}

// WithEntrypoint configures the entrypoint of the container.
func WithEntrypoint(in string) Option {
	return func(c *config) {
		c.ContainerConfig.Entrypoint = &in
	}
}

// WithCmd configures the command of the container.
func WithCmd(in []string) Option {
	return func(c *config) {
		c.ContainerConfig.Cmd = in
	}
}

// WithUlimit configures the ulimits of the container.
func WithUlimit(in []Ulimit) Option {
	return func(c *config) {
		c.ContainerConfig.Ulimits = in
	}
}

// WithLogger configures the logger of the container.
// The containers logs will be logged at Info level to this logger.
// Some errors during closing may also be logged at Error level.
func WithLogger(in Logger) Option {
	return func(c *config) {
		c.logger = in
	}
}

// WithRuntime configures the Runtime to use to launch the container.
// By default, the auto runtime is used.
func WithRuntime(in Runtime) Option {
	return func(c *config) {
		c.runtime = in
	}
}

// WithLivenessCheck defines a function to call repeatedly until it does not
// error, to ascertain the successful startup of the container. The
// function will be retried for 10 seconds, and if it does not return
// a non-nil error before that time, the last error will be returned.
func WithLivenessCheck(lc LivenessCheck) Option {
	return func(c *config) {
		c.liveCheck = lc
	}
}

// WithFileUpload writes the content of the reader to the provided path
// inside the container, before starting the container. This can
// be specified multiple times.
func WithFileUpload(f File) Option {
	return func(c *config) {
		c.Files = append(c.Files, f)
	}
}

// WithExposePort adds extra ports that should be exposed from the
// started container.
func WithExposePort(port string) Option {
	return func(c *config) {
		c.ExtraPorts = append(c.ExtraPorts, port)
	}
}

// LivenessCheck is a type used to check the successful startup
// of a container.
type LivenessCheck func(address string) error

type config struct {
	ContainerConfig

	logger    Logger
	runtime   Runtime
	liveCheck LivenessCheck
}
