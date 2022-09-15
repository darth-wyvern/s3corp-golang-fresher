package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	orderRepo "github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
)

type OrderStatus string

const (
	OrderStatusNew     OrderStatus = "NEW"
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusSuccess OrderStatus = "SUCCESS"
	OrderStatusFailed  OrderStatus = "FAILED"
)

type OrderItemInput struct {
	ProductID    int
	ProductPrice float64
	ProductName  string
	Quantity     int
	Discount     float64
	Note         string
}

type OrderInput struct {
	Note   string
	UserID int
	Items  []OrderItemInput
}

func (serv impl) CreateOrder(ctx context.Context, input OrderInput) error {
	// Check exists user by order user_id
	existed, err := serv.repo.User().ExistsUserByID(ctx, input.UserID)
	if err != nil {
		return err
	}
	if !existed {
		return ErrUserNotExist
	}

	// Check exists product by order item product_id
	for i, item := range input.Items {
		product, err := serv.repo.Product().GetProduct(ctx, item.ProductID)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProductNotExist
		} else if err != nil {
			return err
		}
		input.Items[i].ProductName = product.Title
		input.Items[i].ProductPrice = product.Price
	}

	if err = serv.repo.Tx(ctx, func(tx *sql.Tx) error {
		// Generate order number
		orderNumber, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("error when generate order number: %v", err)
		}

		// Create order
		order, err := serv.repo.Order().CreateOrder(ctx, tx, model.Order{
			OrderNumber: orderNumber.String(),
			OrderDate:   time.Now(),
			Status:      string(OrderStatusNew),
			Note:        input.Note,
			UserID:      input.UserID,
		})
		if err != nil {
			return fmt.Errorf("error when create order: %v", err)
		}

		// Create order item
		for _, item := range input.Items {
			if err := serv.repo.Order().CreateItem(ctx, tx, model.OrderItem{
				OrderID:      order.ID,
				ProductID:    item.ProductID,
				ProductPrice: item.ProductPrice,
				ProductName:  item.ProductName,
				Quantity:     item.Quantity,
				Discount:     item.Discount,
				Note:         item.Note,
			}); err != nil {
				return fmt.Errorf("error when create order item: %v", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// OrderFilter represents filter conditions for orders
type OrderFilter struct {
	ID          int
	OrderNumber string
	Status      string
	UserID      int
}

// OrderSortBy represents sort conditions for orders
type OrderSortBy struct {
	OrderDate string
	CreatedAt string
}

// Pagination represents pagination conditions for order list
type Pagination struct {
	Limit int
	Page  int
}

// OrdersInput represents all conditions for filter order list
type OrdersInput struct {
	Filter     OrderFilter
	SortBy     OrderSortBy
	Pagination Pagination
}

// OrderItem represents order item which will be returned
type OrderItem struct {
	ID           int
	ProductID    int
	ProductPrice float64
	ProductName  string
	Quantity     int
	Discount     float64
	Note         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Order represents result item which will be returned
type Order struct {
	ID          int
	OrderNumber string
	OrderDate   time.Time
	Status      string
	Note        string
	UserID      int
	OrderItems  []OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetOrders returns list of orders which is filterd
func (serv impl) GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error) {
	orders, totalCount, err := serv.repo.Order().GetOrders(ctx, orderRepo.OrdersInput{
		Filter: orderRepo.OrderFilter{
			ID:          input.Filter.ID,
			OrderNumber: input.Filter.OrderNumber,
			Status:      input.Filter.Status,
			UserID:      input.Filter.UserID,
		},
		SortBy: orderRepo.OrderSortBy{
			OrderDate: input.SortBy.OrderDate,
			CreatedAt: input.SortBy.CreatedAt,
		},
		Pagination: orderRepo.Pagination{
			Limit: input.Pagination.Limit,
			Page:  input.Pagination.Page,
		},
	})
	if err != nil {
		return []Order{}, 0, err
	}

	// Mapping orders
	result := make([]Order, len(orders))

	for i, order := range orders {
		orderItems := make([]OrderItem, len(order.OrderItems))
		for j, orderItem := range order.OrderItems {
			orderItems[j] = OrderItem{
				ID:           orderItem.ID,
				ProductID:    orderItem.ProductID,
				ProductPrice: orderItem.ProductPrice,
				ProductName:  orderItem.ProductName,
				Quantity:     orderItem.Quantity,
				Discount:     orderItem.Discount,
				Note:         orderItem.Note,
				CreatedAt:    orderItem.CreatedAt,
				UpdatedAt:    orderItem.UpdatedAt,
			}
		}

		result[i] = Order{
			ID:          order.ID,
			OrderNumber: order.OrderNumber,
			OrderDate:   order.OrderDate,
			Status:      order.Status,
			Note:        order.Note,
			UserID:      order.UserID,
			OrderItems:  orderItems,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}

	return result, totalCount, nil
}
