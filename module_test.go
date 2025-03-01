package pubsub_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/pubsub/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Module(t *testing.T) {
	const PRICE core.Provide = "price"

	type PriceService struct {
		Message interface{}
	}

	service := func(module core.Module) core.Provider {
		service := module.NewProvider(core.ProviderOptions{
			Name: PRICE,
			Factory: pubsub.Listener(module, func(s *pubsub.Subscriber) interface{} {
				priceService := &PriceService{}
				go (func() {
					msg, ok := <-s.GetMessages()
					if ok {
						priceService.Message = msg.GetContent()
					}
				})()

				return priceService
			}),
		})

		return service
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("prices")

		ctrl.Post("", func(ctx core.Ctx) error {
			broker := pubsub.InjectBroker(module)
			go broker.Publish("BTC", "haha")
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

	priceModule := func(module core.Module) core.Module {
		priceModule := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{pubsub.ForFeature("BTC", "ETH", "SOL")},
			Controllers: []core.Controllers{controller},
			Providers:   []core.Providers{service},
		})

		return priceModule
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				pubsub.ForRoot(),
				priceModule,
			},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Post(testServer.URL+"/api/prices", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/prices")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	type Response struct {
		Data interface{} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var response Response
	err = json.Unmarshal(body, &response)
	require.Nil(t, err)

	require.Equal(t, "haha", response.Data)
}
