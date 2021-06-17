# NATSRPC
NATSRPC 是一个基于nats的简单rpc

## Feature
* 使用简单，不需要服务发现
* 代码生成器生成client和server代码
* 支持空间隔离
* 支持定向发送也支持负载均衡(nats的同组内随机)
* 不用手动定义subject
* 支持单协程回调(适用于逻辑单协程模型)

## 使用
1. 引用包 `go get github.com/byebyebruce/natsrpc`
2. 编译代码生成器 `go get github.com/byebyebruce/natsrpc/cmd/natsrpc_codegen`
3. 定义服务接口[示例](testdata/greeter.go)

4. 生成代码
```shell
natsrpc_codegen -s="greeter.go"
```
5. 写服务实现[示例](example/example_greeter.go)
## 示例
* [Client](example/client/main.go)
* [Server](example/server/main.go)
> 运行示例需要部署gnatsd，如果没有可以临时启动`go run cmd/simple_natsserver/main.go`

## 压测工具
1. 广播 `go run bench/pub/main.go -server=nats://127.0.0.1:4222`

2. 请求 `go run bench/req/main.go -server=nats://127.0.0.1:4222`