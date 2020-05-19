//+build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/pavius/impi/cmd/impi"
	_ "github.com/tmthrgd/go-bindata/go-bindata"
	_ "mvdan.cc/gofumpt/gofumports"
)
