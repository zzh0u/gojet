#!/usr/bin/env bash
# 通用 JSON 联调发送。
# 用法：./tests/send-api.sh METHOD /path [json-body]
set -euo pipefail

METHOD="${1:-}"
API_PATH="${2:-}"
BODY="${3:-}"

if [[ -z "$METHOD" || -z "$API_PATH" ]]; then
  echo "用法: $0 METHOD /path [json-body]" >&2
  exit 1
fi

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
CURL_ARGS=(-sS -X "$METHOD" "${BASE_URL}${API_PATH}")

if [[ -n "$BODY" ]]; then
  CURL_ARGS+=(-H "Content-Type: application/json" -d "$BODY")
fi
if [[ -n "${TOKEN:-}" ]]; then
  CURL_ARGS+=(-H "Authorization: Bearer ${TOKEN}")
fi

curl "${CURL_ARGS[@]}"
echo
