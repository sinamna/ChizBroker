package broker

import (
	"context"
	"sync"
	"therealbroker/internal/utils"
	"therealbroker/pkg/broker"
)

var MessageID = AutoIncId{id: 1}


type Topic struct {
	lock sync.Mutex
	Name        string
	Subscribers []*Subscriber
	Messages    map[int]broker.Message
	IDs         map[int]struct{}
	Buffer []broker.Message
	//BufferIndex int
	pubSignal chan struct{}
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
	//for _, sub := range t.Subscribers {
	//	sub.registerMessage(msg)
	//	//TODO: Can we add concurrency here?
	//}

	t.lock.Lock()
	//t.BufferIndex++
	t.Buffer = append (t.Buffer,msg)
	go func(){
		t.pubSignal<- struct{}{}
	}()
	t.lock.Unlock()


	messageId := MessageID.GetID()
	t.IDs[messageId] = struct{}{}

	if msg.Expiration != 0 {
		t.Messages[messageId] = msg
		go utils.WatchForExpiration(t.Messages, messageId, msg.Expiration)
	}
	return messageId
}
func(t *Topic) publishListener(){
	for{
		select {
		case <-t.pubSignal:
			if len(t.Buffer) ==0 {
				continue
			}
			t.lock.Lock()
			messages:= make([]broker.Message,len(t.Buffer))
			copy(messages,t.Buffer)
			t.Buffer = t.Buffer[:0]
			t.lock.Unlock()
			var wg sync.WaitGroup
			for _,sub:= range t.Subscribers{
				sub := sub
				wg.Add(1)
				go func(){
					for _, message:= range messages{
						sub.registerMessage(message)
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}
}
func NewTopic(name string) *Topic {
	subscribers := make([]*Subscriber, 0)
	newTopic := &Topic{
		Name:        name,
		Subscribers: subscribers,
		Messages:    map[int]broker.Message{},
		IDs:         map[int]struct{}{},
		//BufferIndex: -1,
		Buffer: make([]broker.Message,0),
		pubSignal: make(chan struct{}),
	}
	//go newTopic.registerSub()
	go newTopic.publishListener()
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
