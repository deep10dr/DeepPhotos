#!/usr/bin/env bash
set -e

# Kill any existing server process running on port 8080
if lsof -pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "🔄 Stopping existing process running on port 8080..."
    lsof -ti:8080 | xargs kill -9 >/dev/null 2>&1 || true
    sleep 1
fi

if docker info >/dev/null 2>&1; then
    echo "🚀 Starting DeepPhotos Containers (Docker Compose)..."
    docker compose -f docker/docker-compose.yaml up -d
    echo "✨ Services started in Docker:"
    echo "   • Go REST API: http://localhost:8080"
    echo "   • MinIO Storage API: http://localhost:9000"
    echo "   • MinIO Web Console: http://localhost:9001"
else
    echo "ℹ️ Docker daemon is not currently running."
    echo "🚀 Starting DeepPhotos Go REST API Server (Standalone Mode)..."
    if [ ! -f backend/server ]; then
        echo "🔨 Compiling backend binary..."
        (cd backend && go build -o server ./cmd/server)
    fi
    echo "✨ Go REST API active at http://localhost:8080 (Press Ctrl+C to stop)"
    cd backend && ./server
fi
