package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	orderServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
)

func TestHandler_CreateOrder(t *testing.T) {
	type mockData struct {
		input order.OrderInput
		err   error
	}
	type givenData struct {
		reqBody string
		mock    mockData
	}
	type expectedData struct {
		statusCode int
		body       string
	}
	tcs := map[string]struct {
		given        givenData
		expResult    expectedData
		expErr       error
		isCallToServ bool
	}{
		"success": {
			given: givenData{
				reqBody: `{
					"note": "New order",
					"user_id": 2,
					"items": [
						{
							"product_id": 1,
							"quantity": 10,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
				mock: mockData{
					input: order.OrderInput{
						Note:   "New order",
						UserID: 2,
						Items: []order.OrderItemInput{
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
				},
			},
			expResult: expectedData{
				statusCode: http.StatusCreated,
				body:       "Created order successfully",
			},
			expErr:       nil,
			isCallToServ: true,
		},
		"error_invalid_user_id": {
			given: givenData{
				reqBody: `{
					"user_id": -1,
					"items": [
						{
							"product_id": 1,
							"quantity": 10,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidUserID,
		},
		"error_invalid_product_id": {
			given: givenData{
				reqBody: `{
					"user_id": 1,
					"items": [
						{
							"product_id": -1,
							"quantity": 10,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidProductID,
		},
		"error_invalid_quantity": {
			given: givenData{
				reqBody: `{
					"user_id": 1,
					"items": [
						{
							"product_id": 1,
							"quantity": 0,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidQuantity,
		},
		"error_invalid_discount": {
			given: givenData{
				reqBody: `{
					"user_id": 1,
					"items": [
						{
							"product_id": 1,
							"quantity": 10,
							"discount": 2,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidDiscount,
		},
		"error_user_is_not_exist": {
			given: givenData{
				reqBody: `{
					"note": "New order",
					"user_id": 3,
					"items": [
						{
							"product_id": 1,
							"quantity": 10,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
				mock: mockData{
					input: order.OrderInput{
						Note:   "New order",
						UserID: 3,
						Items: []order.OrderItemInput{
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
					err: order.ErrUserNotExist,
				},
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr:       ErrUserNotExist,
			isCallToServ: true,
		},
		"error_product_is_not_exist": {
			given: givenData{
				reqBody: `{
					"note": "New order",
					"user_id": 3,
					"items": [
						{
							"product_id": 9,
							"quantity": 10,
							"discount": 0,
							"note":"item 1"
						},
						{
							"product_id": 3,
							"quantity": 20,
							"discount": 0.5,
							"note": "item 2"
						}
					]
				}`,
				mock: mockData{
					input: order.OrderInput{
						Note:   "New order",
						UserID: 3,
						Items: []order.OrderItemInput{
							{
								ProductID: 9,
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
					err: order.ErrProductNotExist,
				},
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr:       ErrProductNotFound,
			isCallToServ: true,
		},
	}

	for decs, tc := range tcs {
		t.Run(decs, func(t *testing.T) {
			// Given
			r := httptest.NewRequest("POST", "/api/v1/orders", strings.NewReader(tc.given.reqBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			servMock := new(order.Mock)
			if tc.isCallToServ {
				servMock.On("CreateOrder", r.Context(), tc.given.mock.input).Return(tc.given.mock.err)
			}

			h := NewHandler(nil, nil, servMock)

			// When
			h.CreateOrder(w, r)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.Equal(t, tc.expResult.body, w.Body.String())
			}
			if tc.isCallToServ {
				servMock.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_GetOrders(t *testing.T) {
	type mockGetOrdersData struct {
		input            orderServ.OrdersInput
		resultErr        error
		result           []orderServ.Order
		resultTotalCount int64
	}
	type input struct {
		body     string
		mockData mockGetOrdersData
	}
	type output struct {
		result     OrdersResponse
		err        error
		statusCode int
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				body: `{
					"filter":{
						"id": 1,
						"order_number": "abcd",
						"status":"NEW",
						"user_id": 2
					},
					"sort_by":{
						"created_at":"asc",
						"order_date":"asc"
					},
					"pagination":{
						"page":10,
						"limit":2
					}
				}
				`,
				mockData: mockGetOrdersData{
					input: orderServ.OrdersInput{
						Filter: orderServ.OrderFilter{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							UserID:      2,
						},
						SortBy: orderServ.OrderSortBy{
							OrderDate: "asc",
							CreatedAt: "asc",
						},
						Pagination: orderServ.Pagination{
							Limit: 2,
							Page:  10,
						},
					},
					result: []order.Order{
						{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							Note:        "This is note",
							UserID:      2,
							OrderItems: []orderServ.OrderItem{
								{
									ID:           2,
									ProductID:    3,
									ProductPrice: 40000,
									ProductName:  "Product",
									Quantity:     3,
								},
							},
						},
					},
					resultTotalCount: 1,
				},
			},
			expOutput: output{
				result: OrdersResponse{
					Orders: []Order{
						{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							Note:        "This is note",
							UserID:      2,
							OrderItems: []OrderItem{
								{
									ID:           2,
									ProductID:    3,
									ProductPrice: 40000,
									ProductName:  "Product",
									Quantity:     3,
								},
							},
						},
					},
					Pagination: pagination{
						CurrentPage: 10,
						Limit:       2,
						TotalCount:  1,
					},
				},
				statusCode: http.StatusOK,
			},
		},
		"invalid_request_body": {
			input: input{
				body: `{{dklsflks`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidBodyRequest,
			},
		},
		"invalid_order_id": {
			input: input{
				body: `{
					"filter": {
						"id": -1
					}
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidOrderID,
			},
		},
		"invalid_user_id": {
			input: input{
				body: `{
					"filter": {
						"user_id": -1
					}
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidUserID,
			},
		},
		"invalid_order_status": {
			input: input{
				body: `{
					"filter": {
						"status":"Newssss"
					}
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidOrderStatus,
			},
		},
		"invalid_sort_type": {
			input: input{
				body: `{
					"sort_by": {
						"created_at":"asccc"
					}
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidOrderBy,
			},
		},
		"empty_body": {
			input: input{
				mockData: mockGetOrdersData{
					input: orderServ.OrdersInput{
						Pagination: orderServ.Pagination{
							Limit: 20,
							Page:  1,
						},
					},
					result: []order.Order{
						{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							Note:        "This is note",
							UserID:      2,
							OrderItems: []orderServ.OrderItem{
								{
									ID:           2,
									ProductID:    3,
									ProductPrice: 40000,
									ProductName:  "Product",
									Quantity:     3,
								},
							},
						},
					},
					resultTotalCount: 1,
				},
			},
			expOutput: output{
				result: OrdersResponse{
					Orders: []Order{
						{
							ID:          1,
							OrderNumber: "abcd",
							Status:      "NEW",
							Note:        "This is note",
							UserID:      2,
							OrderItems: []OrderItem{
								{
									ID:           2,
									ProductID:    3,
									ProductPrice: 40000,
									ProductName:  "Product",
									Quantity:     3,
								},
							},
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						Limit:       20,
						TotalCount:  1,
					},
				},
				statusCode: http.StatusOK,
			},
		},
	}
	for decs, tc := range tcs {
		t.Run(decs, func(t *testing.T) {
			// GIVEN
			var body io.Reader = nil
			if tc.input.body != "" {
				body = strings.NewReader(tc.input.body)
			}
			r := httptest.NewRequest(http.MethodGet, "/api/v1/orders", body)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			servMock := new(order.Mock)
			servMock.On("GetOrders", r.Context(), tc.input.mockData.input).Return(tc.input.mockData.result, tc.input.mockData.resultTotalCount, tc.input.mockData.resultErr)
			h := NewHandler(nil, nil, servMock)

			// WHEN
			h.GetOrders(w, r)

			// THEN
			if tc.expOutput.err != nil {
				require.Equal(t, tc.expOutput.statusCode, w.Code)
				require.EqualError(t, tc.expOutput.err, w.Body.String())
			} else {
				require.Equal(t, tc.expOutput.statusCode, w.Code)

				var result OrdersResponse
				err := json.Unmarshal(w.Body.Bytes(), &result)
				require.NoError(t, err)

				require.Equal(t, tc.expOutput.result, result)
			}

		})
	}
}
