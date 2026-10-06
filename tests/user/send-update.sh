#!/usr/bin/env bash
# 本地联调：更新用户 PUT /v1/user/:id
# 用法：./tests/user/send-update.sh [id]
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
ID="${1:-1}"
exec "$DIR/../send-api.sh" PUT "/v1/user/${ID}" '{
  "username": "alice-updated",
  "nick_name": "alice Updated",
  "email": "updated@example.com"
}'
