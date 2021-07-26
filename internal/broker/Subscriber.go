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
	//s.Lock()
	//s.BufferIndex++
	//if s.BufferIndex <20 {
	//	s.Buffer[s.BufferIndex]=msg
	//}else{
	//	s.Buffer = append(s.Buffer,msg)
	//}
	//go func(){
	//	s.ready<- struct{}{}
	//}()
	//s.Unlock()

}
//func (s *Subscriber) publishMessage(){
//	for{
//		select {
//		case <-s.Ctx.Done():
//			return
//		case <-s.ready:
//			//s.Lock()
//			//msgs := s.Buffer[:s.BufferIndex+1]
//			//s.BufferIndex=-1
//			go func(){
//				for _,msg := range s.Buffer{
//					s.Channel<-msg
//				}
//			}()
//			//s.Unlock()
//
//		}
//	}

//}
func CreateNewSubscriber(ctx context.Context, ch chan broker.Message)*Subscriber{
	newSub:= &Subscriber{
		Channel: ch,
		Ctx: ctx,
	}
	return newSub
}

