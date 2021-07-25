package broker

import (
	"context"
	"sync"
	"therealbroker/internal/utils"
	"therealbroker/pkg/broker"
)

var MessageID = AutoIncId{id: 1}

type Subscriber struct {
	sync.RWMutex
	Channel chan broker.Message
	Ctx     context.Context
	Buffer []broker.Message
	BufferIndex int
}

func (s *Subscriber) publishMessage(msg broker.Message) {
	select {
	case <-s.Ctx.Done():
	default:
		s.Channel <- msg
	}
}

func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		Ctx: ctx,
		Buffer: make([]broker.Message,20),
		BufferIndex: 0,
	}
	return newSub
}


type Topic struct {
	Name        string
	Subscribers []*Subscriber
	Messages    map[int]broker.Message
	IDs         map[int]struct{}
	subChannel  chan *Subscriber
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx,ch)
	t.subChannel <- newSub
	return ch
}

func (t *Topic) registerSub() {
	for {
		select {
		case sub := <-t.subChannel:
			t.Subscribers = append(t.Subscribers, sub)
		}
	}
}

func (t *Topic) PublishMessage(msg broker.Message) int {
	for _, sub := range t.Subscribers {
		go sub.publishMessage(msg)

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
		subChannel:  make(chan *Subscriber),
	}
	go newTopic.registerSub()
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
