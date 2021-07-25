package utils

import (
	"therealbroker/pkg/broker"
	"time"
)

func WatchForExpiration(messages map[int]broker.Message,id int, expiration time.Duration){
	select{
	case <-time.After(expiration):
		delete(messages,id)
	}
}