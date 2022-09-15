package repository

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) User() user.IUser {
	args := m.Called()
	return args.Get(0).(user.IUser)
}

func (m *Mock) Product() product.IProduct {
	args := m.Called()
	return args.Get(0).(product.IProduct)
}

func (m *Mock) Order() order.IOrder {
	args := m.Called()
	return args.Get(0).(order.IOrder)
}

func (m *Mock) Tx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
