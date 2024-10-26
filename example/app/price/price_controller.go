package price

import (
	"github.com/tinh-tinh/pubsub"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Controller(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("prices")

	ctrl.Post("", func(ctx core.Ctx) error {
		broker := pubsub.InjectBroker(module)
		go broker.Publish("haha", "BTC")
		return ctx.JSON(core.Map{
			"data": "ok",
		})
	})

	ctrl.Get("", func(ctx core.Ctx) error {
		service := module.Ref(PRICE).(*PriceService)
		return ctx.JSON(core.Map{
			"data": service.Message,
		})
	})

	return ctrl
}
