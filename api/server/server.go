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
func(s *Server) Subscribe(req *pb.SubscribeRequest,stream pb.Broker_SubscribeServer) error{
	ch, err :=s.broker.Subscribe(context.Background(),req.GetSubject())
	if err!= nil{
		return status.Errorf(codes.Unavailable,"Broker has been closed bruh.")
	}
	for message := range ch{
		messageResponse := &pb.MessageResponse{Body: []byte(message.Body)}
		stream.Send(messageResponse)
	}
	return nil
}
func(s *Server) Fetch(context.Context, *pb.FetchRequest) (*pb.MessageResponse, error){

}