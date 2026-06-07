# gojet Kubernetes 部署

本目录包含将 gojet 应用部署到 Kubernetes 集群的完整配置清单。采用 Kustomize 作为配置管理工具，实现声明式、可复用的部署方案。

## 架构概览

gojet 在 Kubernetes 中以微服务形态运行，包含三个核心组件：

| 组件 | 类型 | 功能说明 |
|------|------|----------|
| **gojet** | Deployment | Go Web 应用主服务，处理 HTTP 请求和业务逻辑 |
| **postgres** | StatefulSet | PostgreSQL 数据库，提供持久化数据存储 |
| **redis** | Deployment | Redis 缓存服务，用于会话和高速数据存取 |

### 资源组织

```
k8s/
├── namespace.yaml      # 创建独立的命名空间隔离资源
├── configmap.yaml      # 应用配置（非敏感信息）
├── secret.yaml         # 密钥配置（密码、JWT 密钥等）
├── postgres.yaml       # PostgreSQL StatefulSet、Service、PVC
├── redis.yaml          # Redis Deployment、Service
├── app.yaml            # gojet 应用 Deployment、Service
├── kustomization.yaml  # Kustomize 配置入口
└── start.sh            # 一键启动脚本
```

### 设计特点

- **命名空间隔离**：所有资源部署在 `gojet` 命名空间下，避免与其他应用冲突
- **配置与代码分离**：使用 ConfigMap 和 Secret 管理配置，无需重新构建镜像即可调整参数
- **数据持久化**：PostgreSQL 使用 PVC 存储数据，确保容器重启后数据不丢失
- **服务发现**：通过 Kubernetes Service 实现组件间通信，无需硬编码 IP 地址
- **健康检查**：应用配置了存活探针和就绪探针，确保服务稳定性

## 快速开始

### 一键启动（推荐）

```bash
./k8s/start.sh
```

该脚本会自动完成：环境检查、镜像构建、资源部署和状态展示。

### 手动部署

如需了解详细步骤或自定义部署，可执行以下命令：

```bash
# 1. 配置 Docker 环境（Minikube 用户需要执行）
eval $(minikube docker-env)

# 2. 构建应用镜像
docker build -t gojet:latest .

# 3. 应用 Kubernetes 清单
kubectl apply -k k8s/

# 4. 查看部署状态
kubectl get pods,svc -n gojet

# 5. 访问应用（Minikube）
minikube service gojet -n gojet
```

## 访问服务

### Minikube 环境

```bash
# 自动打开浏览器
minikube service gojet -n gojet

# 仅获取访问 URL
minikube service gojet -n gojet --url
```

### 标准 Kubernetes 集群

应用通过 NodePort 服务暴露，默认端口为 30080：

```bash
# 查看 NodePort
kubectl get svc gojet -n gojet -o jsonpath='{.spec.ports[0].nodePort}'

# 或使用端口转发直接访问
kubectl port-forward svc/gojet -n gojet 8080:8080
# 然后访问 http://localhost:8080
```

## 运维命令

### 查看日志

```bash
# 查看应用日志
kubectl logs -f deployment/gojet -n gojet

# 查看数据库日志
kubectl logs -f deployment/postgres -n gojet

# 查看 Redis 日志
kubectl logs -f deployment/redis -n gojet
```

### 调试诊断

```bash
# 查看 Pod 详细状态
kubectl describe pod -l app.kubernetes.io/name=gojet -n gojet

# 查看集群事件
kubectl get events -n gojet --sort-by=.lastTimestamp

# 进入应用容器
kubectl exec -it deployment/gojet -n gojet -- /bin/sh

# 进入数据库容器
kubectl exec -it deployment/postgres -n gojet -- psql -U gojet -d gojet
```

### 扩缩容

```bash
# 扩展应用副本数
kubectl scale deployment/gojet --replicas=3 -n gojet

# 查看扩缩容状态
kubectl get pods -n gojet -w
```

## 清理资源

```bash
# 删除所有 Kubernetes 资源（保留数据）
kubectl delete -k k8s/

# 删除 PVC（会清除数据库和 Redis 数据）
kubectl delete pvc -n gojet --all
```

## 自定义配置

### 修改应用配置

编辑 `configmap.yaml` 调整非敏感配置：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gojet-config
data:
  APP_MODE: "release"
  LOG_LEVEL: "info"
  # 修改其他配置项...
```

### 修改密钥

编辑 `secret.yaml` 更新敏感信息（注意：Secret 需要 base64 编码）：

```bash
# 生成 base64 编码的字符串
echo -n 'your-password' | base64
```

### 使用 Kustomize 覆盖

创建 `k8s/overlays/production/` 目录用于生产环境定制：

```yaml
# k8s/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patchesStrategicMerge:
  - replica-patch.yaml
```

部署时指定 overlay：

```bash
kubectl apply -k k8s/overlays/production/
```

## 兼容性说明

- **Kubernetes**: 1.20+
- **kubectl**: 1.20+
- **Minikube**: 1.25+（本地开发）

## 故障排查

### Pod 无法启动

```bash
# 查看 Pod 事件和错误信息
kubectl describe pod -n gojet <pod-name>

# 查看容器启动日志
kubectl logs -n gojet <pod-name> --previous
```

### 镜像拉取失败

Minikube 环境需要确保镜像已构建到 Minikube 的 Docker 中：

```bash
eval $(minikube docker-env)
docker build -t gojet:latest .
```

### 服务无法访问

```bash
# 检查 Service 配置
kubectl get svc gojet -n gojet

# 检查 Endpoints 是否正常
kubectl get endpoints gojet -n gojet

# 检查 Pod 标签是否匹配
kubectl get pods -n gojet --show-labels
```
