#!/usr/bin/env bash
# DeepPhotos — Build Docker Images
# Usage: ./scripts/mac-linux/docker_build.sh
set -e

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
LOG="$ROOT/docker_build.log"

echo "🐳  Building DeepPhotos Docker Images..."
echo "    Log: $LOG"
echo ""

docker compose -f "$ROOT/docker/docker-compose.yaml" build --no-cache 2>&1 | tee "$LOG"

echo ""
echo "✅  Build complete!"
echo "    Run:  ./scripts/mac-linux/docker_run.sh"
