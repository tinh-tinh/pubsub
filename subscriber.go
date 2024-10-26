package pubsub

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

type Subscriber struct {
	ID       string          // ID of subscriber
	Messages chan *Message   // Message channel
	Topics   map[string]bool // Topics it is subscribed to
	Active   bool            // It given subscriber is active
	mutex    sync.RWMutex
}

func NewSubscriber() (string, *Subscriber) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%X0%X", b[0:4], b[4:8])
	return id, &Subscriber{
		ID:       id,
		Messages: make(chan *Message),
		Topics:   map[string]bool{},
		Active:   true,
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.Topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.Topics, topic)
}

func (s *Subscriber) GetTopic() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	topics := []string{}
	for topic := range s.Topics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber) Destruct() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.Active = false
	close(s.Messages)
}

func (s *Subscriber) Signal(msg *Message) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.Active {
		s.Messages <- msg
		// fmt.Printf("%v Receiver message %s from topic: %s\n", time.Now().Format("2006-01-01 08:00:00"), msg.GetContent(), msg.GetTopic())
	}
}

func (s *Subscriber) GetMessages() chan *Message {
	return s.Messages
}

func (s *Subscriber) Listen() {
	for {
		if msg, ok := <-s.Messages; ok {
			fmt.Printf("Subscriber %s, received: %s from topic: %s\n", s.ID, msg.GetContent(), msg.GetTopic())
		}
	}
}
