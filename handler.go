package pubsub

import "github.com/tinh-tinh/tinhtinh/v2/core"

type Handler struct {
	core.DynamicProvider
	module core.Module
}

func NewHandler(module core.Module) *Handler {
	return &Handler{
		module: module,
	}
}

type HandleFnc func(sub *Subscriber)

func (h *Handler) Listen(factory HandleFnc, topics ...string) {
	broken := InjectBroker(h.module)
	if broken == nil {
		panic("broken not defined")
	}
	sub := broken.AddSubscriber()
	if len(topics) > 0 {
		for _, topic := range topics {
			broken.Subscribe(sub, topic)

		}
	}

	factory(sub)
}
