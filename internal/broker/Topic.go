package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
	"time"
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
	expireSignal chan int
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx,ch)
	t.lock.Lock()
	t.Subscribers = append(t.Subscribers, newSub)
	t.lock.Unlock()
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

	messageId := MessageID.GetID()
	//fmt.Println("publishing",messageId)

	t.lock.Lock()
	//t.BufferIndex++
	t.Buffer = append (t.Buffer,msg)
	go func(){
		t.pubSignal<- struct{}{}
	}()
	t.IDs[messageId] = struct{}{}

	if msg.Expiration != 0 {
		t.Messages[messageId] = msg
		go t.expireMessage(messageId, msg.Expiration)
	}
	t.lock.Unlock()

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
	t.lock.Lock()
	//fmt.Println("fetching",id)
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
	t.lock.Unlock()
	//fmt.Println("found")
	return message, nil
}
func (t *Topic) WatchForExpiration(){
	for{
		select{
		case id := <- t.expireSignal:
			//fmt.Println("tssss",id)
			t.lock.Lock()
			delete(t.Messages,id)
			t.lock.Unlock()
		}
	}

}
func (t *Topic) expireMessage(id int, expiration time.Duration){
	select{
	case <-time.After(expiration):
		t.expireSignal<-id
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
		expireSignal: make(chan int),
	}
	//go newTopic.registerSub()
	go newTopic.publishListener()
	go newTopic.WatchForExpiration()
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
