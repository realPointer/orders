package service

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/realPointer/orders/internal/model"
	"github.com/realPointer/orders/internal/repository"
)

// orderserviceCounter is a counter for the operation on orders.
var OrderServiceCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "orders_service_operations_total",
		Help: "Number of operations on the orderservice",
	},
	[]string{
		"operation",
	},
)

// OrderCreateParams represents parameters for order creation
type OrderCreateParams struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int
	Delivery          model.Delivery
	Payment           model.Payment
	Items             []model.Item
}

type OrderService struct {
	config     *config.Config
	repository *repository.OrderRepository
	cache      *OrderCache
}

func (s *OrderService) Get(ctx context.Context, orderUID string) (model.Order, error) {
	// Try cache first
	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	// If not in cache, get from DB
	order, err := s.repository.Find(ctx, orderUID)
	if err != nil {
		return model.Order{}, err
	}

	// Update cache
	s.cache.Set(order)
	return order, nil
}

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	OrderServiceCounter.WithLabelValues("get_all").Inc()

	return s.repository.FindAll(ctx)
}

func NewOrderService(config *config.Config, repository *repository.OrderRepository) *OrderService {
	return &OrderService{
		config:     config,
		repository: repository,
	}
}

// Create creates a new order with all related data
func (s *OrderService) Create(ctx context.Context, params OrderCreateParams) (string, error) {
	OrderServiceCounter.WithLabelValues("create").Inc()

	order := model.Order{
		OrderUID:          params.OrderUID,
		TrackNumber:       params.TrackNumber,
		Entry:             params.Entry,
		Locale:            params.Locale,
		InternalSignature: params.InternalSignature,
		CustomerID:        params.CustomerID,
		DeliveryService:   params.DeliveryService,
		ShardKey:          params.ShardKey,
		SmID:              params.SmID,
		DateCreated:       time.Now(),
		Delivery:          params.Delivery,
		Payment:           params.Payment,
		Items:             params.Items,
	}

	return s.repository.Create(ctx, order)
}
