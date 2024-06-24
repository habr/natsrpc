```
  _   _       _______ _____   _____  _____   _____ 
 | \ | |   /\|__   __/ ____| |  __ \|  __ \ / ____|
 |  \| |  /  \  | | | (___   | |__) | |__) | |     
 | . ` | / /\ \ | |  \___ \  |  _  /|  ___/| |     
 | |\  |/ ____ \| |  ____) | | | \ \| |    | |____ 
 |_| \_/_/    \_\_| |_____/  |_|  \_\_|     \_____|
```

## What is NATSRPC
> NATSRPC is based on [NATS](https://nats.io/) as a message communication, use [gRPC](https://www.grpc.io/) way to define the RPC framework of the interface

<p align="center">
  <span>English</span> |
  <a href="README.cn.md#readme">中文</a>
</p>

![GitHub release (with filter)](https://img.shields.io/github/v/release/byebyebruce/natsrpc)
![](https://hits.sh/github.com/LeKovr/natsrpc/doc/hits.svg?label=visit)

## Motivation  
NATS needs to manually define cumbersome and error-prone code such as subject, request, reply, handler, etc. to send and receive messages.
gRPC needs to be connected to the server endpoint to send the request.
The purpose of NATRPC is to define the interface like gRPC. Like NATS, it does not care about the specific network location. It only needs to listen and send to complete the RPC call.

## Feature
* Use the gRPC interface to define the method, which is simple to use and generates code with one click
* Support spatial isolation, you can also specify the id to send
* Multiple services can be load balanced (random in the same group of nats)
* Support Header and return error
* Support single coroutine and multi-coroutine handle
* Support middleware
* Support delayed reply to messages
* Support custom encoder

## How It Works
The upper layer pairs nats through Server, Service, and Client.Conn and Subscription are encapsulated.
The underlying layer transmits messages through nats request and publish.A service will create a subscription with the service name as the subject, and if there is a publish method, a sub will be created to receive the publish.
When the client sends a request, the subject will be the name of the service, and the header of nats msg will pass the method name.
After the service receives the message, it takes out the method name, and then calls the corresponding handler. The result returned by the handler will be returned to the client through the reply subject of nats msg.

## Install Tools
1. protoc(v3.17.3) 
   - [Linux](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip)
   - [MacOS](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-osx-x86_64.zip)
   - [Windows](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip)
   
2. protoc-gen-gogo 
   ```shell
   go install github.com/gogo/protobuf/protoc-gen-gogo@v1.3.2
   ```
3. protoc-gen-natsrpc 
   ```shell
   go install github.com/LeKovr/natsrpc/cmd/protoc-gen-natsrpc@v0.7.0
   ```

## Quick Start
* [nats-server](https://github.com/nats-io/nats-server/releases)>=2.2.0
1. Reference package
   ```shell
   go get github.com/LeKovr/natsrpc
   ```
2. Define the service interface example.proto
    ```
    syntax = "proto3";

    package example;
    option go_package = "github.com/LeKovr/natsrpc/example;example";

    message HelloRequest {
      string name = 1;
    }

    message HelloReply {
      string message = 1;
    }

    service Greeter {
      rpc Hello (HelloRequest) returns (HelloReply) {}
    }
    ```
   
3. Generate client and server code
    ```shell
    protoc --proto_path=. \
    --gogo_out=paths=source_relative:. \
    --natsrpc_out=paths=source_relative:. \
    *.proto
    ```
4. Server implements the interface at the end and create a service
   ```
   type HelloSvc struct {
   }

   func (s *HelloSvc) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
       return &example.HelloReply{
           Message: "hello " + req.Name,
       }, nil
   }

   func main() {
       conn, err := nats.Connect(*nats_url)
       defer conn.Close()

       server, err := natsrpc.NewServer(conn)
       defer server.Close(context.Background())

       svc, err := example.RegisterGreetingNRServer(server, &HelloSvc{})
       defer svc.Close()
       
       select{
       }
   }

   ```
   
5. Client calls rpc
   ```
   client:=natsrpc.NewClient(conn)
   
   cli := example.NewGreeterNRClient(client)
   rsp,err:=cli.Hello(context.Background(), &example.HelloRequest{Name: "natsrpc"})
   ```
 
## Examples
[here](./example)

## Bench Tool
1. Request `go run ./example/tool/request_bench -url=nats://127.0.0.1:4222`
2. Broadcast `go run ./example/tool/publish_bench -url=nats://127.0.0.1:4222`

## TODO
-[x] The service definition file is changed to the gRPC standard
-[x] Support return error
-[x] Support Header
-[x] Generate Client interface
-[x] Support middleware
-[x] Default multithreading, support a single thread at the same time
-[] Support goroutine pool
-[] Support byte pool
