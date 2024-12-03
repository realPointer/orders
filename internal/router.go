package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/realPointer/orders/internal/handler"
	"github.com/realPointer/orders/internal/handler/order"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP handlers.
func Router() fx.Option {
	return fx.Options(
		// dashboard handler
		fxhttpserver.AsHandler("GET", "", handler.NewDashboardHandler),
		// orders CRUD handlers group
		fxhttpserver.AsHandlersGroup(
			"/orders",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "/:id", order.NewGetorderHandler),
			},
		),
	)
}
