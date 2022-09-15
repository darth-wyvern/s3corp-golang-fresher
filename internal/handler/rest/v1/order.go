package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	orderServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

type OrderItemRequest struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
	Note      string  `json:"note"`
}

type OrderRequest struct {
	Note   string             `json:"note"`
	UserID int                `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

func validateOrderRequest(r OrderRequest) (orderServ.OrderInput, error) {
	if r.UserID <= 0 {
		return orderServ.OrderInput{}, ErrInvalidUserID
	}
	if len(r.Items) == 0 {
		return order.OrderInput{}, ErrItemsCannotBeBlank
	}
	items := make([]orderServ.OrderItemInput, len(r.Items))
	for i, item := range r.Items {
		if item.ProductID <= 0 {
			return orderServ.OrderInput{}, ErrInvalidProductID
		}
		if item.Quantity <= 0 {
			return orderServ.OrderInput{}, ErrInvalidQuantity
		}
		if item.Discount < 0 || item.Discount > 1 {
			return orderServ.OrderInput{}, ErrInvalidDiscount
		}
		items[i] = orderServ.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Discount:  item.Discount,
			Note:      item.Note,
		}
	}
	return orderServ.OrderInput{
		Note:   r.Note,
		UserID: r.UserID,
		Items:  items,
	}, nil
}

// handleOrderError handles error from order service
func handleOrderError(w http.ResponseWriter, err error) {
	switch err {
	case order.ErrUserNotExist:
		utils.WriteJSONResponse(w, http.StatusBadRequest, ErrUserNotExist)
	case order.ErrProductNotExist:
		utils.WriteJSONResponse(w, http.StatusBadRequest, ErrProductNotFound)
	default:
		utils.WriteJSONResponse(w, http.StatusInternalServerError, ErrInternalServerError)
	}
}

// CreateOrder handles request to create order
func (h Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var orderRequest OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, ErrInvalidBodyRequest)
		return
	}

	// Validate request body
	input, err := validateOrderRequest(orderRequest)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	// Create order
	if err = h.orderServ.CreateOrder(r.Context(), input); err != nil {
		handleOrderError(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, "Created order successfully")
}

// orderFilter represents filter conditions for orders
type orderFilter struct {
	ID          int    `json:"id"`
	OrderNumber string `json:"order_number"`
	Status      string `json:"status"`
	UserID      int    `json:"user_id"`
}

// orderSortBy represents sort conditions for orders
type orderSortBy struct {
	OrderDate string `json:"order_date"`
	CreatedAt string `json:"created_at"`
}

// ordersRequest represents all conditions for filter order list
type ordersRequest struct {
	Filter     orderFilter     `json:"filter"`
	SortBy     orderSortBy     `json:"sort_by"`
	Pagination paginationInput `json:"pagination"`
}

// validOrderStatus validate order status value
func validOrderStatus(status string) bool {
	if status != string(orderServ.OrderStatusFailed) &&
		status != string(orderServ.OrderStatusNew) &&
		status != string(orderServ.OrderStatusPending) &&
		status != string(orderServ.OrderStatusSuccess) {
		return false
	}
	return true
}

// validGetOrdersReq validate request body of get orders handler
func validGetOrdersReq(req ordersRequest) (orderServ.OrdersInput, error) {
	if req.Filter.ID < 0 {
		return orderServ.OrdersInput{}, ErrInvalidOrderID
	}

	status := strings.TrimSpace(req.Filter.Status)
	if status != "" && !validOrderStatus(status) {
		return orderServ.OrdersInput{}, ErrInvalidOrderStatus
	}

	if req.Filter.UserID < 0 {
		return orderServ.OrdersInput{}, ErrInvalidUserID
	}

	pageArgs, err := validatePagination(req.Pagination)
	if err != nil {
		return orderServ.OrdersInput{}, err
	}

	orderByCreatedAt := strings.TrimSpace(req.SortBy.CreatedAt)
	if orderByCreatedAt != "" && orderByCreatedAt != OrderTypeASC && orderByCreatedAt != OrderTypeDESC {
		return orderServ.OrdersInput{}, ErrInvalidOrderBy
	}
	orderByOrderDate := strings.TrimSpace(req.SortBy.OrderDate)
	if orderByOrderDate != "" && orderByOrderDate != OrderTypeASC && orderByOrderDate != OrderTypeDESC {
		return orderServ.OrdersInput{}, ErrInvalidOrderBy
	}

	return orderServ.OrdersInput{
		Filter: orderServ.OrderFilter{
			ID:          req.Filter.ID,
			OrderNumber: strings.TrimSpace(req.Filter.OrderNumber),
			Status:      status,
			UserID:      req.Filter.UserID,
		},
		SortBy: orderServ.OrderSortBy{
			OrderDate: orderByOrderDate,
			CreatedAt: orderByCreatedAt,
		},
		Pagination: orderServ.Pagination{
			Limit: pageArgs.Limit,
			Page:  pageArgs.Page,
		},
	}, nil
}

// OrderItem represents order item for response in order list
type OrderItem struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"product_id"`
	ProductPrice float64   `json:"product_price"`
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	Discount     float64   `json:"discount"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Order represents order for order list response
type Order struct {
	ID          int         `json:"id"`
	OrderNumber string      `json:"order_number"`
	OrderDate   time.Time   `json:"order_date"`
	Status      string      `json:"status"`
	Note        string      `json:"note"`
	UserID      int         `json:"user_id"`
	OrderItems  []OrderItem `json:"order_items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderResponse represents order for order list response
type OrdersResponse struct {
	Orders     []Order    `json:"orders"`
	Pagination pagination `json:"pagination"`
}

func (h Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var orderReq ordersRequest
	if r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
			utils.WriteJSONResponse(w, ErrInvalidBodyRequest.Status, ErrInvalidBodyRequest)
			return
		}
	}

	// Validate request body
	ordersInput, err := validGetOrdersReq(orderReq)
	if err != nil {
		utils.WriteJSONResponse(w, err.(utils.ErrorResponse).Status, err.(utils.ErrorResponse))
		return
	}

	// Get orderList from service
	orders, totalCount, err := h.orderServ.GetOrders(r.Context(), ordersInput)
	if err != nil {
		handleOrderError(w, err)
		return
	}

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

	utils.WriteJSONResponse(w, http.StatusOK,
		OrdersResponse{
			Orders: result,
			Pagination: pagination{
				CurrentPage: ordersInput.Pagination.Page,
				Limit:       ordersInput.Pagination.Limit,
				TotalCount:  totalCount,
			},
		})
}
