# Podrick

[![CircleCI](https://img.shields.io/circleci/project/github/uw-labs/podrick/master.svg?style=flat-square)](https://circleci.com/gh/uw-labs/podrick)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/uw-labs/podrick)
[![Go Report Card](https://goreportcard.com/badge/github.com/uw-labs/podrick?style=flat-square)](https://goreportcard.com/report/github.com/uw-labs/podrick)
[![Code Coverage](https://img.shields.io/codecov/c/github/uw-labs/podrick/master.svg?style=flat-square)](https://codecov.io/gh/uw-labs/podrick)
[![Releases](https://img.shields.io/github/release/uw-labs/podrick.svg?style=flat-square)](https://github.com/uw-labs/podrick/releases)
[![License](https://img.shields.io/github/license/uw-labs/podrick.svg?style=flat-square)](LICENSE)

Dynamically create and destroy containers for tests, within
your Go application. Support for both [Podman](https://podman.io)
and [Docker](https://docker.com) runtimes is built-in.

Inspired by [dockertest](https://github.com/ory/dockertest).

## Usage

```go
package some_test

import (
	"net/http"
	"testing"

	"github.com/uw-labs/podrick"
	_ "github.com/uw-labs/podrick/runtimes/docker" // Register the docker runtime.
	_ "github.com/uw-labs/podrick/runtimes/podman" // Register the podman runtime.
)

func TestDatabase(t *testing.T) {
	ctx := context.Background()

	ctr, err := podrick.StartContainer(ctx, "kennethreitz/httpbin", "latest", "80")
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		err := ctr.Close(ctx) // Stops and removes the container.
		if err != nil {
			t.Error(err.Error())
		}
	}()

	// Use ctr.Address().
	resp, err = http.Get("http://" + ctr.Address() + "/get")
}
```

Podrick automatically selects a runtime if one isn't explicitly specified.
This allows the user to use both podman and docker, depending on the support
of the environment. It's also possible to explicitly specify which runtime to use,
or use a custom runtime implementation.

## Advanced usage

```go
package some_test

import (
	"database/sql"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/uw-labs/podrick"
	"github.com/uw-labs/podrick/runtimes/podman"
	logrusadapter "logur.dev/adapters/logrus"
)

func TestDatabase(t *testing.T) {
	log := logrusadapter.New(logrus.New())
	ctx := context.Background()

	ctr, err := podrick.StartContainer(ctx, "cockroachdb/cockroach", "v19.1.3", "26257",
		podrick.WithUlimit([]podrick.Ulimit{{
			Name: "nofile",
			Soft: 1956,
			Hard: 1956,
		}}),
		podrick.WithCmd([]string{
			"start",
			"--insecure",
		}),
		podrick.WithLogger(log),
		// Use of the podman runtime only.
		// The environment variable PODMAN_VARLINK_ADDRESS
		// can be used to configure where podrick should
		// look for the varlink API.
		podrick.WithRuntime(&podman.Runtime{
			Logger: log,
		}),
		podrick.WithLivenessCheck(func(address string) error {
			_, err := http.Get("http://"+address+"/get")
			return err
		}),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {
		err := ctr.Close(ctx)
		if err != nil {
			t.Error(err.Error())
		}
	}()

	db, err := sql.Open("postgres://root@" + ctr.Address() + "/defaultdb")
	// ... make database calls, etc
}
```

## Using podrick in CI

While `podrick` makes it really easy to run tests locally on users
machines, it can be challenging to run the tests in CI, since
it requires the permission to run containers that are accessible
to the local environment. Circle CI is the only environment confirmed
to work well with `podrick`, both using `podman` and `docker`.

### Circle CI

The `podrick` CI tests themselves illustrate how to use `podrick` with Circle CI.
It requires the use of the `machine` executor type, and there is some setup required
to get a recent Go version. Please see [the CI files](./circleci/config.yml) for
concrete examples of this.

If you are using Circle CI with the `docker` runtime, it is obviously not
necessary to install podman.
