package broker

import (
	"context"
	"fmt"

	//"github.com/prometheus/common/log"
	"therealbroker/pkg/repository"

	//"fmt"
	"sync"
	"therealbroker/pkg/broker"
	"time"
)

var MessageID = AutoIncId{id: 1}

type Topic struct {
	sync.Mutex
	Name          string
	db            repository.Database
	Subscribers   map[int]*Subscriber
	//Messages      map[int]*broker.Message
	expireSignal  chan int
	subDeleteChan chan *Subscriber
	subAddChan    chan *Subscriber
	msgPubChan    chan *broker.Message
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message,70)
	newSub := CreateNewSubscriber(ctx, ch, t.subDeleteChan)
	t.subAddChan <- newSub
	return ch
}

func (t *Topic) PublishMessage(msg broker.Message) int {
	//messageId:= MessageID.GetID()
	var id int
	id = -1
	if msg.Expiration != 0 {
		t.Lock()
		id = t.db.SaveMessage(msg, t.Name)
		t.Unlock()
		if id != -1 {
			go t.expireMessage(id, msg.Expiration)
		}
	}
	t.msgPubChan <- &msg
	//fmt.Println("message published ")
	return id
}
func (t *Topic) actionListener() {
	for {
		select {
		case id := <-t.expireSignal:
			go t.db.DeleteMessage(id, t.Name)
			//fmt.Println("deleting")
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
					sub.RegisterChannel <- msg
					//sub.SendMessages(*msg)
					wg.Done()
				}()
			}
			wg.Wait()

		}
	}
}

func (t *Topic) Fetch(id int) (broker.Message, error) {
	//var fetchedMessage broker.Message

	//t.Lock()
	//defer t.Unlock()

	//message, existed := t.Messages[id]
	//if !existed {
	//	return fetchedMessage, broker.ErrInvalidID
	//} else {
	//	if message == nil {
	//		return broker.Message{}, broker.ErrExpiredID
	//	} else {
	//		fetchedMessage = *message
	//	}
	//}
	fetchedMessage, err:= t.db.FetchMessage(id,t.Name)
	if err!=nil{
		fmt.Println("error in fetching: ",err)
		return broker.Message{},broker.ErrInvalidID
	}
	return fetchedMessage, nil
}
//func (t *Topic) WatchForExpiration() {
//	for {
//		select {}
//	}
//
//}
func (t *Topic) expireMessage(id int, expiration time.Duration) {
	select {
	case <-time.After(expiration):
		t.expireSignal <- id
	}
}
func (t *Topic) SetDB(db repository.Database) {
	t.Lock()
	t.db = db
	t.Unlock()
}
func NewTopic(name string) *Topic {
	newTopic := &Topic{
		Name:          name,
		Subscribers:   map[int]*Subscriber{},
		//Messages:      map[int]*broker.Message{},
		expireSignal:  make(chan int, 3),
		subDeleteChan: make(chan *Subscriber, 3),
		subAddChan:    make(chan *Subscriber),
		msgPubChan:    make(chan *broker.Message),
		//db:            db,
	}
	go newTopic.actionListener()
	//go newTopic.WatchForExpiration()
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
