package order

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/db"
)

func TestOrderRepository_CreateOrder(t *testing.T) {
	tcs := map[string]struct {
		input     model.Order
		expResult model.Order
		expErr    error
	}{
		"success": {
			input: model.Order{
				OrderNumber: "123456789",
				OrderDate:   time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC),
				Status:      "NEW",
				UserID:      10,
				Note:        "Order 1",
			},
			expResult: model.Order{
				OrderNumber: "123456789",
				OrderDate:   time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC),
				Status:      "NEW",
				UserID:      10,
				Note:        "Order 1",
			},
		},
		"error": {
			input: model.Order{
				OrderNumber: "AAA",
				OrderDate:   time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC),
				Status:      "NEW",
				UserID:      10,
				Note:        "Order 1",
			},
			expErr: errors.New("model: unable to insert into orders: pq: duplicate key value violates unique constraint \"order_number_on_orders\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			orderRepo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/order.sql")
			defer dbTest.Exec("DELETE FROM orders; DELETE FROM users;")

			txTest, err := dbTest.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
			require.NoError(t, err)
			defer txTest.Rollback()

			// When
			result, err := orderRepo.CreateOrder(context.Background(), txTest, tc.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
				tc.expResult.ID = result.ID
				tc.expResult.CreatedAt = result.CreatedAt
				tc.expResult.UpdatedAt = result.UpdatedAt
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestOrderRepository_CreateItem(t *testing.T) {
	tcs := map[string]struct {
		input  model.OrderItem
		expErr error
	}{
		"success": {
			input: model.OrderItem{
				OrderID:      10,
				ProductID:    11,
				ProductPrice: 1000,
				ProductName:  "Product 10",
				Quantity:     1,
				Discount:     0,
				Note:         "Order Item 1",
			},
			expErr: nil,
		},
		"error": {
			input: model.OrderItem{
				OrderID:      10,
				ProductID:    10,
				ProductPrice: 1000,
				ProductName:  "Product 10",
				Quantity:     1,
				Discount:     0,
				Note:         "Order Item 1",
			},
			expErr: errors.New("model: unable to insert into order_items: pq: duplicate key value violates unique constraint \"order_id_product_id_on_order_items\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			orderRepo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/order_item.sql")
			defer dbTest.Exec(`DELETE FROM order_items;DELETE FROM products; DELETE FROM orders; DELETE FROM users;`)

			txTest, err := dbTest.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
			require.NoError(t, err)
			defer txTest.Rollback()

			// When
			err = orderRepo.CreateItem(context.Background(), txTest, tc.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrderRepository_GetStatistics(t *testing.T) {
	tcs := map[string]struct {
		expResult []Statistics
		expErr    error
	}{
		"success": {
			expResult: []Statistics{
				{
					Status: "FAILED",
					Count:  1,
				},
				{
					Status: "SUCCESS",
					Count:  2,
				},
				{
					Status: "PENDING",
					Count:  1,
				}, {
					Status: "NEW",
					Count:  2,
				},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			orderRepo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/order.sql")
			defer dbTest.Exec("DELETE FROM orders; DELETE FROM users;")

			// When
			result, err := orderRepo.GetStatistics(context.Background())

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.expResult, result)
			}
		})
	}
}

func TestOrderRepository_GetLatestOrder(t *testing.T) {
	tcs := map[string]struct {
		expResult []OrderInfo
		expErr    error
	}{
		"success": {
			expResult: []OrderInfo{
				{
					OrderID:     10,
					OrderNumber: "AAA",
					OrderDate:   time.Date(2022, 8, 4, 2, 0, 0, 0, time.UTC),
					Status:      "NEW",
					UserID:      10,
					Total:       10000,
				},
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			orderRepo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/order_item.sql")
			defer dbTest.Exec(`DELETE FROM order_items;DELETE FROM products; DELETE FROM orders; DELETE FROM users;`)

			// When
			result, err := orderRepo.GetLatestOrder(context.Background(), 2)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestOderRepository_GetOrders(t *testing.T) {
	type input struct {
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
						Status: "NEW",
						UserID: 10,
					},
					Pagination: Pagination{
						Page:  1,
						Limit: 10,
					},
				},
			},
			expOutput: output{
				totalCount: 2,
				orders: []Order{
					{
						ID:          1,
						OrderNumber: "ORDER_NUMBER_1",
						Status:      "NEW",
						UserID:      10,
						OrderItems: []OrderItem{
							{
								ID:           10,
								ProductPrice: 1000,
								ProductID:    10,
								ProductName:  "Product 10",
								Quantity:     20,
								Discount:     0,
							},
							{
								ID:           11,
								ProductPrice: 1000,
								ProductID:    11,
								ProductName:  "Product 11",
								Quantity:     30,
								Discount:     0,
							},
						},
					},
					{
						ID:          5,
						OrderNumber: "ORDER_NUMBER_5",
						Status:      "NEW",
						UserID:      10,
						OrderItems:  []OrderItem{},
					},
				},
			},
		},
		"empty_result": {
			input: input{
				ctx: context.Background(),
				ordersInput: OrdersInput{
					Filter: OrderFilter{
						OrderNumber: "ORDER_NUMBER_100",
					},
					Pagination: Pagination{
						Page:  1,
						Limit: 10,
					},
				},
			},
			expOutput: output{
				totalCount: 0,
				orders:     []Order{},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			orderRepo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/get_orders.sql")
			defer dbTest.Exec(`DELETE FROM order_items;DELETE FROM products; DELETE FROM orders; DELETE FROM users;`)

			// When
			result, totalCount, err := orderRepo.GetOrders(tc.input.ctx, tc.input.ordersInput)

			// Then
			if tc.expOutput.err != nil {
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.totalCount, totalCount)
				for i, order := range result {
					require.Equal(t, len(tc.expOutput.orders[i].OrderItems), len(order.OrderItems))
					for j, orderItem := range order.OrderItems {
						tc.expOutput.orders[i].OrderItems[j].CreatedAt = orderItem.CreatedAt
						tc.expOutput.orders[i].OrderItems[j].UpdatedAt = orderItem.UpdatedAt
					}
					tc.expOutput.orders[i].OrderDate = order.OrderDate
					tc.expOutput.orders[i].CreatedAt = order.CreatedAt
					tc.expOutput.orders[i].UpdatedAt = order.UpdatedAt
					require.Equal(t, tc.expOutput.orders[i], order)
				}
			}
		})
	}
}
