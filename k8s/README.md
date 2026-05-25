# gojet Kubernetes 部署

这些清单把项目拆成 3 个组件：`gojet` 应用、`postgres` 数据库、`redis` 缓存。

## 部署步骤

```bash
# 1. 让当前 shell 使用 Minikube 内置 Docker，这样镜像不用推送到仓库
eval $(minikube docker-env)

# 2. 构建应用镜像，名称需要和 k8s/app.yaml 中保持一致
docker build -t gojet:latest .

# 3. 应用 Kubernetes 清单
kubectl apply -k k8s

# 4. 查看状态
kubectl get pods,svc,pvc -n gojet

# 5. 打开应用服务
minikube service gojet -n gojet
```

## 常用命令

```bash
kubectl logs -f deployment/gojet -n gojet
kubectl describe pod -l app.kubernetes.io/name=gojet -n gojet
kubectl get events -n gojet --sort-by=.lastTimestamp
```

## 清理

```bash
kubectl delete -k k8s
```

如果要连同数据库和 Redis 数据一起删除，需要额外删除 PVC：

```bash
kubectl delete pvc -n gojet --all
```

