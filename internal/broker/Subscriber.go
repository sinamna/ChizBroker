package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
)

type Subscriber struct {
	sync.RWMutex
	Channel chan broker.Message
	Ctx     context.Context

}

func (s *Subscriber) registerMessage(msg broker.Message){
	s.Channel<-msg
}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		Ctx: ctx,
	}
	return newSub
}

