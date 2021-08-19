package main

import (
	"fmt"
	//"github.com/joho/godotenv"

	//"github.com/pkg/profile"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	pb "therealbroker/api/proto"
	"therealbroker/api/server"
	//"runtime"
	_ "net/http/pprof"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	fmt.Println("couldn't load file")
	//}

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	//defer profile.Start(profile.MemProfile,profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
	//defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	go func(){
		fmt.Println("starting prometheus on 8000")
		http.Handle("/metrics",promhttp.Handler())
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		http.ListenAndServe(":8080",nil)
	}()
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBrokerServer(grpcServer,server.GetServer())
	fmt.Println("starting grpc server on 8086")
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}

}
