#!/bin/bash
#
# 查看 etcd 服务注册信息
# 用法: ./scripts/etcd-info.sh [etcd-endpoint]
# 示例: ./scripts/etcd-info.sh 127.0.0.1:2379
#

ETCD_ENDPOINT="${1:-127.0.0.1:2379}"
SERVICE_KEY="myservice.rpc"

# 检查 etcdctl 是否安装
if ! command -v etcdctl &>/dev/null; then
    echo "错误: 未找到 etcdctl，请先安装:"
    echo "  brew install etcd          # macOS"
    echo "  apt install etcd-client    # Ubuntu/Debian"
    exit 1
fi

export ETCDCTL_API=3

echo "========================================="
echo " etcd 服务信息"
echo " Endpoint: $ETCD_ENDPOINT"
echo "========================================="

# 1. 集群健康状态
echo ""
echo "[集群状态]"
etcdctl --endpoints="$ETCD_ENDPOINT" endpoint health 2>&1

# 2. 集群成员
echo ""
echo "[集群成员]"
etcdctl --endpoints="$ETCD_ENDPOINT" member list -w table 2>&1

# 3. 服务注册信息 (go-zero 注册的 key)
echo ""
echo "[服务注册: $SERVICE_KEY]"
KEYS=$(etcdctl --endpoints="$ETCD_ENDPOINT" get "$SERVICE_KEY" --prefix --keys-only 2>&1)
if [ -z "$KEYS" ]; then
    echo "  (无注册实例，服务未启动或 key 不匹配)"
else
    etcdctl --endpoints="$ETCD_ENDPOINT" get "$SERVICE_KEY" --prefix 2>&1
fi

# 4. 所有 key 列表 (便于排查 key 名称)
echo ""
echo "[全部 key 列表]"
ALL_KEYS=$(etcdctl --endpoints="$ETCD_ENDPOINT" get "" --prefix --keys-only 2>&1)
if [ -z "$ALL_KEYS" ]; then
    echo "  (etcd 中暂无数据)"
else
    echo "$ALL_KEYS"
fi

# 5. 监听服务变化 (可选)
echo ""
echo "========================================="
echo " 提示: 实时监听服务变化:"
echo "   etcdctl --endpoints=$ETCD_ENDPOINT watch $SERVICE_KEY --prefix"
echo "========================================="
