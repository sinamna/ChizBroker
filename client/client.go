
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
	ids := make([]int, 0)
	for {
		func() {
			id, err := c.Publish(context.Background(), &pb.PublishRequest{
				Subject: "test",
				Body:    []byte("bruh"),
				ExpirationSeconds: 10,
			})
			fmt.Println(id)
			ids=append(ids, int(id.Id))
			if err != nil {
				fmt.Println(err)
			}
			counter++
			fmt.Println(counter,"sent")
			ctx := context.WithValue(context.Background(), "a", "b")
			ch, _ := c.Subscribe(ctx, &pb.SubscribeRequest{Subject: "test"})
			go func() {
				_, err := ch.Recv()
				//fmt.Println(response)
				if err!= nil{
					fmt.Println(err)
				}
			}()
			//time.Sleep(time.Second)
		}()
	}
}
