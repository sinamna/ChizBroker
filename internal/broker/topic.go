package broker

import (
	"context"
	//"fmt"
	"sync"
	"therealbroker/pkg/broker"
	"time"
)

var MessageID = AutoIncId{id: 1}

type Topic struct {
	sync.Mutex
	Name string
	//Subscribers     []*Subscriber
	Subscribers     map[int]*Subscriber
	Messages        map[int]*broker.Message
	//IDs             map[int]struct{}
	//Buffer          []broker.Message
	//pubSignal       chan struct{}
	//signalAvailable bool
	expireSignal    chan int
	subDeleteChan   chan *Subscriber
	subAddChan      chan *Subscriber
	msgPubChan      chan *broker.Message

}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message)
	newSub := CreateNewSubscriber(ctx, ch, t.subDeleteChan)
	t.subAddChan <- newSub
	return ch
}

func (t *Topic) PublishMessage(msg broker.Message) int {

	t.Lock()

	messageId := MessageID.GetID()
	//t.Buffer = append(t.Buffer, msg)
	//t.IDs[messageId] = struct{}{}
	t.Messages[messageId] = nil
	if msg.Expiration != 0 {
		t.Messages[messageId] = &msg
		go t.expireMessage(messageId, msg.Expiration)
	}
	t.Unlock()
	//t.pubSignal <- struct{}{}
	t.msgPubChan <- &msg
	return messageId
}
func (t *Topic) actionListener() {
	for {
		select {
		case newSub := <-t.subAddChan:
			t.Subscribers[newSub.Id] = newSub
		case subscriber := <-t.subDeleteChan:
			delete(t.Subscribers, subscriber.Id)
		case msg := <-t.msgPubChan:
			var wg sync.WaitGroup
			for _, sub := range t.Subscribers {
				sub := sub
				wg.Add(1)
				go func() {
					//sub.RegisterChannel <- msg
					sub.SendMessages(*msg)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}
}

func (t *Topic) Fetch(id int) (broker.Message, error) {
	var fetchedMessage broker.Message

	t.Lock()
	defer t.Unlock()

	message, existed := t.Messages[id]
	if !existed{
		return fetchedMessage, broker.ErrInvalidID
	}else{
		if message == nil{
			return broker.Message{}, broker.ErrExpiredID
		}else{
			fetchedMessage = *message
		}
	}
	return fetchedMessage, nil
}
func (t *Topic) WatchForExpiration() {
	for {
		select {
		case id := <-t.expireSignal:
			t.Lock()
			t.Messages[id]=nil
			t.Unlock()
		}
	}

}
func (t *Topic) expireMessage(id int, expiration time.Duration) {
	select {
	case <-time.After(expiration):
		t.expireSignal <- id
	}
}
func NewTopic(name string) *Topic {
	newTopic := &Topic{
		Name:            name,
		Subscribers:     map[int]*Subscriber{},
		Messages:        map[int]*broker.Message{},
		//IDs:             map[int]struct{}{},
		//Buffer:          make([]broker.Message, 0),
		//pubSignal:       make(chan struct{}),
		expireSignal:    make(chan int),
		//signalAvailable: false,
		subDeleteChan:   make(chan *Subscriber),
		subAddChan:      make(chan *Subscriber),
		msgPubChan: make(chan *broker.Message),
	}
	go newTopic.actionListener()
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
