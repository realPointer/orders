package service

import (
	"sync"

	"github.com/realPointer/orders/internal/model"
)

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]model.Order),
	}
}

func (c *OrderCache) Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *OrderCache) Get(orderUID string) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}
