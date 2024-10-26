package pubsub

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var availableTopics = map[string]string{
	"BTC": "BITCOIN",
	"ETH": "ETHEREUM",
	"DOT": "POLKADOT",
	"SOL": "SOLANA",
}

func pricePublisher(broker *Broker) {
	// topicKeys := make([]string, len(availableTopics))
	topicValues := make([]string, len(availableTopics))

	for _, v := range availableTopics {
		// topicKeys = append(topicKeys, k)
		topicValues = append(topicValues, v)
	}

	for {
		randValue := topicValues[rand.Intn(len(topicValues))]
		msg := fmt.Sprintf("%s: $%.2f", randValue, rand.Float64())
		go broker.Publish(msg, randValue)

		r := rand.Intn(4)
		time.Sleep(time.Duration(r) * time.Second)
	}
}

func Test_Broker(t *testing.T) {
	broker := NewBroker()
	s1 := broker.AddSubscriber()

	broker.Subscribe(s1, availableTopics["BTC"])
	broker.Subscribe(s1, availableTopics["ETH"])

	s2 := broker.AddSubscriber()
	broker.Subscribe(s2, availableTopics["ETH"])
	broker.Subscribe(s2, availableTopics["SOL"])

	go (func() {
		// sleep for 5 sec, and then subscribe for topic DOT for s2
		time.Sleep(3 * time.Second)
		broker.Subscribe(s2, availableTopics["DOT"])
	})()

	go (func() {
		// s;eep for 5 sec, and then unsubscribe for topic SOL for s2
		time.Sleep(5 * time.Second)
		broker.Unsubscribe(s2, availableTopics["SOL"])
		fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})()

	go (func() {
		// s;eep for 5 sec, and then unsubscribe for topic SOL for s2
		time.Sleep(10 * time.Second)
		broker.RemoveSubscriber(s2)
		fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})()
	// Concurrently publish the values.
	go pricePublisher(broker)
	// Concurrently listens from s1.
	go s1.Listen()
	// Concurrently listens from s2.
	go s2.Listen()
	// to prevent terminate
	fmt.Scanln()
	fmt.Println("Done!")
}
