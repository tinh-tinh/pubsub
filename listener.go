package pubsub

import "github.com/tinh-tinh/tinhtinh/v2/core"

type Listener struct {
	core.DynamicProvider
	module core.Module
}

func NewListener(module core.Module) *Listener {
	return &Listener{
		module: module,
	}
}

type ListenFnc func(sub *Message)

func (h *Listener) Listen(factory ListenFnc, topics ...string) {
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

	go (func(sub *Subscriber) {
		for msg := range sub.GetMessages() {
			factory(msg)
		}
	})(sub)
}
