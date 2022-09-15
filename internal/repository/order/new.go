package order

import (
	"context"
	"database/sql"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type IOrder interface {
	// CreateOrder create new order
	CreateOrder(ctx context.Context, tx *sql.Tx, order model.Order) (model.Order, error)

	// CreateItem create new order item
	CreateItem(ctx context.Context, tx *sql.Tx, item model.OrderItem) error

	// GetStatistics returns summary statistic of orders.
	GetStatistics(ctx context.Context) ([]Statistics, error)

	// GetLatestOrder returns the list of the latest orders.
	GetLatestOrder(ctx context.Context, limit int) ([]OrderInfo, error)
	// GetOrders
	// GetOrders returns order list from database
	GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error)
}

type impl struct {
	db *sql.DB
}

func New(db *sql.DB) IOrder {
	return &impl{db: db}
}
