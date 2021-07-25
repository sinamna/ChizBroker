package broker

import (
	"context"
	"therealbroker/internal/utils"
	"therealbroker/pkg/broker"
)
var MessageID = AutoIncId{id: 1}
type Subscriber struct {
	Channel chan broker.Message
	Ctx     context.Context
}

func (s *Subscriber) publishMessage(msg broker.Message) {
	select {
	case <-s.Ctx.Done():
	default:
		s.Channel <- msg
	}
}

type Topic struct {
	Name        string
	Subscribers []*Subscriber
	Messages    map[int]broker.Message
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	newSub := &Subscriber{
		Channel: make(chan broker.Message),
		Ctx:     ctx,
	}
	t.Subscribers = append(t.Subscribers, newSub)
	return newSub.Channel
}

func (t *Topic) PublishMessage(msg broker.Message) {
	for _, sub := range t.Subscribers {
		sub.publishMessage(msg)

		//TODO: Can we add concurrency here?
	}
	if msg.Expiration != 0{
		messageId := MessageID.GetID()
		t.Messages[messageId] = msg
		go utils.WatchForExpiration(t.Messages,messageId,msg.Expiration)
	}

}

func NewTopic(name string) *Topic {
	subscribers := make([]*Subscriber, 0)
	return &Topic{
		Name:        name,
		Subscribers: subscribers,
		Messages:    map[int]broker.Message{},
	}
}

type AutoIncId struct{
	id int
}
func (ai *AutoIncId) GetID()(id int){
	id = ai.id
	ai.id++
	return
}
