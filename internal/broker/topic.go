package broker

import (
	"context"
	"fmt"

	//"fmt"
	//"log"

	//"fmt"
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
	Messages      map[int]*broker.Message
	expireSignal  chan int
	subDeleteChan chan *Subscriber
	subAddChan    chan *Subscriber
	msgPubChan    chan *broker.Message
}

func (t *Topic) RegisterSubscriber(ctx context.Context) chan broker.Message {
	ch := make(chan broker.Message, 20)
	newSub := CreateNewSubscriber(ctx, ch, t.subDeleteChan)
	t.subAddChan <- newSub
	return ch
}

func (t *Topic) PublishMessage(msg broker.Message) int {
	Id := MessageID.GetID()
	if msg.Expiration != 0 {
		t.db.SaveMessage(Id, msg, t.Name)
		//if err != nil {
		//	//fmt.Printf("%#v %s\n",msg, t.Name)
		//	fmt.Println(err)
		//}
	}
	t.msgPubChan <- &msg
	return Id
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
	var fetchedMessage broker.Message

	t.Lock()
	defer t.Unlock()

	message, existed := t.Messages[id]
	if !existed {
		return fetchedMessage, broker.ErrInvalidID
	} else {
		if message == nil {
			return broker.Message{}, broker.ErrExpiredID
		} else {
			fetchedMessage = *message
		}
	}
	return fetchedMessage, nil
}
func (t *Topic) WatchForExpiration() {
	for {
		select {
		case id := <-t.expireSignal:
			go t.db.DeleteMessage(id,t.Name)
		}
	}

}
func (t *Topic) expireMessage(id int, expiration time.Duration) {
	select {
	case <-time.After(expiration):
		t.expireSignal <- id
	}
}
func (t *Topic) SetDB (db repository.Database){
	t.db=db
}
func NewTopic(name string) *Topic {
	newTopic := &Topic{
		Name:          name,
		Subscribers:   map[int]*Subscriber{},
		Messages:      map[int]*broker.Message{},
		expireSignal:  make(chan int, 3),
		subDeleteChan: make(chan *Subscriber, 3),
		subAddChan:    make(chan *Subscriber),
		msgPubChan:    make(chan *broker.Message),
		//db:            db,
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
