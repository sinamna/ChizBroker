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
	messageId := MessageID.GetID()
	t.lock.Lock()
	//t.BufferIndex++
	t.Buffer = append (t.Buffer,msg)
	go func(){
		t.pubSignal<- struct{}{}
	}()
	t.IDs[messageId] = struct{}{}
	t.lock.Unlock()




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


func(t *Topic) Fetch(id int)(broker.Message,error){
	message, exists:= t.Messages[id]
	if !exists{
		_, existedInPast := t.IDs[id]
		if !existedInPast{
			//fmt.Println("invalid")
			return broker.Message{}, broker.ErrInvalidID
		}else{
			//fmt.Println("expired")
			return broker.Message{}, broker.ErrExpiredID
		}
	}
	//fmt.Println("found")
	return message, nil
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
