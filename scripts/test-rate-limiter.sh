#!/usr/bin/env bash

set -euo pipefail

if ! command -v vegeta >/dev/null 2>&1; then
  echo "vegeta is not installed"
  echo "install: go install github.com/tsenart/vegeta@latest"
  exit 1
fi

TARGET_URL="${1:-http://localhost:8081/last-comic}"
RATE="${RATE:-200}"
DURATION="${DURATION:-10s}"
CONNECTIONS="${CONNECTIONS:-1}"
WORKERS="${WORKERS:-1}"

TMP_RESULTS="$(mktemp)"
trap 'rm -f "$TMP_RESULTS"' EXIT

echo "Testing rate limiter against: $TARGET_URL"
echo "rate=${RATE}/s duration=${DURATION} connections=${CONNECTIONS} workers=${WORKERS}"
echo

printf "GET %s\n" "$TARGET_URL" \
  | vegeta attack \
      -rate="${RATE}/s" \
      -duration="$DURATION" \
      -connections="$CONNECTIONS" \
      -workers="$WORKERS" \
      -max-workers="$WORKERS" \
  > "$TMP_RESULTS"

vegeta report "$TMP_RESULTS"
echo
echo "Status codes:"
vegeta report -type=text "$TMP_RESULTS" | grep -E "Status Codes|\\[[0-9]{3}\\]" || true
