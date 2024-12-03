// repository/order.go
package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/realPointer/orders/internal/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Find(ctx context.Context, orderUID string) (model.Order, error) {
	var order model.Order

	// Main order query
	orderQuery, orderArgs, err := sq.
		Select(
			"o.order_uid", "o.track_number", "o.entry", "o.locale",
			"o.internal_signature", "o.customer_id", "o.delivery_service",
			"o.shardkey", "o.sm_id", "o.date_created", "o.oof_shard",
			// Delivery fields
			"d.name", "d.phone", "d.zip", "d.city", "d.address", "d.region", "d.email",
			// Payment fields
			"p.transaction", "p.request_id", "p.currency", "p.provider",
			"p.amount", "p.payment_dt", "p.bank", "p.delivery_cost",
			"p.goods_total", "p.custom_fee",
		).
		From("orders o").
		LeftJoin("deliveries d ON d.order_uid = o.order_uid").
		LeftJoin("payments p ON p.order_uid = o.order_uid").
		Where(sq.Eq{"o.order_uid": orderUID}).
		ToSql()

	if err != nil {
		return order, err
	}

	row := r.db.QueryRowContext(ctx, orderQuery, orderArgs...)
	err = row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		// Delivery
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
		// Payment
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)

	if err != nil {
		return order, err
	}

	// Get items separately
	itemsQuery, itemsArgs, err := sq.
		Select(
			"chrt_id", "track_number", "price", "rid", "name",
			"sale", "size", "total_price", "nm_id", "brand", "status",
		).
		From("items").
		Where(sq.Eq{"order_uid": orderUID}).
		ToSql()

	if err != nil {
		return order, err
	}

	rows, err := r.db.QueryContext(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order

	// Main order query
	orderQuery, orderArgs, err := sq.
		Select(
			"o.order_uid", "o.track_number", "o.entry", "o.locale",
			"o.internal_signature", "o.customer_id", "o.delivery_service",
			"o.shardkey", "o.sm_id", "o.date_created", "o.oof_shard",
			// Delivery fields
			"d.name", "d.phone", "d.zip", "d.city", "d.address", "d.region", "d.email",
			// Payment fields
			"p.transaction", "p.request_id", "p.currency", "p.provider",
			"p.amount", "p.payment_dt", "p.bank", "p.delivery_cost",
			"p.goods_total", "p.custom_fee",
		).
		From("orders o").
		LeftJoin("deliveries d ON d.order_uid = o.order_uid").
		LeftJoin("payments p ON p.order_uid = o.order_uid").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, orderQuery, orderArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
			// Delivery
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email,
			// Payment
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		// Get items for this order
		itemsQuery, itemsArgs, err := sq.
			Select(
				"chrt_id", "track_number", "price", "rid", "name",
				"sale", "size", "total_price", "nm_id", "brand", "status",
			).
			From("items").
			Where(sq.Eq{"order_uid": order.OrderUID}).
			ToSql()

		if err != nil {
			return nil, err
		}

		itemRows, err := r.db.QueryContext(ctx, itemsQuery, itemsArgs...)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item model.Item
			err = itemRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
				&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
				&item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) Create(ctx context.Context, order model.Order) (string, error) {
	// Insert order
	orderQuery, orderArgs, err := sq.
		Insert("orders").
		Columns(
			"order_uid", "track_number", "entry", "locale",
			"internal_signature", "customer_id", "delivery_service",
			"shardkey", "sm_id", "date_created", "oof_shard",
		).
		Values(
			order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
			order.InternalSignature, order.CustomerID, order.DeliveryService,
			order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
		).
		ToSql()

	if err != nil {
		return "", err
	}

	_, err = r.db.ExecContext(ctx, orderQuery, orderArgs...)
	if err != nil {
		return "", err
	}

	// Insert delivery
	deliveryQuery, deliveryArgs, err := sq.
		Insert("deliveries").
		Columns(
			"order_uid", "name", "phone", "zip", "city",
			"address", "region", "email",
		).
		Values(
			order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
			order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
		).
		ToSql()

	if err != nil {
		return "", err
	}

	_, err = r.db.ExecContext(ctx, deliveryQuery, deliveryArgs...)
	if err != nil {
		return "", err
	}

	// Insert payment
	paymentQuery, paymentArgs, err := sq.
		Insert("payments").
		Columns(
			"order_uid", "transaction", "request_id", "currency",
			"provider", "amount", "payment_dt", "bank",
			"delivery_cost", "goods_total", "custom_fee",
		).
		Values(
			order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
			order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
			order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
		).
		ToSql()

	if err != nil {
		return "", err
	}

	_, err = r.db.ExecContext(ctx, paymentQuery, paymentArgs...)
	if err != nil {
		return "", err
	}

	// Insert items
	for _, item := range order.Items {
		itemsQuery, itemsArgs, err := sq.
			Insert("items").
			Columns(
				"order_uid", "chrt_id", "track_number", "price", "rid",
				"name", "sale", "size", "total_price", "nm_id", "brand", "status",
			).
			Values(
				order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
				item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
			).
			ToSql()

		if err != nil {
			return "", err
		}

		_, err = r.db.ExecContext(ctx, itemsQuery, itemsArgs...)
		if err != nil {
			return "", err
		}
	}

	return order.OrderUID, nil
}
