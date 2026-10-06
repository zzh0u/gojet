#!/usr/bin/env bash
# 本地联调：获取单个用户 GET /v1/user/:id
# 用法：./tests/user/send-get.sh [id]
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
ID="${1:-1}"
exec "$DIR/../send-api.sh" GET "/v1/user/${ID}"
