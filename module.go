package pubsub

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

const BROKER core.Provide = "broker"
const SUBSCRIBER core.Provide = "subscriber"

// ForRoot returns a module that provides the broker.
//
// The module provides the broker as a dependency to other modules. The broker
// is the central entity of the pub/sub pattern. It is responsible for managing
// the subscribers and topics.
func ForRoot() core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
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
func InjectBroker(module *core.DynamicModule) *Broker {
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
func ForFeature(topics ...string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
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
func InjectSubscriber(module *core.DynamicModule) *Subscriber {
	subscriber, ok := module.Ref(SUBSCRIBER).(*Subscriber)
	if !ok {
		return nil
	}
	return subscriber
}
