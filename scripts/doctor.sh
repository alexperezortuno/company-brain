#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail=0
for command in docker curl; do
  if ! command -v "$command" >/dev/null 2>&1; then
    echo "ERROR: missing command: $command"
    fail=1
  else
    echo "OK: $command"
  fi
done

if command -v docker >/dev/null 2>&1; then
  if docker compose version >/dev/null 2>&1; then
    echo "OK: docker compose"
  else
    echo "ERROR: Docker Compose plugin is unavailable"
    fail=1
  fi
fi

[[ -f .env ]] && echo "OK: .env" || { echo "ERROR: .env missing; run make bootstrap"; fail=1; }
[[ -f config/brain.yaml ]] && echo "OK: config/brain.yaml" || { echo "ERROR: config/brain.yaml missing"; fail=1; }
exit "$fail"
