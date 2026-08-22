#!/usr/bin/env bash
# DeepPhotos — Local Development Environment
# Usage: ./scripts/mac-linux/dev.sh          → Normal dev mode
#        ./scripts/mac-linux/dev.sh --debug  → Dev mode with Delve (dlv) debugger on port 4000
set -e

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
ENABLE_DEBUG="${DEBUG:-0}"

if [ "$1" = "--debug" ] || [ "$1" = "-d" ]; then
    ENABLE_DEBUG=1
fi

echo "⚡  Starting DeepPhotos Local Development Environment..."
if [ "$ENABLE_DEBUG" = "1" ]; then
    echo "🐛  Delve Debugging Mode ACTIVE (Port 4000)"
fi
echo ""

# 1. Start MinIO Object Storage container via Docker Compose
if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
    echo "📦  Starting MinIO Object Storage container..."
    docker compose -f "$ROOT/docker/docker-compose.yaml" up -d minio
    echo "✅  MinIO active at http://localhost:9000 (Web Console: http://localhost:9001)"
else
    echo "ℹ️   Docker is not active — backend will use local disk storage fallback (backend/data/storage/)."
fi

echo ""

# 2. Stop any existing process running on port 8080 or port 4000
if lsof -pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "🔄  Stopping existing backend process on port 8080..."
    lsof -ti:8080 | xargs kill -9 >/dev/null 2>&1 || true
    sleep 1
fi

if [ "$ENABLE_DEBUG" = "1" ] && lsof -pi :4000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "🔄  Stopping existing Delve debugger process on port 4000..."
    lsof -ti:4000 | xargs kill -9 >/dev/null 2>&1 || true
    sleep 1
fi

# 3. Start Go REST API backend (with Delve or normal)
if [ "$ENABLE_DEBUG" = "1" ]; then
    # Ensure dlv is available
    DLV_BIN="$(command -v dlv || echo "$HOME/go/bin/dlv")"
    if [ ! -x "$DLV_BIN" ] && ! command -v dlv >/dev/null 2>&1; then
        echo "🔧  Installing Delve (dlv) debugger..."
        (cd "$ROOT/backend" && go install github.com/go-delve/delve/cmd/dlv@latest)
        DLV_BIN="$(command -v dlv || echo "$HOME/go/bin/dlv")"
    fi

    echo "🐛  Launching Go REST API under Delve Debugger (http://localhost:8080, dlv: :4000)..."
    (
        cd "$ROOT/backend"
        go build -gcflags="all=-N -l" -o server_debug ./cmd/server
        "$DLV_BIN" exec ./server_debug --headless --listen=:4000 --api-version=2 --accept-multiclient --continue
    ) &
    BACKEND_PID=$!
else
    echo "🚀  Starting Go REST API Backend (http://localhost:8080)..."
    (cd "$ROOT/backend" && go run ./cmd/server) &
    BACKEND_PID=$!
fi

# Ensure backend process terminates when exiting dev script (Ctrl+C)
cleanup() {
    echo ""
    echo "🛑  Stopping Go Backend & Debugger (PID: $BACKEND_PID)..."
    kill -9 $BACKEND_PID 2>/dev/null || true
    if [ "$ENABLE_DEBUG" = "1" ]; then
        lsof -ti:4000 | xargs kill -9 2>/dev/null || true
    fi
    exit 0
}
trap cleanup EXIT INT TERM

# 4. Start SvelteKit Vite dev server in current shell
echo "✨  Starting SvelteKit Frontend Dev Server (http://localhost:5173)..."
cd "$ROOT/frontend"
npm run dev
