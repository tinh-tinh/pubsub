package price

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	priceModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{Controller},
		Providers:   []core.Provider{Service},
	})

	return priceModule
}
