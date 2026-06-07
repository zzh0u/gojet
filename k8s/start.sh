#!/bin/bash

# GoJet Kubernetes 一键启动脚本
# 该脚本自动完成镜像构建和 Kubernetes 部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的信息
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        error "$1 未安装，请先安装 $1"
        exit 1
    fi
}

# 检查必要工具
info "检查必要工具..."
check_command kubectl
check_command docker

# 检查 Minikube 是否运行（可选）
if command -v minikube &> /dev/null; then
    if minikube status &> /dev/null; then
        info "检测到 Minikube 正在运行，配置 Docker 环境..."
        eval $(minikube docker-env)
        USING_MINIKUBE=true
    else
        warn "Minikube 未运行，将使用本地 Docker 环境"
        warn "如需使用 Minikube，请先运行: minikube start"
        USING_MINIKUBE=false
    fi
else
    warn "未检测到 Minikube，将使用本地 Docker 环境"
    warn "确保 Kubernetes 集群可以访问本地镜像仓库"
    USING_MINIKUBE=false
fi

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# 构建应用镜像
info "构建 GoJet 应用镜像..."
docker build -t gojet:latest .
success "镜像构建完成"

# 应用 Kubernetes 清单
info "部署 Kubernetes 资源..."
kubectl apply -k k8s/
success "Kubernetes 资源部署完成"

# 等待部署完成
info "等待部署就绪..."
kubectl wait --for=condition=available --timeout=120s deployment/gojet -n gojet 2>/dev/null || true
kubectl wait --for=condition=available --timeout=120s deployment/postgres -n gojet 2>/dev/null || true
kubectl wait --for=condition=available --timeout=120s deployment/redis -n gojet 2>/dev/null || true

# 显示部署状态
echo ""
success "GoJet 应用部署完成！"
echo ""
info "部署状态:"
kubectl get pods,svc -n gojet

echo ""
info "常用命令:"
echo "  查看日志:           kubectl logs -f deployment/gojet -n gojet"
echo "  查看 Pod 详情:      kubectl describe pod -l app.kubernetes.io/name=gojet -n gojet"
echo "  进入应用容器:       kubectl exec -it deployment/gojet -n gojet -- /bin/sh"
echo "  查看事件:           kubectl get events -n gojet --sort-by=.lastTimestamp"

# 如果是 Minikube，提供服务访问方式
if [ "$USING_MINIKUBE" = true ]; then
    echo ""
    info "Minikube 服务访问方式:"
    echo "  命令行打开:         minikube service gojet -n gojet"
    echo "  获取服务 URL:       minikube service gojet -n gojet --url"
else
    echo ""
    info "服务访问方式:"
    echo "  获取 NodePort:      kubectl get svc gojet -n gojet -o jsonpath='{.spec.ports[0].nodePort}'"
    echo "  或使用端口转发:     kubectl port-forward svc/gojet -n gojet 8080:8080"
fi

echo ""
info "清理命令:"
echo "  删除部署:           kubectl delete -k k8s/"
echo "  删除 PVC(数据):     kubectl delete pvc -n gojet --all"
