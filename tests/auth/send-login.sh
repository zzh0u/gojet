#!/usr/bin/env bash
# 本地联调：登录 POST /v1/auth/login
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$DIR/../send-api.sh" POST /v1/auth/login '{
  "username": "包子",
  "password": "123456"
}'
