
package main

import (
"context"
	"fmt"
	"google.golang.org/grpc"
"log"
pb "therealbroker/api/proto"
"time"
)

const address = "localhost:8086"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBrokerClient(conn)
	counter:=0
	for {
		func() {
			c.Publish(context.Background(), &pb.PublishRequest{
				Subject: "test",
				Body:    []byte(fmt.Sprint(counter)),
			})
			counter++
			ctx := context.WithValue(context.Background(), "a", "b")
			_, err = c.Subscribe(ctx, &pb.SubscribeRequest{Subject: "mahdi"})
			time.Sleep(time.Second)
		}()
	}
}
