package pubsub

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
)

const BROKER core.Provide = "broker"
const SUBSCRIBER core.Provide = "subscriber"

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

func InjectBroker(module *core.DynamicModule) *Broker {
	broker, ok := module.Ref(BROKER).(*Broker)
	if !ok {
		return nil
	}
	return broker
}

func ForFeature(topics ...string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		subModule := module.New(core.NewModuleOptions{})
		subModule.NewProvider(core.ProviderOptions{
			Name: SUBSCRIBER,
			Factory: func(param ...interface{}) interface{} {
				fmt.Println(param...)
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

func InjectSubscriber(module *core.DynamicModule) *Subscriber {
	subscriber, ok := module.Ref(SUBSCRIBER).(*Subscriber)
	if !ok {
		return nil
	}
	return subscriber
}
