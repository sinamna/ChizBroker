package repository

import (
	"sync"
	"therealbroker/pkg/broker"
)

type mapMemory struct{
	sync.RWMutex
	messages map[int]*broker.Message
}
func (m *mapMemory) SaveMessage(id int, msg broker.Message){
	m.Lock()
	defer m.Unlock()
	if msg.Expiration != 0 {
		m.messages[id]=&msg
	}else{
		m.messages[id]=nil
	}
}
func (m *mapMemory) FetchMessage(id int)(broker.Message,error){
	m.RLock()
	defer m.RUnlock()
	var fetchedMessage broker.Message
	message, existed := m.messages[id]
	if !existed {
		return fetchedMessage, broker.ErrInvalidID
	} else {
		if message == nil {
			return broker.Message{}, broker.ErrExpiredID
		} else {
			fetchedMessage = *message
		}
	}
	return fetchedMessage, nil
}
func (m *mapMemory) DeleteMessage(id int){
	m.Lock()
	defer m.Unlock()
	delete(m.messages,id)
}

func GetInMemoryDB()Database{
	return &mapMemory{
		messages: map[int]*broker.Message{},
	}
}