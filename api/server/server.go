package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "therealbroker/api/proto"
	broker2 "therealbroker/internal/broker"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/metric"
	"time"
)
type Server struct{
	broker broker.Broker
	pb.UnimplementedBrokerServer
}


func(s Server) Publish(ctx context.Context,publishReq *pb.PublishRequest) (*pb.PublishResponse, error) {
	metric.MethodCalls.WithLabelValues("publish").Inc()
	currentTime := time.Now()
	defer metric.MethodDuration.WithLabelValues("publish").Observe(float64(time.Since(currentTime)))

	message:= broker.Message{
		Body: string(publishReq.GetBody()),
		Expiration: time.Duration(publishReq.ExpirationSeconds),
	}
	messageId, err := s.broker.Publish(ctx,publishReq.GetSubject(),message)
	if err != nil {
		metric.MethodError.WithLabelValues("publish").Inc()
		return nil, status.Errorf(codes.Unavailable,"Broker has been closed bruh.")
	}
	response := &pb.PublishResponse{Id: int32(messageId)}
	return response,nil

}
func(s Server) Subscribe(req *pb.SubscribeRequest,stream pb.Broker_SubscribeServer) error{
	metric.MethodCalls.WithLabelValues("subscribe").Inc()
	currentTime := time.Now()
	defer metric.MethodDuration.WithLabelValues("subscribe").Observe(float64(time.Since(currentTime)))
	metric.ActiveSubscribers.Inc()
	defer metric.ActiveSubscribers.Dec()

	ch, err :=s.broker.Subscribe(context.Background(),req.GetSubject())
	if err!= nil{
		metric.MethodError.WithLabelValues("subscribe").Inc()
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
func(s Server) Fetch(ctx context.Context,fetchReq *pb.FetchRequest) (*pb.MessageResponse, error){
	metric.MethodCalls.WithLabelValues("fetch").Inc()
	currentTime := time.Now()
	defer metric.MethodDuration.WithLabelValues("fetch").Observe(float64(time.Since(currentTime)))

	message, err:= s.broker.Fetch(ctx,fetchReq.GetSubject(),int(fetchReq.GetId()))
	if err!= nil{
		metric.MethodError.WithLabelValues("fetch").Inc()
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

func GetServer() pb.BrokerServer{
	return  &Server{broker: broker2.NewModule()}
}