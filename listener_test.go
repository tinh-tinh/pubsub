package pubsub_test

import (
	"fmt"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/pubsub/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Listener_MultipleMessages(t *testing.T) {
	const TOPIC = "test-topic"
	var receivedCount int32
	var wg sync.WaitGroup
	const MESSAGE_COUNT = 100
	wg.Add(MESSAGE_COUNT) // Expecting 100 messages

	// Define a simple handler that counts messages
	messageHandler := func(module core.Module) core.Provider {
		handler := pubsub.NewListener(module)

		handler.Listen(func(msg *pubsub.Message) {
			newCount := atomic.AddInt32(&receivedCount, 1)
			if newCount <= int32(MESSAGE_COUNT) {
				wg.Done()
			}
		}, TOPIC)

		return handler
	}

	// Create a controller that publishes messages
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			broker := pubsub.InjectBroker(module)

			// Publish 100 times
			go func() {
				for i := range MESSAGE_COUNT {
					broker.Publish(TOPIC, fmt.Sprintf("message %d", i))
					// time.Sleep(1 * time.Millisecond) // Optional: small delay if needed
				}
			}()

			return ctx.JSON(core.Map{"status": "ok"})
		})

		return ctrl
	}

	testModule := func(module core.Module) core.Module {
		return module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
			Providers:   []core.Providers{messageHandler},
		})
	}

	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				pubsub.ForRoot(pubsub.BrokerOptions{}),
				testModule,
			},
		})
	}

	app := core.CreateFactory(appModule)
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	// Trigger the publish
	_, err := testClient.Post(testServer.URL+"/test", "text/plain", nil)
	require.Nil(t, err)

	// Wait for messages with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		// Timeout
		t.Fatal("Timeout waiting for messages")
	}

	require.Equal(t, int32(MESSAGE_COUNT), atomic.LoadInt32(&receivedCount), "Should have received 100 messages")
}
