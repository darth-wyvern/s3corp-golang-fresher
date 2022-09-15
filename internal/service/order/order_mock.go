package order

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateOrder(ctx context.Context, input OrderInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *Mock) GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]Order), args.Get(1).(int64), args.Error(2)
}
