#!/usr/bin/env bash
# 本地联调：删除用户 DELETE /v1/user/:id
# 用法：./tests/user/send-delete.sh [id]
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
ID="${1:-1}"
exec "$DIR/../send-api.sh" DELETE "/v1/user/${ID}"
