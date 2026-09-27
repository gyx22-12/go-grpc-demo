# go-grpc-demo

[![CI](https://github.com/gyx22-12/go-grpc-demo/actions/workflows/ci.yml/badge.svg)](https://github.com/gyx22-12/go-grpc-demo/actions/workflows/ci.yml)

一个最小可运行的 Go gRPC 示例，配套 GitHub Actions CI，用于学习 gRPC（protobuf）与 CI/CD。

## 功能

- `Greeter.SayHello`：一元调用（unary）
- `Greeter.SayHelloStream`：服务端流式（server streaming）

## 目录结构

- `pb/greeter.proto`：protobuf 服务定义
- `pb/*.pb.go`：protoc 生成的 Go 代码（已提交，CI 无需重新生成）
- `server/`：gRPC 服务端（含 bufconn 内存单测）
- `client/`：gRPC 客户端
- `Dockerfile`：多阶段构建镜像
- `.github/workflows/ci.yml`：CI（gofmt / go vet / go build / go test / docker build）

## 运行

```bash
# 服务端
go run ./server
# 客户端（另开终端）
go run ./client
```

## 重新生成 protobuf 代码

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pb/greeter.proto
```
