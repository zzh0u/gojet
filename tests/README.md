# 本地联调

用于检查健康检查、登录和用户接口。这些脚本不参与编译。

除健康检查和登录外，请求需要 `TOKEN`。

## 发送示例

```bash
go run ./cmd/api

./tests/health/send-health.sh
./tests/auth/send-login.sh

export TOKEN=""
./tests/user/send-list.sh
./tests/user/send-get.sh 1
./tests/user/send-create.sh
./tests/user/send-update.sh 1
./tests/user/send-delete.sh 1
```

默认请求 `http://127.0.0.1:8080`。换地址时设置 `BASE_URL`。
