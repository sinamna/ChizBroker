package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

type Subscriber struct{
	Channel chan broker.Message
	Ctx context.Context
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
func NewTopic(name string)*Topic{
	subscribers := make([]Subscriber,0)
	return &Topic{
		Name: name,
		Subscribers: subscribers,
		Messages: map[string]broker.Message{},
	}
}