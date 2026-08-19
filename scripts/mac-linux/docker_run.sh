#!/usr/bin/env bash
# DeepPhotos — Start / Stop Containers
# Usage: ./scripts/mac-linux/docker_run.sh         → start
#        ./scripts/mac-linux/docker_run.sh stop     → stop
#        ./scripts/mac-linux/docker_run.sh restart  → restart
set -e

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
COMPOSE="$ROOT/docker/docker-compose.yaml"
CMD="${1:-up}"

case "$CMD" in
  stop)
    echo "🛑  Stopping DeepPhotos..."
    docker compose -f "$COMPOSE" down
    echo "✅  Stopped."
    ;;
  restart)
    echo "🔄  Restarting DeepPhotos..."
    docker compose -f "$COMPOSE" down
    docker compose -f "$COMPOSE" up -d
    echo "✅  Restarted."
    ;;
  up|*)
    echo "🚀  Starting DeepPhotos..."
    docker compose -f "$COMPOSE" up -d
    echo ""
    echo "✅  Running at:"
    echo "    Frontend  →  http://localhost:5173"
    echo "    Backend   →  http://localhost:8080"
    echo "    MinIO     →  http://localhost:9001"
    echo ""
    echo "    Stop:     ./scripts/mac-linux/docker_run.sh stop"
    echo "    Restart:  ./scripts/mac-linux/docker_run.sh restart"
    ;;
esac
