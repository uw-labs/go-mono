[![CircleCI](https://circleci.com/gh/uw-labs/go-mono.svg?style=shield&circle-token=53ab4342cde1e547f400c27d21dbc3e8cd9de66f)](https://circleci.com/gh/uw-labs/go-mono)

# Utility Warehouse template Go monorepo

This repo is an abbreviated copy of one used by one of the teams inside Utility Warehouse.
It's been built for tight CI integration and developer productivity.

## Making it your own

If you're creating a new repo from this template, you'll want to do a search-replace on
`github.com/uw-labs/go-mono` and `uwlabs` (replacing `yourorg`, `yourrepo` and `yourscm`
with your own):

```shell
$ find . -type d -name "uwlabs" -exec sh -c 'mv {} $(dirname {})/yourorg' \;
$ find . -type f -exec sed -i 's;github.com/uw-labs/go-mono;yourscm.com/yourorg/yourrepo;g' {} +
$ find . -type f -exec sed -i 's;uwlabs;yourorg;g' {} +
```

Configure `DOCKER_USER` and `DOCKER_PASSWORD` in your CircleCI settings for release
builds to work. The registry path to use is configured in the `release` CI job.

## CI Setup

The CI setup uses Circle CI, partly because of the powerful caching features, and partly
because Circle CI machine users have access to a local docker socket, which
allows us to run integration tests using `go test` which work both locally and in CI.
The CI jobs rely on Go build and test caching to be efficient even with a larger number of services.
The same CI setup is used to efficiently test and publish over 130 applications within UW.

### Linters

* [golangci-lint](https://github.com/golangci/golangci-lint) for Go code.
* [buf](https://github.com/bufbuild/buf) for protobuf linting and breaking change detection.
* [gofumports](https://github.com/mvdan/gofumpt) for formatting and imports.

### Automatic publishing

The `release` CI job automatically figures out what needs building and publishes a docker
image to the configured registry. It requires the setting of `DOCKER_USER` and `DOCKER_PASSWORD`
in the Circle CI configuration environment variables.

The Dockerfile used to build the images is [here](./cmd/deploy/internal/docker/static/Dockerfile).
It can be edited as necessary, just make sure to run `make generate` after changing it.

For an example of this, the [user-api](./cmd/user-api/main.go) is published automatically to
[the local GitHub docker registry](https://github.com/uw-labs/go-mono/packages/237911)
whenever it requires rebuilding.

## Repository Layout

* [cmd](cmd) - Utilities and service applications.
  * [cmd/calculate-releases](cmd/calculate-releases/main.go) Script for calculating applications to publish.
  * [cmd/deploy](cmd/deploy/main.go) Script for building and publishing docker images.
  * [cmd/user-api](cmd/user-api/main.go) Example application with deploy.yml.
* [pkg](pkg) - Shared packages.
* [proto](proto) - Protobuf definitions & generated code.
* [vendor](vendor) - Vendored third-party dependencies.

## Makefile
A top level [Makefile](./Makefile) exists to help you perform common actions within the
monorepo. Recipes include:

* `format` - Formats your `.go` source code
* `install-generators` - Installs all necessary generators for the `generate` step
* `generate` - Runs all generators required within the monorepo
* `lint-imports` - Runs the import linter.

## How do I add a new application?

1. Create a new folder in `cmd` for your service
   E.g. `cmd/my-new-service`.
1. Create a `main.go`
1. If you want to automatically build a docker container, add a `deploy.yml` file.

### Adding protofiles

Add your own protofiles under [proto/uwlabs/](proto/uwlabs). Follow the folder
structure laid out in [the buf documentation](https://buf.build/docs/style-guide#files-and-packages).

## The calculate-releases script

[calculate-releases](./cmd/calculate-releases/main.go) is the magic that calculates
exactly what applications need to be rebuilt based on file changes to the application
directly, or any of its transative dependencies. It is called in CI on every branch push,
to calculate which applications to build docker images for.

It relies on `go list -json ./...` for dependency information.

## The deploy script

[deploy](./cmd/deploy/main.go) is responsible for building and publishing
a docker image to the specified registry. It's usually run on the output of
the `calculate-releases` script. It defines a custom format for build configuration
and metadata. This can be extended to include things such as kubernetes
deployment targets, extra application metadata and more. It is currently run
automatically against every branch push in CI.

### The deploy.yml file

Use a `deploy.yml` together with any main packages that you want to deploy
to configure automatic docker container building and publishing. The `deploy.yml`
file allows a single configuration parameter:

* `name`

   Used to configure the name of the docker image pushed to the registry.

## Why a vendor directory?

When evaluating solutions to two problems, the vendor directory became the primary
candidate:

1. How do we keep CI builds as fast as possible?

   We first implemented this using module caching, where
   the first job would download all the modules and cache them for
   future jobs. It meant a lot of extra boilerplate in the CircleCI
   configuration files, and it never worked well for the Docker builds.
   Vendoring means we have all the source code available at all time,
   and completely removes the need for caching. This sped up builds by
   roughly 50% in testing.

1. How do we ensure we only release applications that have changed
   when we perform dependency updates?

   This could be done with some custom tooling that can discover file changes
   between module updates, but this is nontrivial, and we already had
   a solution that worked with the vendor directory.

The CI pipeline will fail on your pull request if you add a new dependency that is not
within the `/vendor` directory. If you've added a new dependency, make sure you run `go mod tidy`
and `go mod vendor` to ensure your dependencies are up-to-date and vendored.
