package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

type Subscriber struct {
	Channel chan broker.Message
	Ctx     context.Context
	unSubSignal chan *Subscriber
}

func (s *Subscriber) registerMessage(msg broker.Message){
	select{
	case <-s.Ctx.Done():
		go func() {s.unSubSignal <- s}()
		return
	default:
		s.Channel<- msg
	}
}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message, unSubSignal chan *Subscriber)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		Ctx: ctx,
		unSubSignal: unSubSignal,
	}
	return newSub
}

