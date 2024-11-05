package app

import (
	"github.com/tinh-tinh/pubsub"
	"github.com/tinh-tinh/pubsub/example/app/price"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Module() *core.DynamicModule {
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			pubsub.ForRoot(),
			price.Module,
		},
	})

	return module
}
