package pubsub

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
)

type Subscriber struct {
	ID       string          // ID of subscriber
	messages chan *Message   // Message channel
	topics   map[string]bool // Topics it is subscribed to
	active   bool            // It given subscriber is active
	mutex    sync.RWMutex
}

// NewSubscriber creates and returns a new Subscriber with a unique ID.
//
// It generates a random ID for the subscriber and initializes the subscriber
// with an empty message channel and an empty topic map. The subscriber is
// marked as active. If there's an error during ID generation, the program
// will log a fatal error and terminate.
func NewSubscriber() (string, *Subscriber) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%X0%X", b[0:4], b[4:8])
	return id, &Subscriber{
		ID:       id,
		messages: make(chan *Message),
		topics:   map[string]bool{},
		active:   true,
	}
}

// AddTopic adds the given topic to the subscriber.
//
// The subscriber is added to the list of subscribers for the specified topic.
// When a message is published to the topic, the subscriber will receive the
// message.
//
// The subscriber is not added if it is already subscribed to the topic.
func (s *Subscriber) AddTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.topics[topic] = true
}

// RemoveTopic removes the given topic from the subscriber.
//
// This function deletes the topic from the subscriber's list of subscribed topics.
// After this operation, the subscriber will no longer receive messages from the
// specified topic. If the topic does not exist in the subscriber's list, the
// function performs no action.
func (s *Subscriber) RemoveTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.topics, topic)
}

// GetTopic returns the list of topics to which the subscriber is subscribed.
func (s *Subscriber) GetTopic() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	topics := []string{}
	for topic := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

// Destruct marks the subscriber as inactive and closes the message channel.
//
// The subscriber will no longer receive messages and resources associated with
// the subscriber are released.
func (s *Subscriber) Destruct() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.active = false
	close(s.messages)
}

// Signal sends the given message to the subscriber.
//
// The message is sent to the subscriber only if the subscriber is active.
// If the subscriber is inactive, the message is not sent.
func (s *Subscriber) Signal(msg *Message) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.active {
		s.messages <- msg
		// fmt.Printf("%v Receiver message %s from topic: %s\n", time.Now().Format("2006-01-01 08:00:00"), msg.GetContent(), msg.GetTopic())
	}
}

// GetMessages returns the message channel of the subscriber.
//
// The channel can be used to receive messages sent to the subscriber.
func (s *Subscriber) GetMessages() chan *Message {
	return s.messages
}

// Listen is a goroutine that listens to the subscriber's message channel.
//
// The goroutine will run indefinitely and will receive messages from the
// subscriber's message channel. Each message is printed to the standard output.
func (s *Subscriber) Listen() {
	for {
		if msg, ok := <-s.messages; ok {
			fmt.Printf("Subscriber %s, received: %s from topic: %s\n", s.ID, msg.GetContent(), msg.GetTopic())
		}
	}
}
