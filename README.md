# go-grpc-http-template

基于 [go-zero](https://go-zero.dev/) 的 gRPC + HTTP Gateway 生产级模板项目。

## 特性

- gRPC 服务 + HTTP Gateway 同进程运行
- `google.api.http` 注解自动生成 HTTP 路由
- HTTP 响应统一包装为 `{code, msg, data}` 格式
- 生产级配置：超时、限流、自适应降载、日志、健康检查
- 可选：Prometheus 指标、OpenTelemetry 链路追踪、Etcd 服务发现
- `protoc-gen-doc` 自动生成接口文档
- `air` 热重载开发模式

## 快速开始

### 前置依赖

- Go 1.21+
- protoc
- [goctl](https://go-zero.dev/docs/tasks/installation/goctl)

### 初始化项目

```bash
# 克隆模板
git clone <repo-url> my-project && cd my-project

# 重命名为你的项目 (module + 服务名)
./scripts/rename.sh github.com/yourname/my-project myService

# 安装工具 (protoc-gen-doc, air)
make install-tools
```

### 开发

```bash
make run     # 启动服务 (gRPC :8080 + HTTP :8081)
make air     # 热重载模式，修改代码自动重启
```

### 测试接口

```bash
# HTTP
curl http://localhost:8081/api/v1/health
# 返回: {"code":200,"msg":"ok","data":{"status":"ok"}}

# gRPC (需要 dev 模式开启 reflection)
grpcurl -plaintext localhost:8080 myService.myService/Health
```

### 代码生成

```bash
make gen     # 从 proto 生成 Go 代码
make desc    # 生成 proto 描述符 (gateway 需要)
make doc     # 生成接口文档到 doc/
make all     # 以上全部
```

## 项目结构

```
├── myService.proto          # proto 定义 (google.api.http 注解)
├── myService.go             # 入口，启动 gRPC + HTTP Gateway
├── etc/
│   └── myService.yaml       # 配置文件
├── internal/
│   ├── config/              # 配置结构体
│   ├── logic/myservice/     # 业务逻辑
│   ├── middleware/          # HTTP 中间件 (响应包装)
│   ├── server/myservice/    # gRPC server 实现
│   └── svc/                 # 服务上下文 (依赖注入)
├── myService/               # proto 生成的 Go 代码
├── client/myservice/        # gRPC 客户端
├── third_party/google/      # google/api proto 依赖
├── scripts/
│   └── rename.sh            # 项目重命名脚本
├── makefile
└── .air.toml                # air 热重载配置
```

## 添加新接口

1. 在 `myService.proto` 中定义 rpc 方法和 `google.api.http` 注解：

```protobuf
rpc GetUser(GetUserRequest) returns(GetUserResponse) {
    option (google.api.http) = {
        get: "/api/v1/users/{id}"
    };
}
```

2. 生成代码并实现逻辑：

```bash
make gen                          # 生成代码框架
# 编辑 internal/logic/myservice/  # 实现业务逻辑
make desc                         # 更新描述符
```

## 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `ListenOn` | gRPC 监听地址 | `0.0.0.0:8080` |
| `Gateway.Port` | HTTP 网关端口 | `8081` |
| `Mode` | 运行模式 (`dev`/`test`/`pro`) | `pro` |
| `Timeout` | gRPC 超时 (ms) | `3000` |
| `CpuThreshold` | CPU 降载阈值 (0-1000) | `900` |
| `Health` | gRPC 健康检查 | `true` |
| `Etcd` | 服务发现 (可选) | 关闭 |
| `Prometheus` | 指标采集 (可选) | 关闭 |
| `Telemetry` | OpenTelemetry 追踪 (可选) | 关闭 |
