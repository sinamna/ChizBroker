package broker

import (
	"context"
	"sync"
	"therealbroker/internal/utils"
	"therealbroker/pkg/broker"
)

var MessageID = AutoIncId{id: 1}


type Topic struct {
	sync.Mutex
	Name        string
	Subscribers []*Subscriber
	Messages    map[int]broker.Message
	IDs         map[int]struct{}
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx,ch)
	t.Subscribers = append(t.Subscribers, newSub)
	return ch
}

//func (t *Topic) registerSub() {
//	for {
//		select {
//		case sub := <-t.subChannel:
//			t.Lock()
//
//			t.Unlock()
//		}
//	}
//}

func (t *Topic) PublishMessage(msg broker.Message) int {
	for _, sub := range t.Subscribers {
		sub.registerMessage(msg)
		//TODO: Can we add concurrency here?
	}
	messageId := MessageID.GetID()
	t.IDs[messageId] = struct{}{}

	if msg.Expiration != 0 {
		t.Messages[messageId] = msg
		go utils.WatchForExpiration(t.Messages, messageId, msg.Expiration)
	}
	return messageId
}

func NewTopic(name string) *Topic {
	subscribers := make([]*Subscriber, 0)
	newTopic := &Topic{
		Name:        name,
		Subscribers: subscribers,
		Messages:    map[int]broker.Message{},
		IDs:         map[int]struct{}{},

	}
	//go newTopic.registerSub()
	return newTopic
}

type AutoIncId struct {
	id int
}

func (ai *AutoIncId) GetID() (id int) {
	id = ai.id
	ai.id++
	return
}
