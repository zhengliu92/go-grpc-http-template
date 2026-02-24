.PHONY: gen desc doc all install-tools run air

gen:
	goctl rpc protoc \
	myService.proto \
	-style goZero \
	--go_out=. \
	--go-grpc_out=. \
	--zrpc_out=. \
	--proto_path=. \
	--proto_path=third_party \
	-m

# 生成 proto 描述符文件 (gateway 需要，包含 google.api.http 注解)
desc:
	protoc --proto_path=. --proto_path=third_party \
	--descriptor_set_out=etc/myService.pb \
	--include_imports \
	myService.proto

# 生成接口文档到 doc/ 文件夹
doc:
	@mkdir -p doc
	protoc --proto_path=. --proto_path=third_party --doc_out=doc --doc_opt=markdown,api.md myService.proto
	protoc --proto_path=. --proto_path=third_party --doc_out=doc --doc_opt=html,api.html myService.proto

# 一键生成所有 (代码 + 描述符 + 文档)
all: gen desc doc

# 启动服务
run: desc
	go run myService.go -f etc/myService.yaml

# 热重载开发模式 (修改代码自动重启)
air: desc
	air

# 安装文档生成工具
install-tools:
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
	go install github.com/air-verse/air@latest
