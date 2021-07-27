package broker

import (
	"context"
	"therealbroker/pkg/broker"
)

type Subscriber struct {
	Channel chan broker.Message
	Ctx     context.Context

}

func (s *Subscriber) registerMessage(msg broker.Message){
	select{
	case <-s.Ctx.Done():
		return
	default:
		s.Channel<- msg
	}
}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		Ctx: ctx,
	}
	return newSub
}

