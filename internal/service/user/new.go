package user

import (
	"context"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
)

type IService interface {
	// CreateUser creates a new user by given InputUser param, and returns the created user.
	CreateUser(ctx context.Context, input InputUser) (model.User, error)

	// GetUsers returns all users by given InputGetUser param.
	GetUsers(ctx context.Context, input InputGetUser) ([]model.User, int64, error)

	// UpdateUser updates a user by given InputUser param.
	UpdateUser(ctx context.Context, input InputUser) error

	// GetUser returns a user by given "id" param.
	GetUser(ctx context.Context, id int) (model.User, error)

	// DeleteUser delete a user by given "id" param.
	DeleteUser(ctx context.Context, id int) error

	// Login authenticate login data
	Login(ctx context.Context, input LoginInput) (LoginResponse, error)

	// GetStatistics returns statistic of users
	GetStatistics(ctx context.Context, orderLimit int) (SummaryStatistics, error)
}

type impl struct {
	repo repository.IRepo
}

func New(repo repository.IRepo) IService {
	return impl{repo: repo}
}
