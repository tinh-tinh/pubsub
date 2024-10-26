package pubsub

import (
	"fmt"
	"time"
)

type Message struct {
	Topic   string
	Content interface{}
}

func NewMessage(topic string, content interface{}) *Message {
	return &Message{
		Topic:   topic,
		Content: content,
	}
}

func (m *Message) GetTopic() string {
	return m.Topic
}

func (m *Message) GetContent() interface{} {
	return m.Content
}

type MessageChannel chan Message

func Subscribe(channel MessageChannel) {
	go func() {
		for {
			message := <-channel
			fmt.Printf("%v: Received message on topic %s: %v\n", time.Now().Format("2006-01-02 15:04:05"), message.GetTopic(), message.GetContent())
		}
	}()
}
