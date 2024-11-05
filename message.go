package pubsub

type Message struct {
	topic   string
	content interface{}
}

// NewMessage returns a new Message with the given topic and content.
func NewMessage(topic string, content interface{}) *Message {
	return &Message{
		topic:   topic,
		content: content,
	}
}

// GetTopic returns the topic of the message.
func (m *Message) GetTopic() string {
	return m.topic
}

// GetContent returns the content of the message.
func (m *Message) GetContent() interface{} {
	return m.content
}

type MessageChannel chan Message
