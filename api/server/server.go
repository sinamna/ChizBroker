package server

import (
	"context"
	pb "therealbroker/api/proto"
)
type Server struct{}

func(s *Server) Publish(context.Context, *pb.PublishRequest) (*pb.PublishResponse, error) {


}
func(s *Server) Subscribe(*pb.SubscribeRequest, pb.Broker_SubscribeServer) error{


}
func(s *Server) Fetch(context.Context, *pb.FetchRequest) (*pb.MessageResponse, error){

}