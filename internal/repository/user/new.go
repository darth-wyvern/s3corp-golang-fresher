package user

import (
	"context"
	"database/sql"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type IUser interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user model.User) (model.User, error)

	// ExistsUserByEmail returns true if the user exists
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)

	// ExistsUserByID returns true if the user exists
	ExistsUserByID(ctx context.Context, id int) (bool, error)

	// GetUsers returns all users by given Filter param
	GetUsers(ctx context.Context, input Filter) (model.UserSlice, int64, error)

	// UpdateUser updates the user
	UpdateUser(ctx context.Context, updateUser model.User) (int64, error)

	// GetUser returns a user by input "id" param
	GetUser(ctx context.Context, id int) (model.User, error)

	// DeleteUser deletes the user with the given id
	DeleteUser(ctx context.Context, id int) (int64, error)

	// GetUserByEmail returns a user with the given email
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	// GetStatistics returns summary statistic of users
	GetStatistics(ctx context.Context) (SummaryStatistics, error)
}

type impl struct {
	db *sql.DB
}

func New(db *sql.DB) IUser {
	return impl{
		db: db,
	}
}
