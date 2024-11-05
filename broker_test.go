package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/pubsub"
)

func Test_Broker(t *testing.T) {
	broker := pubsub.NewBroker()

	sub := broker.AddSubscriber()
	sub2 := broker.AddSubscriber()

	topic := "topic"
	broker.Subscribe(sub, topic)
	broker.Subscribe(sub2, topic)
	require.NotEmpty(t, sub.GetTopic())
	require.NotEmpty(t, sub2.GetTopic())

	require.Equal(t, broker.GetSubscribers(topic), 2)
	broker.RemoveSubscriber(sub2)
	require.Equal(t, broker.GetSubscribers(topic), 1)

	broker.Unsubscribe(sub, topic)
	require.Empty(t, sub.GetTopic())
}

func Test_Pubsub(t *testing.T) {
	broker := pubsub.NewBroker()

	sub := broker.AddSubscriber()
	sub2 := broker.AddSubscriber()

	topic := "topic"
	broker.Subscribe(sub, topic)

	topic2 := "topic2"
	broker.Subscribe(sub2, topic2)
	require.NotEmpty(t, sub.GetTopic())
	require.NotEmpty(t, sub2.GetTopic())

	go sub.Listen()
	go sub2.Listen()

	go (func() {
		broker.Publish(topic, "hello")
	})()

	go (func() {
		broker.Broadcast("hello")
	})()

	require.NotNil(t, sub.GetMessages())
	require.NotNil(t, sub2.GetMessages())

	fmt.Println(sub.GetMessages())
}

func Test_Options(t *testing.T) {
	broker := pubsub.NewBroker(pubsub.BrokerOptions{
		MaxSubscribers: 10,
	})

	for i := 0; i < 15; i++ {
		s := broker.AddSubscriber()
		if i < 10 {
			require.NotNil(t, s)
		} else {
			require.Nil(t, s)
		}
	}
}
