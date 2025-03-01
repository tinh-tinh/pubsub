package pubsub

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const BROKER core.Provide = "broker"
const SUBSCRIBER core.Provide = "subscriber"

// ForRoot returns a module that provides the broker.
//
// The module provides the broker as a dependency to other modules. The broker
// is the central entity of the pub/sub pattern. It is responsible for managing
// the subscribers and topics.
func ForRoot() core.Modules {
	return func(module core.Module) core.Module {
		pubModule := module.New(core.NewModuleOptions{})
		pubModule.NewProvider(core.ProviderOptions{
			Name:  BROKER,
			Value: NewBroker(),
		})
		pubModule.Export(BROKER)

		return pubModule
	}
}

// InjectBroker returns the broker from the given module.
//
// The broker is the central entity of the pub/sub pattern. It is responsible for
// managing the subscribers and topics. The broker is injected into a module
// by calling ForRoot() and then can be retrieved from the module by calling
// InjectBroker.
func InjectBroker(module core.Module) *Broker {
	broker, ok := module.Ref(BROKER).(*Broker)
	if !ok {
		return nil
	}
	return broker
}

// ForFeature returns a module that provides a subscriber for the specified topics.
//
// The module creates and provides a new subscriber that is automatically subscribed
// to the given topics. The subscriber is associated with the broker, which is injected
// into the module. This allows the subscriber to receive messages published to the
// specified topics.
func ForFeature(topics ...string) core.Modules {
	return func(module core.Module) core.Module {
		subModule := module.New(core.NewModuleOptions{})
		subModule.NewProvider(core.ProviderOptions{
			Name: SUBSCRIBER,
			Factory: func(param ...interface{}) interface{} {
				broker := param[0].(*Broker)
				s := broker.AddSubscriber()
				if len(topics) > 0 {
					for _, topic := range topics {
						broker.Subscribe(s, topic)
					}
				}
				return s
			},
			Inject: []core.Provide{BROKER},
		})
		subModule.Export(SUBSCRIBER)

		return subModule
	}
}

// InjectSubscriber returns the subscriber from the given module.
//
// The subscriber is the entity that receives messages published to the topics it
// is subscribed to. The subscriber is injected into a module by calling
// ForFeature() and then can be retrieved from the module by calling
// InjectSubscriber.
func InjectSubscriber(module core.Module) *Subscriber {
	subscriber, ok := module.Ref(SUBSCRIBER).(*Subscriber)
	if !ok {
		return nil
	}
	return subscriber
}

func Listener(module core.Module, fnc func(*Subscriber) interface{}) core.Factory {
	return func(param ...interface{}) interface{} {
		subscriber := InjectSubscriber(module)
		return fnc(subscriber)
	}
}
