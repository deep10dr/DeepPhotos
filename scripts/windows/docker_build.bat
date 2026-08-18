@echo off
REM DeepPhotos — Build Docker Images
REM Usage: scripts\windows\docker_build.bat

cd /d "%~dp0..\.."
set LOG=%cd%\docker_build.log

echo [*] Building DeepPhotos Docker Images...
echo     Log: %LOG%
echo.

docker compose -f docker\docker-compose.yaml build --no-cache 2>&1 | tee "%LOG%"

echo.
echo [OK] Build complete!
echo      Run: scripts\windows\docker_run.bat
