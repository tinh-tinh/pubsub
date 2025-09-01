# PubSub for Tinh Tinh

<div align="center">
<img alt="GitHub Release" src="https://img.shields.io/github/v/release/tinh-tinh/pubsub">
<img alt="GitHub License" src="https://img.shields.io/github/license/tinh-tinh/pubsub">
<a href="https://codecov.io/gh/tinh-tinh/pubsub" > 
 <img src="https://codecov.io/gh/tinh-tinh/pubsub/graph/badge.svg?token=TS4B5QAO3T"/> 
</a>
<a href="https://pkg.go.dev/github.com/tinh-tinh/pubsub"><img src="https://pkg.go.dev/badge/github.com/tinh-tinh/pubsub.svg" alt="Go Reference"></a>
</div>

<div align="center">
    <img src="https://avatars.githubusercontent.com/u/178628733?s=400&u=2a8230486a43595a03a6f9f204e54a0046ce0cc4&v=4" width="200" alt="Tinh Tinh Logo">
</div>

## Overview

PubSub is a modular, in-memory publish/subscribe (pubsub) system for the Tinh Tinh framework. It enables decoupled communication between different parts of your application using topics and message subscribers.

## Install 

```bash
go get -u github.com/tinh-tinh/pubsub/v2
```

## Features

- **Central Broker:** Manages topics and subscribers.
- **Dynamic Subscribers:** Subscribe to one or many topics dynamically.
- **Asynchronous Messaging:** Message delivery is non-blocking.
- **Topic Patterns:** Supports wildcards and topic delimiters for pattern-based subscriptions.
- **Broadcast Support:** Broadcast a message to all subscribers.
- **Integration with Tinh Tinh Modules:** Use dependency injection for broker and subscribers.
- **Subscriber Limits:** Optionally limit the max number of subscribers.
- **Handler Utility:** Simplifies listening to topics with concise functions and is the recommended way to consume messages.

## Basic Usage

### 1. Register the Broker

```go
import "github.com/tinh-tinh/pubsub/v2"

pubsubModule := pubsub.ForRoot(pubsub.BrokerOptions{
    // Optional: MaxSubscribers, Wildcard, Delimiter, etc.
})
```

### 2. Register Subscribers (Feature Modules)

Subscribe to specific topics:

```go
priceSubModule := pubsub.ForFeature("BTC", "ETH") // subscribe to multiple topics
```

### 3. Use Handler for Consumption (Recommended)

**Instead of using `subscriber.GetMessages()` or `Listener`, use the Handler utility for subscribing to topics and consuming messages.**

```go
import (
    "github.com/tinh-tinh/pubsub/v2"
    "github.com/tinh-tinh/tinhtinh/v2/core"
)

type PriceService struct {
    Message interface{}
}

priceService := func(module core.Module) core.Provider {
    return module.NewProvider(core.ProviderOptions{
        Name:  "PriceService",
        Value: &PriceService{},
    })
}

priceHandler := func(module core.Module) core.Provider {
    handler := pubsub.NewHandler(module)
    priceService := module.Ref("PriceService").(*PriceService)

    handler.Listen(func(msg *pubsub.Message) {
        priceService.Message = msg.GetContent()
    }, "BTC", "ETH", "SOL")

    return handler
}
```

### 4. Use in Controllers

```go
controller := func(module core.Module) core.Controller {
    ctrl := module.NewController("prices")

    ctrl.Post("", func(ctx core.Ctx) error {
        broker := pubsub.InjectBroker(module)
        go broker.Publish("BTC", "hihi")
        return ctx.JSON(core.Map{"data": "ok"})
    })

    ctrl.Get("", func(ctx core.Ctx) error {
        service := module.Ref("PriceService").(*PriceService)
        return ctx.JSON(core.Map{"data": service.Message})
    })

    return ctrl
}
```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or need help, you can:
- Open an issue in the GitHub repository
- Check our documentation
- Join our community discussions
