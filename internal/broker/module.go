package broker

import (
	"context"
	"log"
	"therealbroker/pkg/broker"
)

type Module struct {
	// TODO: Add required fields
}

func NewModule() broker.Broker {
	return &Module{}
}

func (m *Module) Close() error {
	panic("implement me")
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.closed {
		return -1, broker.ErrUnavailable
	}
	return 1,nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.closed {
		return nil, broker.ErrUnavailable
	}
	return nil, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.closed {
		return broker.Message{}, broker.ErrUnavailable
	}
	topic, exists:= m.Topics[subject]
	if !exists{
		log.Fatalln("invalid topic")
	}
	return topic.Fetch(id)
}
