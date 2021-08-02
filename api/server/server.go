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
		err := stream.Send(messageResponse)
		if err != nil {
			return err
		}
	}
	return nil
}
func(s *Server) Fetch(ctx context.Context,fetchReq *pb.FetchRequest) (*pb.MessageResponse, error){
	message, err:= s.broker.Fetch(ctx,fetchReq.GetSubject(),int(fetchReq.GetId()))
	if err!= nil{
		switch err{
		case broker.ErrUnavailable:
			return nil, status.Errorf(codes.Unavailable,"broker has been closed bruh")
		case broker.ErrInvalidID:
			return nil, status.Errorf(codes.InvalidArgument,"invalid ID has been entered")
		case broker.ErrExpiredID:
			return nil, status.Errorf(codes.DeadlineExceeded,"message has been expired")
		}
	}
	messageResponse := &pb.MessageResponse{Body: []byte(message.Body)}
	return messageResponse, nil
}