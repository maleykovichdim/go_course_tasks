/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//GRPC CLIENT

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "my_grpc/grpc_messages/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newMessage := []pb.Message{
		{
			ID:   1,
			Text: "Hello, this is a message from the gRPC client1",
		},
		{
			ID:   2,
			Text: "Hello, this is a message from the gRPC client2",
		},
		{
			ID:   3,
			Text: "Hello, this is a message from the gRPC client3",
		},
	}

	for i := range newMessage {
		newMessage[i].Timestamp = time.Now().UnixNano()
		_, err := c.Send(ctx, &newMessage[i])
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(100 * time.Millisecond)
	}

	r, err := c.Messages(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Get From Server: %s", r.String())
}
