# Logur adapter for [Logrus](https://github.com/sirupsen/logrus)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/logur/adapter-logrus/CI?style=flat-square)](https://github.com/logur/adapter-logrus/actions?query=workflow%3ACI)
[![Codecov](https://img.shields.io/codecov/c/github/logur/adapter-logrus?style=flat-square)](https://codecov.io/gh/logur/adapter-logrus)
[![Go Report Card](https://goreportcard.com/badge/logur.dev/adapter/logrus?style=flat-square)](https://goreportcard.com/report/logur.dev/adapter/logrus)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.11-61CFDD.svg?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/logur.dev/adapter/logrus)


## Installation

```bash
go get logur.dev/adapter/logrus
```


## Usage

```go
package main

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
)

func main() {
	logger := logrusadapter.New(logrus.New())
}
```


## Development

When all coding and testing is done, please run the test suite:

```bash
$ make check
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
