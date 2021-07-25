package broker

import (
	"context"
	"fmt"
	"sync"
	"therealbroker/internal/utils"
	"therealbroker/pkg/broker"
)

var MessageID = AutoIncId{id: 1}

type Subscriber struct {
	sync.RWMutex
	Channel chan broker.Message
	ready chan struct{}
	Ctx     context.Context
	Buffer []broker.Message
	BufferIndex int
}

func (s *Subscriber) registerMessage(msg broker.Message){
	s.Lock()
	defer s.Unlock()
	s.BufferIndex++
	if s.BufferIndex <20 {
			s.Buffer[s.BufferIndex]=msg
	}else{
		s.Buffer = append(s.Buffer,msg)
	}
}
func (s *Subscriber) publishMessage(){
	for{
		select {
		case <-s.Ctx.Done():
			return
		case <-s.ready:
			fmt.Println("got ready")
			s.Lock()
			msgs := s.Buffer[:s.BufferIndex]
			s.BufferIndex=-1
			s.Unlock()
			go func(){
				for _,msg := range msgs{
					s.Channel<-msg
				}
			}()
		default:
			fmt.Println("default")
		}
	}

}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		ready: make(chan struct{}),
		Ctx: ctx,
		Buffer: make([]broker.Message,20),
		BufferIndex: -1,
	}
	go newSub.publishMessage()
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
		sub.registerMessage(msg)
		sub := sub
		go func(){sub.ready<- struct{}{}
		}()
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
