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

echo "🚀 Starting DeepPhotos Build Process..."
> build.log

(cd frontend && npm install >> ../build.log 2>&1) &
spinner $! "Installing Frontend Dependencies"

(cd frontend && npm run build >> ../build.log 2>&1) &
spinner $! "Building Frontend (SvelteKit)"

(cd backend && go build -o server ./cmd/server >> ../build.log 2>&1) &
spinner $! "Compiling Backend (Go)"

echo "✅ DeepPhotos Build Complete! (See build.log for details)"
