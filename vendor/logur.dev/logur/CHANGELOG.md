# Change Log


All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased]


## [0.16.1] - 2020-01-16

### Fixed

- Log event assertion output


## [0.16.0] - 2020-01-15

### Added

- `LoggerContext` interface
- `LoggerFacade` interface (combination of `Logger` and `LoggerContext`)
- `LoggerContextFunc` logger function wrapper
- `NoopHandler` no-op error handler
- `TestLoggerContext`, `TestLoggerSet` test error handlers

### Changed

- Improved conformance tests

### Fixed

- gRPC format log functions

### Deprecated

- `NewNoopLogger` no-op logger. Use `NoopLogger` instead.
- `NewTestHandler` test handler factory. Use `TestHandler` instead.
- `LoggerTestSuite` test suite. Use the new `conformance` package instead.
- `LogEventsEqual` function. Use `LogEvent.Equals` and `LogEvent.AssertEquals` instead.


## [0.15.1] - 2019-11-12

### Added

- `WithField` as a shortcut for `WithFields`
- Github Actions workflow


## [0.15.0] - 2019-08-22

### Changed

- Import path to `logur.dev/logur`


## [0.14.0] - 2019-08-22

### Removed

- Adapter implementations. Use the ones from the new [organization](https://github.com/logur?utf8=%E2%9C%93&q=adapter-&type=&language=)
- Integration implementations (with external dependencies). Use the ones from the new [organization](https://github.com/logur?utf8=%E2%9C%93&q=integration-&type=&language=)


## [0.13.0] - 2019-08-22

### Deprecated

- Adapter implementations. Use the ones from the new [organization](https://github.com/logur?utf8=%E2%9C%93&q=adapter-&type=&language=)
- Integration implementations (with external dependencies). Use the ones from the new [organization](https://github.com/logur?utf8=%E2%9C%93&q=integration-&type=&language=)


## [0.12.0] - 2019-08-16

### Changed

- Renamed `ContextualLogger` to `fieldLogger`
- Examples are moved to a separate module

### Removed

- Error handler (use [emperror.dev/handler/logur](https://emperror.dev/handler/logur) instead)


## [0.11.2] - 2019-07-18

### Fixed

- Minimum Logrus version ([#49](https://github.com/logur/logur/pull/49))


## [0.11.1] - 2019-07-10

### Added

- `logrus`: `NewFromEntry` method to create a logger from a custom entry


## [0.11.0] - 2019-02-26

### Added

- [zap](https://github.com/uber-go/zap) logger integration

### Changed

- Renamed `testing` package directory to `logtesting`


## [0.10.0] - 2019-02-08

### Added

- Separate interface for error logging
- Error handler interface to Watermill

### Changed

- Update Watermill logger to prepare for the next version
- Export the Watermill logger type
- Export the Invision logger type


## [0.9.0] - 2019-01-10

### Added

- [logr](https://github.com/go-logr/logr) integration

### Changed

- Make the log context map optional


## [0.8.0] - 2018-12-29

### Added

- Constructor for standard logger for errors
- `PrintLogger` that logs messages using `fmt.Print*` semantics

### Changed

- Renamed `logtesting.AssertLogEvents` to `AssertLogEventsEqual`
- Renamed `AssertLogEventsEqual` to `LogEventsEqual`

### Removed

- [MySQL driver](https://github.com/go-sql-driver/mysql) integration (use `PrintLogger` instead)


## [0.7.1] - 2018-12-22

### Added

- Simplified message logger without contextual logging
- Some tests for integrations to ensure interface compatibility


## [0.7.0] - 2018-12-21

### Added

- Public test log event comparison function
- Example package

### Changed

- Exported the log testing library so it can be used for testing in other libraries
- Unexport noop logger


## [0.6.0] - 2018-12-21

### Added

- Contextual logger (instead of `Logger.WithFields`)
- Field parameter to log functions
- [gRPC log](https://godoc.org/google.golang.org/grpc/grpclog) integration
- [MySQL driver](https://github.com/go-sql-driver/mysql) integration

### Changed

- Replace log func variadic arguments with a single message argument
- Check if level is enabled (to prevent unwanted context conversions) when the underlying logger supports it
- Export all log adapter types (in accordance with [Go interface](https://github.com/golang/go/wiki/CodeReviewComments#interfaces) guidelines)

### Removed

- format functions from `Logger` interface
- ln functions from `Logger` interface
- Simple log adapter (implementing format and ln functions)
- `Logger.WithFields` method (use field parameter of log functions instead)


## [0.5.0] - 2018-12-17

### Added

- [Watermill](https://watermill.io) compatible logger

### Changed

- Dropped the custom `Fields` type from the `Logger` interface (replaced with `map[string]interface{}`)


## [0.4.0] - 2018-12-11

### Added

- Benchmarks
- [github.com/rs/zerolog](https://github.com/rs/zerolog) adapter
- [github.com/go-kit/kit](https://github.com/go-kit/kit) adapter


## [0.3.0] - 2018-12-11

### Added

- [github.com/goph/emperror](https://github.com/goph/emperror) compatible error handler
- Uber Zap adapter

### Changed

- Removed *Level* suffix from level constants


## [0.2.0] - 2018-12-10

### Added

- [github.com/InVisionApp/go-logger](https://github.com/InVisionApp/go-logger) integration
- `simplelogadapter` to make logger library integration easier
- [github.com/hashicorp/go-hclog](https://github.com/hashicorp/go-hclog) adapter

### Changed

- Renamed `logrusshim` to `logrusadapter`

## 0.1.0 - 2018-12-09

- Initial release


[Unreleased]: https://github.com/logur/logur/compare/v0.16.1...HEAD
[0.16.1]: https://github.com/logur/logur/compare/v0.16.0...v0.16.1
[0.16.0]: https://github.com/logur/logur/compare/v0.15.1...v0.16.0
[0.15.1]: https://github.com/logur/logur/compare/v0.15.0...v0.15.1
[0.15.0]: https://github.com/logur/logur/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/logur/logur/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/logur/logur/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/logur/logur/compare/v0.11.2...v0.12.0
[0.11.2]: https://github.com/logur/logur/compare/v0.11.1...v0.11.2
[0.11.1]: https://github.com/logur/logur/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/logur/logur/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/logur/logur/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/logur/logur/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/logur/logur/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/logur/logur/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/logur/logur/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/logur/logur/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/logur/logur/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/logur/logur/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/logur/logur/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/logur/logur/compare/0.1.0...v0.2.0
