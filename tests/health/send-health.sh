#!/usr/bin/env bash
# 本地联调：健康检查 GET /v1/health
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$DIR/../send-api.sh" GET /v1/health
