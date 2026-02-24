#!/bin/bash
#
# 重命名模板项目
# 用法: ./scripts/rename.sh <go-module> <serviceName>
# 示例: ./scripts/rename.sh github.com/zheng/order-service orderService
#

set -e

if [ $# -ne 2 ]; then
    echo "用法: $0 <go-module> <service-name>"
    echo "示例: $0 github.com/zheng/order-service orderService"
    echo ""
    echo "  go-module     新的 Go module 路径"
    echo "  service-name  新的服务名 (camelCase，如 orderService)"
    exit 1
fi

NEW_MODULE="$1"
NEW_SERVICE="$2"

OLD_MODULE="go-grpc-http-template"
OLD_SERVICE="myService"
OLD_SERVICE_LOWER="myservice"

# 服务名转小写
NEW_SERVICE_LOWER=$(echo "$NEW_SERVICE" | tr '[:upper:]' '[:lower:]')

echo "========================================="
echo " 重命名模板项目"
echo "========================================="
echo " Module:  $OLD_MODULE -> $NEW_MODULE"
echo " Service: $OLD_SERVICE -> $NEW_SERVICE"
echo "          $OLD_SERVICE_LOWER -> $NEW_SERVICE_LOWER"
echo "========================================="
echo ""

# 确认
read -p "确认重命名? (y/N) " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "已取消"
    exit 0
fi

# 1. 替换文件内容 (先替换 module，再替换 service 名)
echo "[1/4] 替换文件内容..."

# 需要替换的文件类型
FILES=$(find . -type f \( -name "*.go" -o -name "*.proto" -o -name "*.yaml" -o -name "*.toml" -o -name "*.mod" -o -name "makefile" \) \
    ! -path "./third_party/*" ! -path "./tmp/*" ! -path "./.git/*")

for f in $FILES; do
    if grep -q "$OLD_MODULE\|$OLD_SERVICE\|$OLD_SERVICE_LOWER" "$f" 2>/dev/null; then
        sed -i "s|$OLD_MODULE|$NEW_MODULE|g" "$f"
        sed -i "s|$OLD_SERVICE_LOWER|$NEW_SERVICE_LOWER|g" "$f"
        sed -i "s|$OLD_SERVICE|$NEW_SERVICE|g" "$f"
    fi
done

# 2. 重命名文件
echo "[2/4] 重命名文件..."

rename_file() {
    local old="$1"
    local new="$(echo "$old" | sed "s|$OLD_SERVICE|$NEW_SERVICE|g")"
    if [ "$old" != "$new" ]; then
        mv "$old" "$new"
        echo "  $old -> $new"
    fi
}

# 先重命名文件 (从深到浅)
find . -type f -name "*${OLD_SERVICE}*" ! -path "./.git/*" ! -path "./tmp/*" | sort -r | while read -r f; do
    rename_file "$f"
done

# 3. 重命名目录 (从深到浅)
echo "[3/4] 重命名目录..."

rename_dir() {
    local old="$1"
    local new="$(echo "$old" | sed "s|$OLD_SERVICE_LOWER|$NEW_SERVICE_LOWER|g" | sed "s|$OLD_SERVICE|$NEW_SERVICE|g")"
    if [ "$old" != "$new" ]; then
        mv "$old" "$new"
        echo "  $old -> $new"
    fi
}

find . -type d \( -name "*${OLD_SERVICE}*" -o -name "*${OLD_SERVICE_LOWER}*" \) ! -path "./.git/*" | sort -r | while read -r d; do
    rename_dir "$d"
done

# 4. 更新 go.sum
echo "[4/4] 更新依赖..."
go mod tidy 2>/dev/null || echo "  go mod tidy 需要手动执行"

echo ""
echo "========================================="
echo " 重命名完成!"
echo "========================================="
echo " 后续步骤:"
echo "   make desc    # 生成 proto 描述符"
echo "   make run     # 启动服务"
echo "========================================="
