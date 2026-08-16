#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/../.."

echo "🔨 Building DeepPhotos Frontend..."
cd frontend
npm install
npm run build
cd ..

echo "🔨 Building DeepPhotos Go Backend..."
cd backend
go build -o server ./cmd/server
cd ..

echo "✅ DeepPhotos Build Complete!"
