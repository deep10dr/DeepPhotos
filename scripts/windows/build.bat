@echo off
setlocal

echo 🔨 Building DeepPhotos Frontend...
cd frontend
call npm install
call npm run build
cd ..

echo 🔨 Building DeepPhotos Go Backend...
cd backend
go build -o server.exe ./cmd/server
cd ..

echo ✅ DeepPhotos Build Complete!
