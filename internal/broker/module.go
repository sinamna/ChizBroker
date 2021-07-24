package broker

import (
	"context"
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
	panic("implement me")
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	panic("implement me")
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	panic("implement me")
}
