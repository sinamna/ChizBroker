package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

type Subscriber struct{
	Channel chan broker.Message
	Ctx context.Context
}
func(s *Subscriber) publishMessage(msg broker.Message)error{
	select{
	case <-s.Ctx.Done():
		return broker.ErrExpiredID
	default:
		s.Channel<-msg
		return nil
	}
}
type Topic struct{
	Name string
	Subscribers []*Subscriber
	Messages map[string]broker.Message
}
func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message{
	newSub := &Subscriber{
		Channel: make(chan broker.Message),
		Ctx: ctx,
	}
	t.Subscribers = append(t.Subscribers, newSub)
	return newSub.Channel
}
func(t *Topic) PublishMessage(msg broker.Message){

}
func NewTopic(name string)*Topic{
	subscribers := make([]Subscriber,0)
	return &Topic{
		Name: name,
		Subscribers: subscribers,
		Messages: map[string]broker.Message{},
	}
}