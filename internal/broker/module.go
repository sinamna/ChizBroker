package broker

import (
	"context"
	"log"
	"sync"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
)

type Module struct {
	sync.Mutex
	closed bool
	Topics map[string]*Topic
	DB repository.Database
}

func NewModule() broker.Broker {
	db, err:= repository.GetPostgreDB()
	if err!=nil{
		log.Fatalln(err)
		return nil
	}
	return &Module{
		closed: false,
		Topics: map[string]*Topic{},
		DB: db,
	}
}

func (m *Module) Close() error {
	if m.closed {
		return broker.ErrUnavailable
	}
	m.closed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.closed {
		return -1, broker.ErrUnavailable
	}
	m.Lock()
	topic, exists := m.Topics[subject]
	if !exists {

		topic = NewTopic(subject, m.DB)
		m.Topics[subject]=topic
	}
	m.Unlock()
	id := topic.PublishMessage(msg)
	return id,nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.closed {
		return nil, broker.ErrUnavailable
	}
	//channel := make(chan broker.Message)
	m.Lock()
	topic, exists := m.Topics[subject]
	if !exists {
		topic = NewTopic(subject, m.DB)
		m.Topics[subject]=topic
	}
	m.Unlock()

	channel:= topic.RegisterSubscriber(ctx)
	return channel, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.closed {
		return broker.Message{}, broker.ErrUnavailable
	}
	m.Lock()
	topic, exists:= m.Topics[subject]
	m.Unlock()

	if !exists{
		log.Fatalln("invalid topic")
	}
	return topic.Fetch(id)
}
