#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/../.."

echo "⚡ Starting DeepPhotos Development Environment..."

# 1. Start backend in background if desired
(cd backend && go run ./cmd/server) &
BACKEND_PID=$!

# 2. Start frontend Vite dev server
cd frontend
npm run dev

# Cleanup on exit
trap "kill $BACKEND_PID" EXIT
