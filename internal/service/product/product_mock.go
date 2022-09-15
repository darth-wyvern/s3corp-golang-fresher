package product

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

type Mock struct {
	mock.Mock
}

func (p *Mock) UpdateProduct(ctx context.Context, id int, input ProductInput) error {
	args := p.Called(ctx, id, input)
	return args.Error(0)
}

func (p *Mock) GetProduct(ctx context.Context, id int) (model.Product, error) {

	args := p.Called(ctx, id)

	var (
		product model.Product
		err     error
	)
	product = args.Get(0).(model.Product)
	err = args.Error(1)

	return product, err
}

func (p *Mock) CreateProduct(ctx context.Context, product ProductInput) (model.Product, error) {

	args := p.Called(ctx, product)

	return args.Get(0).(model.Product), args.Error(1)
}
func (p *Mock) DeleteProduct(ctx context.Context, id int) error {
	args := p.Called(ctx, id)
	return args.Error(0)
}

func (p *Mock) GetProducts(ctx context.Context, input GetProductsInput) ([]ProductItem, int64, error) {
	args := p.Called(ctx, input)
	return args.Get(0).([]ProductItem), args.Get(1).(int64), args.Error(2)
}

func (p *Mock) ImportProductCSV(ctx context.Context, fileName string, csvFile io.Reader) error {
	args := p.Called(ctx, fileName, csvFile)
	return args.Error(0)
}

func (p *Mock) ExportProductsCSV(ctx context.Context, input GetProductsInput, dir string) (string, error) {
	args := p.Called(ctx, input, dir)
	return args.String(0), args.Error(1)
}

func (p *Mock) DownloadCSV(ctx context.Context, filePath string) ([]byte, error) {
	args := p.Called(ctx, filePath)
	return args.Get(0).([]byte), args.Error(1)
}

func (p *Mock) GetStatistics(ctx context.Context) (SummaryStatistics, error) {
	args := p.Called(ctx)
	return args.Get(0).(SummaryStatistics), args.Error(1)
}
