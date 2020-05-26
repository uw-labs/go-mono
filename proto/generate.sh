#!/usr/bin/env bash

# Enable logging of commands.
set -ex

# Get current directory.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Find all directories containing at least one prototfile.
# Based on: https://buf.build/docs/migration-prototool#prototool-generate.
for dir in $(find ${DIR}/uwlabs -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq); do
  files=$(find "${dir}" -name '*.proto')

  # Generate all files with go.
  protoc -I ${DIR} --go_out=plugins=grpc,paths=source_relative:${DIR}/gen/go ${files}

  # Generate with grpc-gateway and openapiv2 if any of the files imports the annotations proto.
  if grep -q 'import "google/api/annotations.proto";' ${files}; then
    protoc -I ${DIR} \
      --grpc-gateway_out=paths=source_relative:${DIR}/gen/go \
      --openapiv2_out=:${DIR}/gen/openapiv2 \
      ${files}
  fi
done
