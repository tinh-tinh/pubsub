package pubsub

import (
	"fmt"
	"sync"
)

type Subscribers map[string]*Subscriber

type BrokerOptions struct {
	Wildcard    bool
	Delimiter   string
	NewListener bool
}

type Broker struct {
	Subscribers Subscribers
	Topics      map[string]Subscribers
	mutex       sync.RWMutex
}

// NewBroker returns a new instance of Broker.
//
// The broker is a central entity of the pub/sub pattern. It is responsible for
// managing the subscribers and topics.
//
// When a subscriber subscribes to a topic, the subscriber is added to a list
// of subscribers for that topic. When a message is published to a topic, all
// subscribers of that topic will receive the message.
func NewBroker() *Broker {
	return &Broker{
		Subscribers: Subscribers{},
		Topics:      map[string]Subscribers{},
	}
}

// AddSubscriber creates a new subscriber and adds it to the broker.
//
// The subscriber is created and registered with the broker. The subscriber
// can then be used to subscribe to topics and receive messages.
//
// The subscriber is returned by AddSubscriber.
func (b *Broker) AddSubscriber() *Subscriber {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	id, s := NewSubscriber()
	b.Subscribers[id] = s
	return s
}

// Subscribe adds the subscriber to the specified topic.
//
// The subscriber is added to the list of subscribers for the specified topic.
// When a message is published to the topic, the subscriber will receive the
// message.
//
// The subscriber is not added if it is already subscribed to the topic.
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

// Unsubscribe removes the subscriber from the specified topic.
//
// The subscriber is removed from the list of subscribers for the specified
// topic. When a message is published to the topic, the subscriber will no
// longer receive the message.
//
// The subscriber is not removed if it is not currently subscribed to the
// topic.
func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.Topics[topic], s.ID)
	s.RemoveTopic(topic)
	fmt.Printf("%s unsubscribed for topic: %s\n", s.ID, topic)
}

// RemoveSubscriber removes the subscriber from all topics and the broker.
//
// The subscriber is removed from the list of subscribers for all topics.
// When a message is published to the topic, the subscriber will no longer
// receive the message.
//
// The subscriber is then removed from the broker and the resources associated
// with the subscriber are released.
func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.Topics {
		b.Unsubscribe(s, topic)
	}

	b.mutex.Lock()
	delete(b.Subscribers, s.ID)
	b.mutex.Unlock()

	s.Destruct()
}

// GetSubscribers returns the number of subscribers for the given topic.
//
// The number of subscribers includes only active subscribers. Inactive
// subscribers are not counted.
func (b *Broker) GetSubscribers(topic string) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.Topics[topic])
}

// Broadcast sends the given message to all subscribers in all topics.
//
// The message is delivered to all active subscribers in all topics.
//
// The message is sent to the subscribers asynchronously.
func (b *Broker) Broadcast(msg string) {
	for topic := range b.Topics {
		for _, s := range b.Topics[topic] {
			m := NewMessage(topic, msg)
			go (func(s *Subscriber) {
				s.Signal(m)
			})(s)
		}
	}
}

// Publish sends the given message to all subscribers of the specified topic.
//
// The message is delivered to all active subscribers of the specified topic.
// The message is sent to the subscribers asynchronously. If a subscriber is
// inactive, it will not receive the message.
func (b *Broker) Publish(topic string, msg string) {
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
