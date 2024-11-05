package pubsub

import (
	"fmt"
	"strings"
	"sync"
)

type Subscribers map[string]*Subscriber

type BrokerOptions struct {
	// set this to `true` to use wildcards
	Wildcard bool
	// the delimiter used to segment namespaces
	Delimiter string
	// the maximum number of subscribers per topic
	MaxSubscribers int
}

type Broker struct {
	subscribers Subscribers
	topics      map[string]Subscribers
	mutex       sync.RWMutex
	opt         BrokerOptions
}

// NewBroker returns a new instance of Broker.
//
// The broker is a central entity of the pub/sub pattern. It is responsible for
// managing the subscribers and topics.
//
// When a subscriber subscribes to a topic, the subscriber is added to a list
// of subscribers for that topic. When a message is published to a topic, all
// subscribers of that topic will receive the message.
func NewBroker(opt ...BrokerOptions) *Broker {
	broker := &Broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
	if len(opt) > 0 {
		broker.opt = opt[0]
	}

	return broker
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

	if b.opt.MaxSubscribers != 0 && len(b.subscribers)+1 > b.opt.MaxSubscribers {
		return nil
	}

	id, s := NewSubscriber()
	b.subscribers[id] = s
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

	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}

	s.AddTopic(topic)
	b.topics[topic][s.ID] = s
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

	delete(b.topics[topic], s.ID)
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
	for topic := range s.topics {
		b.Unsubscribe(s, topic)
	}

	b.mutex.Lock()
	delete(b.subscribers, s.ID)
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

	return len(b.topics[topic])
}

// Broadcast sends the given message to all subscribers in all topics.
//
// The message is delivered to all active subscribers in all topics.
//
// The message is sent to the subscribers asynchronously.
func (b *Broker) Broadcast(msg string) {
	for topic := range b.topics {
		for _, s := range b.topics[topic] {
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
	topics := []string{topic}
	if b.opt.Wildcard {
		patterns := strings.Split(topic, b.opt.Delimiter)
		for i := 0; i < len(patterns)-1; i++ {
			topics = append(topics, fmt.Sprintf("%s%s*", patterns[i], b.opt.Delimiter))
		}
	}
	var subscribers []*Subscriber
	for _, tp := range topics {
		for _, subscriber := range b.topics[tp] {
			subscribers = append(subscribers, subscriber)
		}
	}
	b.mutex.RUnlock()

	for _, s := range subscribers {
		m := NewMessage(topic, msg)
		if !s.active {
			return
		}

		go (func(s *Subscriber) {
			s.Signal(m)
		})(s)
	}
}
