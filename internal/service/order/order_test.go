package order

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
	orderRepo "github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
)

func TestOrderService_CreateOrder(t *testing.T) {
	type mockData struct {
		txFn       mock.AnythingOfTypeArgument
		txErr      error
		userExist  bool
		userErr    error
		product    []model.Product
		productErr []error
	}
	type givenData struct {
		input OrderInput
		mock  mockData
	}
	tcs := map[string]struct {
		given  givenData
		expErr error
	}{
		"success": {
			given: givenData{
				input: OrderInput{
					Note:   "New order",
					UserID: 2,
					Items: []OrderItemInput{
						{
							ProductID: 1,
							Quantity:  10,
							Discount:  0,
							Note:      "item 1",
						},
						{
							ProductID: 3,
							Quantity:  20,
							Discount:  0.5,
							Note:      "item 2",
						},
					},
				},
				mock: mockData{
					txFn:      mock.AnythingOfType("func(*sql.Tx) error"),
					txErr:     nil,
					userExist: true,
					userErr:   nil,
					product: []model.Product{
						{
							ID:    1,
							Title: "product 1",
							Price: 100,
						},
						{
							ID:    2,
							Title: "product 2",
							Price: 200,
						},
					},
					productErr: []error{nil, nil},
				},
			},
			expErr: nil,
		},
		"error_user_is_not_exists": {
			given: givenData{
				input: OrderInput{
					Note:   "New order",
					UserID: 2,
					Items: []OrderItemInput{
						{
							ProductID: 1,
							Quantity:  10,
							Discount:  0,
							Note:      "item 1",
						},
						{
							ProductID: 3,
							Quantity:  20,
							Discount:  0.5,
							Note:      "item 2",
						},
					},
				},
				mock: mockData{
					txFn:      mock.AnythingOfType("func(*sql.Tx) error"),
					txErr:     nil,
					userExist: false,
					userErr:   nil,
					product: []model.Product{
						{
							ID:    1,
							Title: "product 1",
							Price: 100,
						},
						{
							ID:    2,
							Title: "product 2",
							Price: 200,
						},
					},
					productErr: []error{nil, nil},
				},
			},
			expErr: ErrUserNotExist,
		},
		"error_product_is_not_exists": {
			given: givenData{
				input: OrderInput{
					Note:   "New order",
					UserID: 2,
					Items: []OrderItemInput{
						{
							ProductID: 1,
							Quantity:  10,
							Discount:  0,
							Note:      "item 1",
						},
						{
							ProductID: 3,
							Quantity:  20,
							Discount:  0.5,
							Note:      "item 2",
						},
					},
				},
				mock: mockData{
					txFn:      mock.AnythingOfType("func(*sql.Tx) error"),
					txErr:     nil,
					userExist: true,
					userErr:   nil,
					product: []model.Product{
						{
							ID:    1,
							Title: "product 1",
							Price: 100,
						},
						{
							ID:    2,
							Title: "product 2",
							Price: 200,
						},
					},
					productErr: []error{sql.ErrNoRows, nil},
				},
			},
			expErr: ErrProductNotExist,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			repoMock.On("Tx", context.Background(), tc.given.mock.txFn).Return(tc.given.mock.txErr)
			userRepo := new(user.Mock)
			userRepo.On("ExistsUserByID", context.Background(), tc.given.input.UserID).Return(tc.given.mock.userExist, tc.given.mock.userErr)
			repoMock.On("User").Return(userRepo)
			productRepo := new(product.Mock)
			for i, item := range tc.given.input.Items {
				productRepo.On("GetProduct", context.Background(), item.ProductID).Return(tc.given.mock.product[i], tc.given.mock.productErr[i])
			}
			repoMock.On("Product").Return(productRepo)
			orderRepoMock := new(order.Mock)
			repoMock.On("Order").Return(orderRepoMock)

			orderServ := New(repoMock)

			// When
			err := orderServ.CreateOrder(context.Background(), tc.given.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrderService_GetOrders(t *testing.T) {
	type mockData struct {
		inputCTX         context.Context
		input            orderRepo.OrdersInput
		mockResultTotal  int64
		mockResultOrders []orderRepo.Order
		mockResultError  error
	}
	type input struct {
		mockData    mockData
		ctx         context.Context
		ordersInput OrdersInput
	}
	type output struct {
		totalCount int64
		orders     []Order
		err        error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx: context.Background(),
				ordersInput: OrdersInput{
					Filter: OrderFilter{
						ID:          1,
						OrderNumber: "abcd",
						Status:      "NEW",
						UserID:      2,
					},
					SortBy: OrderSortBy{
						OrderDate: "asc",
						CreatedAt: "asc",
					},
					Pagination: Pagination{
						Limit: 2,
						Page:  10,
					},
				},
				mockData: mockData{
					inputCTX: context.Background(),
					input: orderRepo.OrdersInput{
						Filter: orderRepo.OrderFilter{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							UserID:      2,
						},
						SortBy: orderRepo.OrderSortBy{
							OrderDate: "asc",
							CreatedAt: "asc",
						},
						Pagination: orderRepo.Pagination{
							Limit: 2,
							Page:  10,
						},
					},
					mockResultTotal: 2,
					mockResultOrders: []orderRepo.Order{
						{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							UserID:      2,
							OrderItems: []orderRepo.OrderItem{
								{
									ID:           2,
									ProductPrice: 20000,
									ProductID:    10,
									ProductName:  "Product 1",
									Quantity:     10,
									Discount:     0,
								},
							},
						},
						{
							ID:          2,
							OrderNumber: "abcdef",
							Status:      "NEW",
							UserID:      2,
							OrderItems: []orderRepo.OrderItem{
								{
									ID:           3,
									ProductPrice: 20000,
									ProductID:    11,
									ProductName:  "Product 2",
									Quantity:     20,
									Discount:     0,
								},
							},
						},
					},
				},
			},
			expOutput: output{
				totalCount: 2,
				orders: []Order{
					{
						ID:          1,
						OrderNumber: "abcd",
						Status:      "NEW",
						UserID:      2,
						OrderItems: []OrderItem{
							{
								ID:           2,
								ProductPrice: 20000,
								ProductID:    10,
								ProductName:  "Product 1",
								Quantity:     10,
								Discount:     0,
							},
						},
					},
					{
						ID:          2,
						OrderNumber: "abcdef",
						Status:      "NEW",
						UserID:      2,
						OrderItems: []OrderItem{
							{
								ID:           3,
								ProductPrice: 20000,
								ProductID:    11,
								ProductName:  "Product 2",
								Quantity:     20,
								Discount:     0,
							},
						},
					},
				},
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			orderRepoMock := new(order.Mock)
			orderRepoMock.On("GetOrders", tc.input.mockData.inputCTX, tc.input.mockData.input).Return(tc.input.mockData.mockResultOrders, tc.input.mockData.mockResultTotal, tc.input.mockData.mockResultError)

			repoMock.On("Order").Return(orderRepoMock)

			orderServ := New(repoMock)

			// When
			result, totalCount, err := orderServ.GetOrders(tc.input.ctx, tc.input.ordersInput)

			// THhen
			if err != nil {
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.totalCount, totalCount)
				require.Equal(t, tc.expOutput.orders, result)
			}
		})
	}
}
