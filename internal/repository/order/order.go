package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

func (r impl) CreateOrder(ctx context.Context, tx *sql.Tx, order model.Order) (model.Order, error) {
	if err := order.Insert(context.Background(), tx, boil.Infer()); err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (r impl) CreateItem(ctx context.Context, tx *sql.Tx, item model.OrderItem) error {
	return item.Insert(context.Background(), tx, boil.Infer())
}

type Statistics struct {
	Status string `boil:"status"`
	Count  int64  `boil:"count"`
}

func (r impl) GetStatistics(ctx context.Context) ([]Statistics, error) {
	qms := []qm.QueryMod{
		qm.Select("count(*) as count", model.OrderColumns.Status),
		qm.From("orders"),
		qm.GroupBy(model.OrderColumns.Status),
	}

	var result []Statistics
	if err := model.NewQuery(qms...).Bind(ctx, r.db, &result); err != nil {
		return result, err
	}

	return result, nil
}

type OrderInfo struct {
	OrderID     int       `boil:"order_id"`
	OrderNumber string    `boil:"order_number"`
	OrderDate   time.Time `boil:"order_date"`
	Status      string    `boil:"status"`
	UserID      int       `boil:"user_id"`
	Total       float64   `boil:"total"`
}

func (r impl) GetLatestOrder(ctx context.Context, limit int) ([]OrderInfo, error) {
	qms := []qm.QueryMod{
		qm.Select(
			"o.id as order_id",
			model.OrderColumns.OrderNumber,
			model.OrderColumns.OrderDate,
			model.OrderColumns.Status,
			model.OrderColumns.UserID,
			"sum(oi.product_price * oi.quantity)*(1 - oi.discount) as total",
		),
		qm.From("orders o"),
		qm.InnerJoin("order_items oi on oi.order_id = o.id"),
		qm.GroupBy("o.id, oi.discount"),
		qm.Limit(limit),
	}

	var result []OrderInfo
	if err := model.NewQuery(qms...).Bind(context.Background(), r.db, &result); err != nil {
		return nil, err
	}

	return result, nil
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

// GetOrders represents all conditions for filter order list
type OrdersInput struct {
	Filter     OrderFilter
	SortBy     OrderSortBy
	Pagination Pagination
}

// OrderItem represents order item for return in order list
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

// Order represents order for order list which will be returned
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

// GetOrders returns order list from database
func (r impl) GetOrders(ctx context.Context, input OrdersInput) ([]Order, int64, error) {
	var qms = []qm.QueryMod{
		qm.Load(model.OrderRels.OrderItems),
	}

	// Add filter condition.
	if input.Filter.ID > 0 {
		qms = append(qms, model.OrderWhere.ID.EQ(input.Filter.ID))
	}
	if input.Filter.OrderNumber != "" {
		qms = append(qms, model.OrderWhere.OrderNumber.EQ((input.Filter.OrderNumber)))
	}
	if input.Filter.UserID > 0 {
		qms = append(qms, model.OrderWhere.UserID.EQ(input.Filter.UserID))
	}
	if input.Filter.Status != "" {
		qms = append(qms, model.OrderWhere.Status.EQ(input.Filter.Status))
	}

	totalCount, err := model.Orders(qms...).Count(context.Background(), r.db)
	if err != nil {
		return []Order{}, 0, err
	}

	// SortBy
	if input.SortBy != (OrderSortBy{}) {
		if input.SortBy.OrderDate != "" {
			qms = append(qms, qm.OrderBy(model.OrderColumns.OrderDate+" "+input.SortBy.OrderDate))
		}
		if input.SortBy.CreatedAt != "" {
			qms = append(qms, qm.OrderBy(model.OrderColumns.CreatedAt+" "+input.SortBy.CreatedAt))
		}
	} else {
		qms = append(qms, qm.OrderBy(model.OrderColumns.UpdatedAt+" desc"))
	}

	// Paging
	if input.Pagination.Limit > 0 && input.Pagination.Page > 0 {
		qms = append(
			qms,
			qm.Offset(input.Pagination.Limit*(input.Pagination.Page-1)),
			qm.Limit(input.Pagination.Limit))
	}

	orders, err := model.Orders(qms...).All(context.Background(), r.db)
	if err != nil {
		return []Order{}, 0, nil
	}

	var result = make([]Order, len(orders))
	// Map the data
	for i, order := range orders {
		// TODO: Check value of order.R
		orderItems := make([]OrderItem, len(order.R.OrderItems))
		for j, orderItem := range order.R.OrderItems {
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
