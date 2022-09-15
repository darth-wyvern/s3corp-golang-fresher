package v1

import (
	"net/http"

	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

const (
	StatisticsOrderLimit = 10
)

type UserSummary struct {
	Total         int64 `json:"total"`
	TotalInactive int64 `json:"total_inactive"`
}

type ProductSummary struct {
	Total         int64 `json:"total"`
	TotalInactive int64 `json:"total_inactive"`
}

type OrderSummary struct {
	TotalNew     int64 `json:"total_new"`
	TotalPending int64 `json:"total_pending"`
	TotalSuccess int64 `json:"total_success"`
	TotalFailed  int64 `json:"total_failed"`
}

type OrderInfo struct {
	OrderID     int     `json:"order_id"`
	OrderNumber string  `json:"order_number"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	UserID      int     `json:"user_id"`
	Total       float64 `json:"total"`
}

type StatisticsResponse struct {
	Users       UserSummary    `json:"users"`
	Products    ProductSummary `json:"products"`
	Orders      OrderSummary   `json:"orders"`
	LatestOrder []OrderInfo    `json:"latest_orders"`
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	summary, err := h.userServ.GetStatistics(r.Context(), StatisticsOrderLimit)
	if err != nil {
		handleStatisticError(w, err)
		return
	}

	latestOrders := make([]OrderInfo, len(summary.LatestOrders))
	for i, order := range summary.LatestOrders {
		latestOrders[i] = OrderInfo{
			OrderID:     order.OrderID,
			OrderNumber: order.OrderNumber,
			OrderDate:   order.OrderDate.String(),
			Status:      order.Status,
			UserID:      order.UserID,
			Total:       order.Total,
		}
	}

	utils.WriteJSONResponse(w, http.StatusOK, StatisticsResponse{
		Users: UserSummary{
			Total:         summary.Users.Total,
			TotalInactive: summary.Users.TotalInactive,
		},
		Products: ProductSummary{
			Total:         summary.Products.Total,
			TotalInactive: summary.Products.TotalInactive,
		},
		Orders: OrderSummary{
			TotalNew:     summary.Orders.TotalNew,
			TotalPending: summary.Orders.TotalPending,
			TotalSuccess: summary.Orders.TotalSuccess,
			TotalFailed:  summary.Orders.TotalFailed,
		},
		LatestOrder: latestOrders,
	})
}
