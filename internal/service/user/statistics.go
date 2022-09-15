package user

import (
	"context"
	"time"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
)

type UserSummary struct {
	Total         int64
	TotalInactive int64
}

type ProductSummary struct {
	Total         int64
	TotalInactive int64
}

type OrderSummary struct {
	TotalNew     int64
	TotalPending int64
	TotalSuccess int64
	TotalFailed  int64
}

type OrderInfo struct {
	OrderID     int
	OrderNumber string
	OrderDate   time.Time
	Status      string
	UserID      int
	Total       float64
}

type SummaryStatistics struct {
	Users        UserSummary
	Products     ProductSummary
	Orders       OrderSummary
	LatestOrders []OrderInfo
}

func (serv impl) GetStatistics(ctx context.Context, orderLimit int) (SummaryStatistics, error) {
	userSummary, err := serv.repo.User().GetStatistics(ctx)
	if err != nil {
		return SummaryStatistics{}, err
	}

	productSummary, err := serv.repo.Product().GetStatistics(ctx)
	if err != nil {
		return SummaryStatistics{}, err
	}

	orderStatistics, err := serv.repo.Order().GetStatistics(ctx)
	if err != nil {
		return SummaryStatistics{}, err
	}
	var orderSummary OrderSummary
	for _, stat := range orderStatistics {
		switch stat.Status {
		case string(order.OrderStatusNew):
			orderSummary.TotalNew = stat.Count
		case string(order.OrderStatusPending):
			orderSummary.TotalPending = stat.Count
		case string(order.OrderStatusSuccess):
			orderSummary.TotalSuccess = stat.Count
		case string(order.OrderStatusFailed):
			orderSummary.TotalFailed = stat.Count
		}
	}

	orderList, err := serv.repo.Order().GetLatestOrder(ctx, orderLimit)
	if err != nil {
		return SummaryStatistics{}, err
	}

	latestOrders := make([]OrderInfo, len(orderList))
	for i, order := range orderList {
		latestOrders[i] = OrderInfo{
			OrderID:     order.OrderID,
			OrderNumber: order.OrderNumber,
			OrderDate:   order.OrderDate,
			Status:      order.Status,
			UserID:      order.UserID,
			Total:       order.Total,
		}
	}

	return SummaryStatistics{
		Users: UserSummary{
			Total:         userSummary.Total,
			TotalInactive: userSummary.TotalInactive,
		},
		Products: ProductSummary{
			Total:         productSummary.Total,
			TotalInactive: productSummary.TotalInactive,
		},
		Orders:       orderSummary,
		LatestOrders: latestOrders,
	}, nil
}
