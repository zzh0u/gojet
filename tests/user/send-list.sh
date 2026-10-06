#!/usr/bin/env bash
# 本地联调：用户列表 GET /v1/user
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
exec "$DIR/../send-api.sh" GET /v1/user
