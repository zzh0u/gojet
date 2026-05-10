# API 测试命令清单

直接复制下面的命令到终端执行即可。

## 预设环境变量

```bash
export BASE_URL=http://localhost:8080
export TOKEN=
```

## 健康检查

```bash
curl -X GET "${BASE_URL}/v1/health"
```

## 用户登录

```bash
curl -X POST "${BASE_URL}/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "包子",
    "password": "123456"
  }'
```

## 用户列表

```bash
curl -X GET "${BASE_URL}/v1/user" \
  -H "Authorization: Bearer ${TOKEN}"
```

## 获取单个用户

```bash
curl -X GET "${BASE_URL}/v1/user/1" \
  -H "Authorization: Bearer ${TOKEN}"
```

## 创建用户

```bash
curl -X POST "${BASE_URL}/v1/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "username": "alice",
    "nick_name": "Alice",
    "password": "123456",
    "email": "alice@example.com"
  }'
```

## 更新用户

```bash
curl -X PUT "${BASE_URL}/v1/user/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "username": "alice-updated",
    "nick_name": "alice Updated",
    "email": "updated@example.com"
  }'
```

## 删除用户

```bash
curl -X DELETE "${BASE_URL}/v1/user/1" \
  -H "Authorization: Bearer ${TOKEN}"
```
