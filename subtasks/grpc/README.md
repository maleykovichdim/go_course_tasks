# Task: Developing an RPC Service

## Overview
Practice developing an RPC service. You need to develop an RPC application client and server. You can use either the RPC technology implementation from the standard library or gRPC at your discretion.

---

## Task #1: Develop the Server Application

### Requirements
The server has two functions:

1. **Send()**: Allows the client to send an array of messages to the server. The messages are added to an array.
2. **Messages()**: Returns all accumulated messages to the client. Each message must contain:
   - An identifier
   - Sending time
   - Text

---

## Task #2: Develop the RPC Service Client Application

### Requirements
The client application, when launched, must do the following:
1. Send several messages to the server.
2. Request all messages accumulated on the server.
3. Display those messages on the screen.

---

## some gRPC Setup Instructions

   Before writting code:
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

 -- in power shell:
$env:PATH += ";$(go env GOPATH)\bin"

 -- from grpc folder execute this to build prito files:
protoc --go_out=./grpc_messages --go_opt=paths=source_relative   --go-grpc_out=./grpc_messages --go-grpc_opt=paths=source_relative   proto/grpc_messages.proto
















