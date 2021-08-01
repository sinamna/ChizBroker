package broker

import (
	"context"
	"fmt"
	"sync"
	"therealbroker/pkg/broker"
	"time"
)

var MessageID = AutoIncId{id: 1}

type Topic struct {
	sync.Mutex
	Name        string
	Subscribers []*Subscriber
	Messages    map[int]broker.Message
	IDs         map[int]struct{}
	Buffer []broker.Message
	pubSignal chan struct{}
	signalAvailable bool
	expireSignal chan int
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx,ch)
	t.Lock()
	t.Subscribers = append(t.Subscribers, newSub)
	t.Unlock()
	return ch
}

func (t *Topic) PublishMessage(msg broker.Message) int {

	t.Lock()
	messageId := MessageID.GetID()
	t.Buffer = append (t.Buffer,msg)
	t.IDs[messageId] = struct{}{}
	if msg.Expiration != 0 {
		t.Messages[messageId] = msg
		go t.expireMessage(messageId, msg.Expiration)
	}
	if !t.signalAvailable {
		go func(){
			t.pubSignal<- struct{}{}
		}()
		t.signalAvailable= true
	}

	t.Unlock()

	return messageId
}
func(t *Topic) publishListener(){
	for{
		select {
		case <-t.pubSignal:

			t.Lock()
			if len(t.Buffer) ==0 {
				continue
			}
			messages:= make([]broker.Message,len(t.Buffer))
			copy(messages,t.Buffer)
			t.Buffer = t.Buffer[:0]
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
			t.signalAvailable=false
			t.Unlock()
			wg.Wait()
		}
	}
}


func(t *Topic) Fetch(id int)(broker.Message,error){
	var fetchedMessage broker.Message

	t.Lock()
	_, existedInPast := t.IDs[id]
	message, exists := t.Messages[id]

	if existedInPast{
		if !exists{
			fmt.Println("invalid")
			return broker.Message{}, broker.ErrExpiredID
		}else{
			fetchedMessage = message
		}
	}else{
		return fetchedMessage,broker.ErrInvalidID
	}
	t.Unlock()

	return fetchedMessage, nil
}
func (t *Topic) WatchForExpiration(){
	for{
		select{
		case id := <- t.expireSignal:
			t.Lock()
			delete(t.Messages,id)
			t.Unlock()
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
		Buffer: make([]broker.Message,0),
		pubSignal: make(chan struct{}),
		expireSignal: make(chan int),
		signalAvailable: false,
	}
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
