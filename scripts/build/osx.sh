#!/usr/bin/env bash
set -eu -o pipefail

BRANCH=$(git rev-parse --abbrev-ref HEAD)
GITREV=$(git rev-parse HEAD)
REV=${GITREV:0:7}-$BRANCH-$(date +%Y%m%d-%H:%M:%S)

CGO_ENABLED=0
GOOS=darwin
GOARCH=amd64

SOURCE="github.com/asnelzin/translate/cmd/translate"
TARGET="target/translate-$GOOS-$GOARCH"
echo "Building $TARGET"

go build -ldflags "-X main.revision=$REV" -o ./target/translate-$GOOS-$GOARCH "$SOURCE"