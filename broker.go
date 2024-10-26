package pubsub

import (
	"fmt"
	"sync"
)

type Subscribers map[string]*Subscriber

type Broker struct {
	Subscribers Subscribers
	Topics      map[string]Subscribers
	mutex       sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		Subscribers: Subscribers{},
		Topics:      map[string]Subscribers{},
	}
}

func (b *Broker) AddSubscriber() *Subscriber {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	id, s := NewSubscriber()
	b.Subscribers[id] = s
	return s
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.Topics[topic] == nil {
		b.Topics[topic] = Subscribers{}
	}

	s.AddTopic(topic)
	b.Topics[topic][s.ID] = s
	fmt.Printf("%s subscribed for topic: %s\n", s.ID, topic)
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.Topics[topic], s.ID)
	s.RemoveTopic(topic)
	fmt.Printf("%s unsubscribed for topic: %s\n", s.ID, topic)
}

func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.Topics {
		b.Unsubscribe(s, topic)
	}

	b.mutex.Lock()
	delete(b.Subscribers, s.ID)
	b.mutex.Unlock()

	s.Destruct()
}

func (b *Broker) GetSubscribers(topic string) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.Topics[topic])
}

func (b *Broker) Broadcast(msg string, topics ...string) {
	for _, topic := range topics {
		for _, s := range b.Topics[topic] {
			m := NewMessage(msg, topic)
			go (func(s *Subscriber) {
				s.Signal(m)
			})(s)
		}
	}
}

func (b *Broker) Publish(msg string, topic string) {
	b.mutex.RLock()
	bTopics := b.Topics[topic]
	b.mutex.RUnlock()

	for _, s := range bTopics {
		m := NewMessage(topic, msg)
		if !s.Active {
			return
		}

		go (func(s *Subscriber) {
			s.Signal(m)
		})(s)
	}
}
