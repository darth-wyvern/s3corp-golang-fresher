package user

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *Mock) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(bool), args.Error(1)
}

func (m *Mock) ExistsUserByID(ctx context.Context, id int) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *Mock) GetUsers(ctx context.Context, input Filter) (model.UserSlice, int64, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.UserSlice), args.Get(1).(int64), args.Error(2)
}

func (m *Mock) UpdateUser(ctx context.Context, updateUser model.User) (int64, error) {
	args := m.Called(ctx, updateUser)
	return args.Get(0).(int64), args.Error(1)
}

func (m *Mock) GetUser(ctx context.Context, id int) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *Mock) DeleteUser(ctx context.Context, id int) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *Mock) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *Mock) GetStatistics(ctx context.Context) (SummaryStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(SummaryStatistics), args.Error(1)
}
