package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
)

func TestStatisticsHandler_GetStatistics(t *testing.T) {
	type expectedData struct {
		statusCode int
		result     StatisticsResponse
	}

	type mockData struct {
		summary user.SummaryStatistics
		err     error
	}

	tcs := map[string]struct {
		mock   mockData
		exp    expectedData
		expErr error
	}{
		"success": {
			mock: mockData{
				summary: user.SummaryStatistics{
					Users: user.UserSummary{
						Total:         2,
						TotalInactive: 1,
					},
					Products: user.ProductSummary{
						Total:         2,
						TotalInactive: 0,
					},
					Orders: user.OrderSummary{
						TotalNew:     2,
						TotalPending: 1,
						TotalSuccess: 1,
						TotalFailed:  0,
					},
					LatestOrders: []user.OrderInfo{
						{
							OrderID:     1,
							OrderNumber: "TEST",
							OrderDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							Status:      string(order.OrderStatusNew),
							UserID:      10,
							Total:       15000,
						},
					},
				},
			},
			exp: expectedData{
				statusCode: http.StatusOK,
				result: StatisticsResponse{
					Users: UserSummary{
						Total:         2,
						TotalInactive: 1,
					},
					Products: ProductSummary{
						Total:         2,
						TotalInactive: 0,
					},
					Orders: OrderSummary{
						TotalNew:     2,
						TotalPending: 1,
						TotalSuccess: 1,
						TotalFailed:  0,
					},
					LatestOrder: []OrderInfo{
						{
							OrderID:     1,
							OrderNumber: "TEST",
							OrderDate:   "2020-01-01 00:00:00 +0000 UTC",
							Status:      string(order.OrderStatusNew),
							UserID:      10,
							Total:       15000,
						},
					},
				},
			},
		},
		"error": {
			mock: mockData{
				summary: user.SummaryStatistics{},
				err:     fmt.Errorf("UNKNOWN_ERROR"),
			},
			exp: expectedData{
				statusCode: http.StatusInternalServerError,
			},
			expErr: ErrInternalServerError,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			r := httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
			w := httptest.NewRecorder()

			userServiceMock := new(user.Mock)
			userServiceMock.On("GetStatistics", r.Context(), 10).Return(tc.mock.summary, tc.mock.err)

			handler := NewHandler(userServiceMock, nil, nil)

			// When
			handler.GetStatistics(w, r)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.exp.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				require.Equal(t, tc.exp.statusCode, w.Code)

				var actualResult StatisticsResponse
				if err := json.NewDecoder(w.Body).Decode(&actualResult); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.exp.result, actualResult)
			}
		})
	}
}
