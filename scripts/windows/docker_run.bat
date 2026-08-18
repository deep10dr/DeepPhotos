@echo off
REM DeepPhotos — Start / Stop Containers
REM Usage: scripts\windows\docker_run.bat          → start
REM        scripts\windows\docker_run.bat stop     → stop
REM        scripts\windows\docker_run.bat restart  → restart

cd /d "%~dp0..\.."
set COMPOSE=docker\docker-compose.yaml
set CMD=%1

if /I "%CMD%"=="stop" (
    echo [*] Stopping DeepPhotos...
    docker compose -f %COMPOSE% down
    echo [OK] Stopped.
    goto :end
)

if /I "%CMD%"=="restart" (
    echo [*] Restarting DeepPhotos...
    docker compose -f %COMPOSE% down
    docker compose -f %COMPOSE% up -d
    echo [OK] Restarted.
    goto :end
)

REM Default: start
echo [*] Starting DeepPhotos...
docker compose -f %COMPOSE% up -d
echo.
echo [OK] Running at:
echo      Frontend  -^>  http://localhost:5173
echo      Backend   -^>  http://localhost:8080
echo      MinIO     -^>  http://localhost:9001
echo.
echo      Stop:     scripts\windows\docker_run.bat stop
echo      Restart:  scripts\windows\docker_run.bat restart

:end
