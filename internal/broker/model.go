package broker

import "therealbroker/pkg/broker"

type Topic struct{
	Name string
	Subscribers []chan broker.Message
	Messages map[string]broker.Message
}

func NewTopic(name string)*Topic{
	channels := make([]chan broker.Message,0)
	return &Topic{
		Name: name,
		Subscribers: channels,
		Messages: map[string]broker.Message{},
	}
}