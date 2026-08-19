#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/../.."

scripts/mac-linux/build.sh
scripts/mac-linux/run.sh
