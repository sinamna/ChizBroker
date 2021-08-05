package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "therealbroker/api/proto"
	"therealbroker/api/server"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

func main() {
	//http.Handle("/metrics",promhttp.Handler())
	//err := http.ListenAndServe(":8000", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}

	lis, err := net.Listen("tcp", "localhost:8086")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBrokerServer(grpcServer,server.Server{})
	fmt.Println("server started")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}


}
