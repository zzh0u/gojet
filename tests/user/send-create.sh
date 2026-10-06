#!/usr/bin/env bash
# 本地联调：创建用户 POST /v1/user
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$DIR/../send-api.sh" POST /v1/user '{
  "username": "alice",
  "nick_name": "Alice",
  "password": "123456",
  "email": "alice@example.com"
}'
