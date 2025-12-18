#!/bin/bash
# By José Puga. 2025.
set -euo pipefail
app_name=mytmux
# To create a version tag:
#   Use `git tag -a vX.Y.Z -m "Release v.X.Y.Z`
#   (a git commit is necessary first). Then `git push --tags`
version=$(git describe --tags)
ldflags="-s -w -X main.VERSION=${version}"
go fmt && go mod tidy && \
go build -ldflags="${ldflags}" -o "bin/${app_name}"