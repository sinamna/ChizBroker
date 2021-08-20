
package main

import (
"context"
	"fmt"
	"google.golang.org/grpc"
"log"
pb "therealbroker/api/proto"
//"time"
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
			_, err := c.Publish(context.Background(), &pb.PublishRequest{
				Subject: "test",
				Body:    []byte("bruh"),
				ExpirationSeconds: 10,
			})
			if err != nil {
				fmt.Println(err)
			}
			counter++
			fmt.Println(counter,"sent")
			ctx := context.WithValue(context.Background(), "a", "b")
			ch, _ := c.Subscribe(ctx, &pb.SubscribeRequest{Subject: "test"})
			go func() {
				response, _ := ch.Recv()
				fmt.Println(response.GetBody())
			}()
			//time.Sleep(time.Second)
		}()
	}
}
