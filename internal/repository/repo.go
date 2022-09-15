package repository

import (
	"context"
	"database/sql"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
)

type IRepo interface {

	// User returns user repository
	User() user.IUser

	// Product returns product repository
	Product() product.IProduct

	// Order returns order repository
	Order() order.IOrder

	// Tx commits the given function in a transaction.
	Tx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

func New(db *sql.DB) IRepo {
	return impl{
		db:      db,
		order:   order.New(db),
		user:    user.New(db),
		product: product.New(db),
	}
}

type impl struct {
	db      *sql.DB
	order   order.IOrder
	user    user.IUser
	product product.IProduct
}

func (i impl) User() user.IUser {
	return i.user
}

func (i impl) Product() product.IProduct {
	return i.product
}

func (i impl) Order() order.IOrder {
	return i.order
}

func (i impl) Tx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := i.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
