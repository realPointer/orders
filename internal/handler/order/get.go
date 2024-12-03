package order

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/realPointer/orders/internal/service"
)

// GetOrderHandler is the http handler to get a order.
type GetOrderHandler struct {
	service *service.OrderService
}

// NewGetorderHandler returns a new GetorderHandler.
func NewGetorderHandler(service *service.OrderService) *GetOrderHandler {
	return &GetOrderHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *GetOrderHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		orderUID := c.Param("id")

		order, err := h.service.Get(c.Request().Context(), orderUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("cannot find order with id %v: %v", orderUID, err))
			}

			return err
		}

		return c.JSON(http.StatusOK, order)
	}
}
