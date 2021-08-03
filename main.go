package therealbroker

import (

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
	lis, err := net.Listen("tcp", "localhost:80086")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBrokerServer(grpcServer,server.Server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
