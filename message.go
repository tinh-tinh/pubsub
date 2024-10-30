package pubsub

type Message struct {
	Topic   string
	Content interface{}
}

// NewMessage returns a new Message with the given topic and content.
func NewMessage(topic string, content interface{}) *Message {
	return &Message{
		Topic:   topic,
		Content: content,
	}
}

// GetTopic returns the topic of the message.
func (m *Message) GetTopic() string {
	return m.Topic
}

// GetContent returns the content of the message.
func (m *Message) GetContent() interface{} {
	return m.Content
}

type MessageChannel chan Message
