package graph

import (
	"context"
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/handler/gql/graph/mod"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"

	productServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
)

func TestProductResolver_CreateProduct(t *testing.T) {
	type input struct {
		ctx              context.Context
		productInput     mod.CreateProductInput
		mockInputProduct productServ.ProductInput
		mockInputCTX     context.Context
		mockResult       model.Product
		mockError        error
	}
	type output struct {
		product mod.Product
		err     error
	}

	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "test product",
					Price:    100000,
					Quantity: 100,
					IsActive: mod.ActiveTypeNo,
					UserID:   2,
				},
				mockInputProduct: productServ.ProductInput{
					Title:    "test product",
					Price:    100000,
					Quantity: 100,
					IsActive: false,
					UserID:   2,
				},
				mockInputCTX: context.Background(),
				mockResult: model.Product{
					Title:    "test product",
					Price:    100000,
					Quantity: 100,
					IsActive: false,
					UserID:   2,
				},
			},
			expOutput: output{
				product: mod.Product{
					Title:    "test product",
					Price:    100000,
					Quantity: 100,
					IsActive: false,
					UserID:   2,
				},
			},
		},
		"title_is_required": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "",
					Price:    100000,
					Quantity: 100,
					IsActive: mod.ActiveTypeYes,
					UserID:   2,
				},
			},
			expOutput: output{
				err: errTitleCannotBeBlank,
			},
		},
		"invalid_user_id": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "test",
					Price:    100000,
					Quantity: 100,
					IsActive: mod.ActiveTypeYes,
					UserID:   -2,
				},
			},
			expOutput: output{
				err: errInvalidUserID,
			},
		},
		"invalid_price": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "test",
					Price:    -100000,
					Quantity: 100,
					IsActive: mod.ActiveTypeYes,
					UserID:   2,
				},
			},
			expOutput: output{
				err: errInvalidPrice,
			},
		},
		"invalid_quantity": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "test",
					Price:    100000,
					Quantity: -10,
					IsActive: mod.ActiveTypeYes,
					UserID:   2,
				},
			},
			expOutput: output{
				err: errInvalidQuantity,
			},
		},
		"user_is_not_exist": {
			input: input{
				ctx: context.Background(),
				productInput: mod.CreateProductInput{
					Title:    "test",
					Price:    100000,
					Quantity: 100,
					IsActive: mod.ActiveTypeYes,
					UserID:   2,
				},
				mockInputProduct: productServ.ProductInput{
					Title:    "test",
					Price:    100000,
					Quantity: 100,
					IsActive: true,
					UserID:   2,
				},
				mockInputCTX: context.Background(),
				mockError:    productServ.ErrUserNotExist,
			},
			expOutput: output{
				err: errUserNotExist,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			productServiceMock := new(productServ.Mock)
			productServiceMock.On("CreateProduct", tc.input.mockInputCTX, tc.input.mockInputProduct).Return(tc.input.mockResult, tc.input.mockError)

			resolver := NewResolver(nil, productServiceMock)

			// WHEN
			result, err := resolver.Mutation().CreateProduct(tc.input.ctx, tc.input.productInput)

			// THEN
			if tc.expOutput.err != nil {
				// must be error
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				// must be success
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.product, *result)
			}
		})
	}
}

func TestProductResolver_GetProducts(t *testing.T) {
	type mockData struct {
		input      productServ.GetProductsInput
		products   []productServ.ProductItem
		err        error
		totalCount int64
	}

	type givenData struct {
		input mod.GetProductsInput
		mock  mockData
	}

	tcs := map[string]struct {
		given     givenData
		expResult mod.GetProductsOutput
		expErr    error
	}{
		"success": {
			given: givenData{
				input: mod.GetProductsInput{
					Title: stringToPtr("test"),
					PriceRange: &mod.PriceRange{
						MinPrice: 100,
						MaxPrice: 3000,
					},
					IsActive: boolToPtr(true),
					UserID:   intToPtr(1),
					OrderBy: &mod.OrderBy{
						Title:     stringToPtr("desc"),
						CreatedAt: stringToPtr("desc"),
					},
				},
				mock: mockData{
					input: productServ.GetProductsInput{
						Title:      "test",
						PriceRange: productServ.PriceRange{MinPrice: 100, MaxPrice: 3000},
						IsActive:   null.NewBool(true, true),
						UserID:     1,
						OrderBy: productServ.OrderInput{
							Title:     "desc",
							CreatedAt: "desc",
						},
						Pagination: productServ.Pagination{
							Limit: 20,
							Page:  1,
						},
					},
					products: []productServ.ProductItem{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
						{
							ID:          2,
							Title:       "test 2",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					totalCount: 2,
				},
			},
			expResult: mod.GetProductsOutput{
				Products: []*mod.Product{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
					{
						ID:          2,
						Title:       "test 2",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
				},
				Pagination: &mod.Pagination{
					CurrentPage: intToPtr(1),
					Limit:       intToPtr(20),
					TotalCount:  int64ToPtr(2),
				},
			},
		},
		"invalid_id": {
			given: givenData{
				input: mod.GetProductsInput{
					ID: intToPtr(-1),
				},
			},
			expErr: errInvalidID,
		},
		"invalid_price": {
			given: givenData{
				input: mod.GetProductsInput{
					PriceRange: &mod.PriceRange{
						MinPrice: -1,
						MaxPrice: 2000,
					},
				},
			},
			expErr: errInvalidPriceRange,
		},
		"invalid_user_id": {
			given: givenData{
				input: mod.GetProductsInput{
					UserID: intToPtr(-1),
				},
			},
			expErr: errInvalidUserID,
		},
		"invalid_order_by": {
			given: givenData{
				input: mod.GetProductsInput{
					OrderBy: &mod.OrderBy{
						Title: stringToPtr("invalid"),
					},
				},
			},
			expErr: errInvalidOrderBy,
		},
		"empty_body": {
			given: givenData{
				mock: mockData{
					input: productServ.GetProductsInput{
						Pagination: productServ.Pagination{
							Limit: 20,
							Page:  1,
						},
					},
					products: []productServ.ProductItem{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
						{
							ID:          2,
							Title:       "test 2",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					totalCount: 2,
				},
			},
			expResult: mod.GetProductsOutput{
				Products: []*mod.Product{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
					{
						ID:          2,
						Title:       "test 2",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
				},
				Pagination: &mod.Pagination{
					CurrentPage: intToPtr(1),
					Limit:       intToPtr(20),
					TotalCount:  int64ToPtr(2),
				},
			},
		},
		"pagination:": {
			given: givenData{
				input: mod.GetProductsInput{
					Title: stringToPtr("test"),
					PriceRange: &mod.PriceRange{
						MinPrice: 100,
						MaxPrice: 3000,
					},
					IsActive: boolToPtr(true),
					UserID:   intToPtr(1),
					OrderBy: &mod.OrderBy{
						Title:     stringToPtr("desc"),
						CreatedAt: stringToPtr("desc"),
					},
					Pagination: &mod.PaginationInput{
						Page:  intToPtr(2),
						Limit: intToPtr(20),
					},
				},
				mock: mockData{
					input: productServ.GetProductsInput{
						Title:      "test",
						PriceRange: productServ.PriceRange{MinPrice: 100, MaxPrice: 3000},
						IsActive:   null.NewBool(true, true),
						UserID:     1,
						OrderBy: productServ.OrderInput{
							Title:     "desc",
							CreatedAt: "desc",
						},
						Pagination: productServ.Pagination{
							Limit: 20,
							Page:  2,
						},
					},
					products: []productServ.ProductItem{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
						{
							ID:          2,
							Title:       "test 2",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: productServ.CreatedBy{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					totalCount: 2,
				},
			},
			expResult: mod.GetProductsOutput{
				Products: []*mod.Product{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
					{
						ID:          2,
						Title:       "test 2",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						UserID:      1,
						IsActive:    true,
					},
				},
				Pagination: &mod.Pagination{
					CurrentPage: intToPtr(2),
					Limit:       intToPtr(20),
					TotalCount:  int64ToPtr(2),
				},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			serviceMock := new(productServ.Mock)
			serviceMock.On("GetProducts", context.Background(), tc.given.mock.input).Return(tc.given.mock.products, tc.given.mock.totalCount, tc.given.mock.err)
			resolver := NewResolver(nil, serviceMock)

			// WHEN
			result, err := resolver.Query().GetProducts(context.Background(), tc.given.input)

			// THEN
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				for i := 0; i < len(result.Products); i++ {
					require.Equal(t, tc.expResult.Products[i], result.Products[i])
				}
				require.Equal(t, *tc.expResult.Pagination.CurrentPage, *result.Pagination.CurrentPage)
				require.Equal(t, *tc.expResult.Pagination.Limit, *result.Pagination.Limit)
				require.Equal(t, *tc.expResult.Pagination.TotalCount, *result.Pagination.TotalCount)
			}
		})
	}
}
