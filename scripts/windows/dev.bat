@echo off
REM DeepPhotos — Local Development Environment (Windows)
REM Usage: scripts\windows\dev.bat          → Normal dev mode
REM        scripts\windows\dev.bat --debug  → Dev mode with Delve (dlv) debugger on port 4000

setlocal
cd /d "%~dp0..\.."

set ENABLE_DEBUG=0
if "%1"=="--debug" set ENABLE_DEBUG=1
if "%1"=="-d" set ENABLE_DEBUG=1

echo ⚡  Starting DeepPhotos Local Development Environment...
if "%ENABLE_DEBUG%"=="1" (
    echo 🐛  Delve Debugging Mode ACTIVE ^(Port 4000^)
)
echo.

:: 1. Start MinIO Object Storage container via Docker Compose
docker info >nul 2>&1
if %errorlevel% equ 0 (
    echo [*] Starting MinIO Object Storage container...
    docker compose -f docker\docker-compose.yaml up -d minio
    echo [OK] MinIO active at http://localhost:9000 ^(Web Console: http://localhost:9001^)
) else (
    echo [i] Docker is not active - backend will use local disk storage fallback.
)

echo.

:: 2. Start Go Backend in a new window
if "%ENABLE_DEBUG%"=="1" (
    echo [*] Starting Go REST API Backend with Delve Debugger on port 4000...
    start "DeepPhotos Backend (Delve Debugger :4000)" cmd /c "cd backend && go build -gcflags=\"all=-N -l\" -o server_debug.exe ./cmd/server && dlv exec ./server_debug.exe --headless --listen=:4000 --api-version=2 --accept-multiclient --continue"
) else (
    echo [*] Starting Go REST API Backend ^(http://localhost:8080^)...
    start "DeepPhotos Backend" cmd /c "cd backend && go run ./cmd/server"
)

:: 3. Start SvelteKit Frontend in current window
echo [*] Starting SvelteKit Frontend Dev Server ^(http://localhost:5173^)...
cd frontend
call npm run dev
