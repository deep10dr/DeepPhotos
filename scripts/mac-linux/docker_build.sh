#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/../.."

spinner() {
    local pid=$1
    local msg=$2
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    while kill -0 $pid 2>/dev/null; do
        for i in $(seq 0 9); do
            printf "\r\033[36m${spinstr:$i:1}\033[0m %s..." "$msg"
            sleep 0.1
        done
    done
    wait $pid
    local status=$?
    if [ $status -eq 0 ]; then
        printf "\r\033[32m✔\033[0m %s... Done! \n" "$msg"
    else
        printf "\r\033[31m✖\033[0m %s... Failed! \n" "$msg"
        exit $status
    fi
}

echo "🐳 Building DeepPhotos Docker Containers..."
> docker_build.log

(docker compose -f docker/docker-compose.yaml build >> docker_build.log 2>&1) &
spinner $! "Building Docker Images (This may take a while)"

echo "✅ Docker Build Complete! (See docker_build.log for details)"
