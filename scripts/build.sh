#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY_NAME="agent-forensic"

# Read version
VERSION="$(cat "$ROOT_DIR/version" 2>/dev/null || echo "dev")"
LDFLAGS="-X main.Version=$VERSION"

PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

cd "$ROOT_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUT="$ROOT_DIR/bin/${GOOS}-${GOARCH}/${BINARY_NAME}"
    [[ "$GOOS" == "windows" ]] && OUT="${OUT}.exe"
    mkdir -p "$(dirname "$OUT")"
    GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "$LDFLAGS" -o "$OUT" .
    echo "Built: bin/${GOOS}-${GOARCH}/${BINARY_NAME}$([[ "$GOOS" == "windows" ]] && echo ".exe") v${VERSION}"
done
