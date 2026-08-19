@echo off
setlocal

echo ⚡ Starting DeepPhotos Development Environment...

:: Start backend in a new command prompt window
start "DeepPhotos Backend" cmd /c "cd backend && go run ./cmd/server"

:: Start frontend Vite dev server in the current window
cd frontend
call npm run dev
