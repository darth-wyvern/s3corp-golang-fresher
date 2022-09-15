package product

import (
	"context"
	"database/sql"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type IProduct interface {
	// GetProduct get 1 product by id
	GetProduct(ctx context.Context, id int) (model.Product, error)

	// ExistsProductByID check product exist by id
	ExistsProductByID(ctx context.Context, id int) (bool, error)

	// CreateProduct create new product
	CreateProduct(ctx context.Context, newProduct model.Product) (model.Product, error)

	// InsertAll insert all the given products and return the error
	InsertAll(ctx context.Context, tx *sql.Tx, products []model.Product) error

	// UpdateProduct update product with id
	UpdateProduct(ctx context.Context, product model.Product) (int64, error)

	// DeleteProduct delete a product by product id
	DeleteProduct(ctx context.Context, id int) (int64, error)

	// GetProducts returns list of products (filtered by filter obj)
	GetProducts(ctx context.Context, filter Filter) ([]ProductItem, int64, error)

	// GetStatistics returns summary statistic of products
	GetStatistics(ctx context.Context) (SummaryStatistics, error)
}

type impl struct {
	db *sql.DB
}

func New(db *sql.DB) IProduct {
	return &impl{db: db}
}
