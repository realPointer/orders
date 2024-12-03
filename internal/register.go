package internal

import (
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"github.com/realPointer/orders/internal/repository"
	"github.com/realPointer/orders/internal/service"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// services
		fx.Provide(
			repository.NewOrderRepository,
			service.NewOrderService,
		),
		// metrics
		fxmetrics.AsMetricsCollector(service.OrderServiceCounter),
		// probes
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
	)
}
