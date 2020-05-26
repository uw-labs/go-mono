format:
	go install mvdan.cc/gofumpt/gofumports
	grep -L -R "^\/\/ Code generated .* DO NOT EDIT\.$$" --exclude-dir=.git --exclude-dir=vendor --include="*.go" . | xargs -n 1 gofumports -w -local github.com/uw-labs/go-mono

lint-imports:
	go install github.com/pavius/impi/cmd/impi
	impi --local github.com/uw-labs/go-mono --scheme stdThirdPartyLocal --ignore-generated=true ./...

generate:
	./proto/generate.sh
	go generate -x ./...

install-generators:
	go install \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/tmthrgd/go-bindata/go-bindata
