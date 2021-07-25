package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

type Module struct {
	// TODO: Add required fields
	closed bool
	Topics map[string]*Topic
}

func NewModule() broker.Broker {
	return &Module{
		closed: false,
		Topics: map[string]*Topic{},
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
	topic, exists := m.Topics[subject]
	if !exists {
		topic = NewTopic(subject)
		m.Topics[subject]=topic
	}
	id := topic.PublishMessage(msg)
	return id,nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.closed {
		return nil, broker.ErrUnavailable
	}



	return make(chan broker.Message), nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.closed {
		return broker.Message{}, broker.ErrUnavailable
	}
	return broker.Message{}, nil
}
