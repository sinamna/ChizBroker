package broker

import (
	"context"
	"sync"
	"therealbroker/pkg/broker"
)

type Subscriber struct {
	sync.RWMutex
	Channel chan broker.Message
	ready chan struct{}
	Ctx     context.Context
	Buffer []broker.Message
	BufferIndex int
}

func (s *Subscriber) registerMessage(msg broker.Message){
	s.Lock()
	s.BufferIndex++
	if s.BufferIndex <20 {
		s.Buffer[s.BufferIndex]=msg
	}else{
		s.Buffer = append(s.Buffer,msg)
	}
	go func(){
		s.ready<- struct{}{}
	}()
	s.Unlock()

}
func (s *Subscriber) publishMessage(){
	for{
		select {
		case <-s.Ctx.Done():
			return
		case <-s.ready:
			s.Lock()
			msgs := s.Buffer[:s.BufferIndex+1]
			s.BufferIndex=-1
			go func(){
				for _,msg := range msgs{
					s.Channel<-msg
				}
			}()
			s.Unlock()

		}
	}

}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		ready: make(chan struct{}),
		Ctx: ctx,
		Buffer: make([]broker.Message,20),
		BufferIndex: -1,
	}
	go newSub.publishMessage()
	return newSub
}

