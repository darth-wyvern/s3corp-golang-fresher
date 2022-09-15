package product

import (
	"context"
	"io"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
)

type IProduct interface {
	// GetProduct get id and return product
	GetProduct(ctx context.Context, id int) (model.Product, error)

	// CreateProduct create new product from product input
	CreateProduct(ctx context.Context, newProduct ProductInput) (model.Product, error)

	// ImportProductCSV insert all data from csvFile
	ImportProductCSV(ctx context.Context, fileName string, csvFile io.Reader) error

	// ExportProductsCSV returns name of the file which is exported
	ExportProductsCSV(ctx context.Context, input GetProductsInput, dir string) (string, error)

	// DownloadCSV return file based on file path
	DownloadCSV(ctx context.Context, filePath string) ([]byte, error)

	// UpdateProduct update product with id and product input
	UpdateProduct(ctx context.Context, id int, product ProductInput) error

	// DeleteProduct get id and delete product
	DeleteProduct(ctx context.Context, id int) error

	// GetProducts returns list of products
	GetProducts(ctx context.Context, input GetProductsInput) ([]ProductItem, int64, error)
}

type impl struct {
	repo repository.IRepo
}

// New create new service with repo parameter(dependency)
func New(repo repository.IRepo) IProduct {
	return &impl{repo: repo}
}
