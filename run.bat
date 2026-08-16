@echo off
setlocal

:: Kill any existing server process running on port 8080
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080') do (
    if "%%a" neq "0" (
        echo 🔄 Stopping existing process running on port 8080 (PID: %%a)...
        taskkill /F /PID %%a >nul 2>&1
    )
)
timeout /t 1 /nobreak >nul

:: Check if Docker is running
docker info >nul 2>&1
if %errorlevel% equ 0 (
    echo 🚀 Starting DeepPhotos Containers (Docker Compose)...
    docker compose -f docker\docker-compose.yaml up -d
    echo ✨ Services started in Docker:
    echo    • Go REST API: http://localhost:8080
    echo    • MinIO Storage API: http://localhost:9000
    echo    • MinIO Web Console: http://localhost:9001
) else (
    echo ℹ️ Docker daemon is not currently running.
    echo 🚀 Starting DeepPhotos Go REST API Server (Standalone Mode)...
    if not exist backend\server.exe (
        echo 🔨 Compiling backend binary...
        cd backend
        go build -o server.exe ./cmd/server
        cd ..
    )
    echo ✨ Go REST API active at http://localhost:8080 (Press Ctrl+C to stop)
    cd backend
    server.exe
)
