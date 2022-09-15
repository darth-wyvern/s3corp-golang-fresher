package product

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

// Mock repo implement the repo interface, inherit Mock struct from testify lib
type Mock struct {
	mock.Mock
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

func (p *Mock) ExistsProductByID(ctx context.Context, id int) (bool, error) {
	args := p.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (p *Mock) UpdateProduct(ctx context.Context, product model.Product) (int64, error) {
	args := p.Called(ctx, product)
	return args.Get(0).(int64), args.Error(1)
}

func (p *Mock) CreateProduct(ctx context.Context, newProduct model.Product) (model.Product, error) {

	args := p.Called(ctx, newProduct)

	var (
		product model.Product
		err     error
	)
	product = args.Get(0).(model.Product)
	err = args.Error(1)

	return product, err
}
func (p *Mock) DeleteProduct(ctx context.Context, id int) (int64, error) {
	args := p.Called(ctx, id)
	return int64(args.Int(0)), args.Error(1)
}
func (p *Mock) GetProducts(ctx context.Context, filter Filter) ([]ProductItem, int64, error) {
	args := p.Called(ctx, filter)
	return args.Get(0).([]ProductItem), args.Get(1).(int64), args.Error(2)
}

func (p *Mock) InsertAll(ctx context.Context, tx *sql.Tx, products []model.Product) error {
	args := p.Called(ctx, tx, products)
	return args.Error(0)
}

func (m *Mock) GetStatistics(ctx context.Context) (SummaryStatistics, error) {
	args := m.Called(ctx)
	return args.Get(0).(SummaryStatistics), args.Error(1)
}
