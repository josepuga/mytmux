#!/bin/bash
# By José Puga. 2025.
set -euo pipefail
app_name=mytmux
# Use `git tag -a vX.Y.Z -m "Release v.X.Y.Z` to create a tag
version=$(git describe --tags)
ldflags="-s -w -X main.VERSION=${version}"
go fmt && go mod tidy && \
go build "${ldflags}" -o "bin/${app_name}"