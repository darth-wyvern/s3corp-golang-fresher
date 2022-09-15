package order

import (
	"context"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
)

type IOrder interface {
	// CreateOrder create new order from order input
	CreateOrder(ctx context.Context, input OrderInput) error

	// GetOrders returns list of orders which is filterd
	GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error)
}

type impl struct {
	repo repository.IRepo
}

func New(repo repository.IRepo) IOrder {
	return &impl{repo: repo}
}
