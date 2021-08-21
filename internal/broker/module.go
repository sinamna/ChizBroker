package broker

import (
	"context"
	//"fmt"
	"sync"

	//"fmt"
	"log"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/repository"
)

type Module struct {
	closed bool
	sync.RWMutex
	topicStorage *TopicStorage
	//topics map[string]*Topic
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
		//topics: map[string]*Topic{},
		topicStorage: CreateTopicStorage(),
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
	//m.Lock()
	//topic, exists := m.topics[subject]
	//if !exists {
	//	topic = NewTopic(subject)
	//	m.topics[subject]=topic
	//	topic.SetDB(m.DB)
	//}
	//m.Unlock()
	topic, exists:= m.topicStorage.GetTopic(subject)
	if !exists{
		topic = m.topicStorage.CreateTopic(subject)
	}
	topic.SetDB(m.DB)

	id := topic.PublishMessage(msg)
	return id,nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.closed {
		return nil, broker.ErrUnavailable
	}
	//m.Lock()
	//topic, exists := m.topics[subject]
	//if !exists {
	//	topic = NewTopic(subject)
	//	m.topics[subject]=topic
	//	topic.SetDB(m.DB)
	//}
	//m.Unlock()
	topic, exists:= m.topicStorage.GetTopic(subject)
	if !exists{
		topic = m.topicStorage.CreateTopic(subject)
	}
	topic.SetDB(m.DB)

	channel:= topic.RegisterSubscriber(ctx)
	return channel, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.closed {
		return broker.Message{}, broker.ErrUnavailable
	}

	topic, exists:= m.topicStorage.GetTopic(subject)
	if !exists{
		log.Fatalln("invalid topic")
	}
	//m.RLock()
	//topic, exists := m.topics[subject]
	//if !exists{
	//	fmt.Println("invalid topic")
	//}
	//m.RUnlock()
	return topic.Fetch(id)
}
