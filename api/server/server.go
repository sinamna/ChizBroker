package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "therealbroker/api/proto"
	"therealbroker/pkg/broker"
	"time"

)
type Server struct{
	broker broker.Broker
}

func(s *Server) Publish(ctx context.Context,publishReq *pb.PublishRequest) (*pb.PublishResponse, error) {
	message:= broker.Message{
		Body: string(publishReq.GetBody()),
		Expiration: time.Duration(publishReq.ExpirationSeconds),
	}
	messageId, err := s.broker.Publish(ctx,publishReq.GetSubject(),message)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable,"Broker has been closed bruh.")
	}
	response := &pb.PublishResponse{Id: int32(messageId)}
	return response,nil

}
func(s *Server) Subscribe(*pb.SubscribeRequest, pb.Broker_SubscribeServer) error{


}
func(s *Server) Fetch(context.Context, *pb.FetchRequest) (*pb.MessageResponse, error){

}