package price

import (
	"github.com/tinh-tinh/pubsub"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	priceModule := module.New(core.NewModuleOptions{
		Imports:     []core.Module{pubsub.ForFeature("BTC", "ETH", "SOL")},
		Controllers: []core.Controller{Controller},
		Providers:   []core.Provider{Service},
	})

	return priceModule
}
