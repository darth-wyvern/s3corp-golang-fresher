package user

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetUsers(ctx context.Context, input InputGetUser) ([]model.User, int64, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *Mock) CreateUser(ctx context.Context, input InputUser) (model.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.User), args.Error(1)
}
func (m *Mock) UpdateUser(ctx context.Context, input InputUser) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *Mock) GetUser(ctx context.Context, id int) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *Mock) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Mock) Login(ctx context.Context, input LoginInput) (LoginResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(LoginResponse), args.Error(1)
}

func (m *Mock) GetStatistics(ctx context.Context, orderLimit int) (SummaryStatistics, error) {
	args := m.Called(ctx, orderLimit)
	return args.Get(0).(SummaryStatistics), args.Error(1)
}
