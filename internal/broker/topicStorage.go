package broker

import "sync"

type TopicStorage struct{
	topics map[string]*Topic
	sync.RWMutex
}

func (ts *TopicStorage) GetTopic(name string)(*Topic,bool){
	ts.RLock()
	defer ts.RUnlock()
	topic, err := ts.topics[name]
	return topic,err
}
func (ts *TopicStorage) CreateTopic(name string)*Topic{
	ts.Lock()
	defer ts.Unlock()
	newTopic:= NewTopic(name)
	ts.topics[name] = newTopic
	return newTopic
}
func CreateTopicStorage()*TopicStorage{
	return &TopicStorage{
		topics: map[string]*Topic{},
	}
}

