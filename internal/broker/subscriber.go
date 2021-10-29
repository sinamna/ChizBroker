package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

var subscriberId = AutoIncId{id: 1}

type Subscriber struct {
	//sync.Mutex
	Id              int
	Channel         chan broker.Message
	Ctx             context.Context
	unSubSignal     chan *Subscriber
	RegisterChannel chan *broker.Message
	messages        []*broker.Message
}

func (s *Subscriber) SendMessages() {
	for {

		select {
		case <-s.Ctx.Done():
			go func() { s.unSubSignal <- s }()
			return
		case msg := <-s.RegisterChannel:
			s.Channel<-*msg
			//fmt.Println("message published")
		}

	}
}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message, unSubSignal chan *Subscriber) *Subscriber {
	newSub := &Subscriber{
		Id:              subscriberId.GetID(),
		Channel:         ch,
		Ctx:             ctx,
		unSubSignal:     unSubSignal,
		RegisterChannel: make(chan *broker.Message),
		messages:        make([]*broker.Message, 0),
	}
	go newSub.SendMessages()
	return newSub
}
