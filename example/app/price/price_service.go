package price

import (
	"github.com/tinh-tinh/pubsub"
	"github.com/tinh-tinh/tinhtinh/core"
)

const PRICE core.Provide = "price"

type PriceService struct {
	Message interface{}
}

func Service(module *core.DynamicModule) *core.DynamicProvider {
	service := module.NewProvider(core.ProviderOptions{
		Name: PRICE,
		Factory: func(param ...interface{}) interface{} {
			subscriber := param[0].(*pubsub.Subscriber)
			priceService := &PriceService{}
			go (func() {
				// sleep for 5 sec, and then subscribe for topic DOT for s2
				msg, ok := <-subscriber.GetMessages()
				if ok {
					priceService.Message = msg.GetContent()
				}
			})()

			return priceService
		},
		Inject: []core.Provide{pubsub.SUBSCRIBER},
	})

	return service
}
