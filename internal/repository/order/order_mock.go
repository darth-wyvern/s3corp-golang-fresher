package order

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateOrder(ctx context.Context, tx *sql.Tx, order model.Order) (model.Order, error) {
	args := m.Called(ctx, tx, order)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *Mock) CreateItem(ctx context.Context, tx *sql.Tx, item model.OrderItem) error {
	args := m.Called(ctx, tx, item)
	return args.Error(0)
}

func (m *Mock) GetStatistics(ctx context.Context) ([]Statistics, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Statistics), args.Error(1)
}

func (m *Mock) GetLatestOrder(ctx context.Context, limit int) ([]OrderInfo, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]OrderInfo), args.Error(1)
}

func (m *Mock) GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]Order), args.Get(1).(int64), args.Error(2)
}
