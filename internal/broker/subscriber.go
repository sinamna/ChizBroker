package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
)

var subscriberId = AutoIncId{id: 1}

type Subscriber struct {
	sync.Mutex
	Id              int
	Channel         chan broker.Message
	Ctx             context.Context
	unSubSignal     chan *Subscriber
	RegisterChannel chan *broker.Message
	//SendChannel     chan *broker.Message
	//messageSignal      chan struct{}
	//Messages        *LinkedList
	//lastSent        *broker.Message
	//messageCounter  int
}

func (s *Subscriber) SendMessages() {
	for {
		select {
		case <-s.Ctx.Done():
			go func() { s.unSubSignal <- s }()
			return
		case msg := <-s.RegisterChannel:
			s.Channel<-*msg
			//s.Messages.Add(msg)
			//s.messageCounter++
		//default:
		//	if s.messageCounter > 0 {
		//		msg := s.Messages.GetMessage()
		//		if s.lastSent != msg && msg != nil {
		//			//go func() {
		//			fmt.Println("ts")
		//			s.SendChannel <- msg
		//			s.messageCounter--
		//			s.Messages.NextMessage()
		//			//}()
		//		}
		//	}

		}
	}
}
//func (s *Subscriber) broadcastMessage(queue chan *broker.Message){
//	for {
//		select {
//		case msgToSend := <-queue:
//			s.Channel <- *msgToSend
//			//s.messageCounter--
//			//s.lastSent = msgToSend
//		}
//	}
//}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message, unSubSignal chan *Subscriber) *Subscriber {
	newSub := &Subscriber{
		Id:              subscriberId.GetID(),
		Channel:         ch,
		Ctx:             ctx,
		unSubSignal:     unSubSignal,
		//Messages:        CreateNewList(),
		//SendChannel:     make(chan *broker.Message),
		RegisterChannel: make(chan *broker.Message),
		//messageSignal:      make(chan struct{}),
		//lastSent:        nil,
		//messageCounter:  0,
	}
	go newSub.SendMessages()
	return newSub
}
