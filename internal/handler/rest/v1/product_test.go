package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	productService "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
)

func TestProductHandler_GetProduct(t *testing.T) {
	type input struct {
		productID         int
		mockInput         int
		mockResultProduct model.Product
		mockResultError   error
	}
	type output struct {
		expBody   model.Product
		expStatus int
		expErr    error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				productID: 1,
				mockInput: 1,
				mockResultProduct: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				mockResultError: nil,
			},
			expOutput: output{
				expBody: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				expStatus: http.StatusOK,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			// Define http test request
			// Define http test response
			r := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+strconv.Itoa(tc.input.productID), nil)
			w := httptest.NewRecorder()
			// Init chi route context
			// Set id to chi route context
			// Add chi route context to request
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.input.productID))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			productServiceMock := new(productService.Mock)
			productServiceMock.On("GetProduct", r.Context(), tc.input.productID).Return(tc.input.mockResultProduct, tc.input.mockResultError)
			handler := NewHandler(nil, productServiceMock, nil)

			// WHEN
			handler.GetProduct(w, r)

			//THEN
			if tc.expOutput.expErr != nil {
				//must be error
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				require.EqualError(t, tc.expOutput.expErr, w.Body.String()) // Equal body if error case
			} else {
				// must be success
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				// Read file
				var result model.Product
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Error("Error on creating result")
				}
				require.Equal(t, tc.expOutput.expBody, result) // Equal body if success case
			}
			productServiceMock.AssertExpectations(t)
		})
	}

}

func TestProductHandler_CreateProduct(t *testing.T) {
	type input struct {
		reqBody           string //json
		mockInput         productService.ProductInput
		mockResultProduct model.Product
		mockResultError   error
	}
	type output struct {
		expBody   model.Product
		expStatus int
		expErr    error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				reqBody: `{
					"title": "test",
					"description": "",
					"price":20000,
					"quantity": 10,
					"is_active": true,
					"user_id": 1
					}`,

				mockInput: productService.ProductInput{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				mockResultProduct: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
			},
			expOutput: output{
				expBody: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				expStatus: http.StatusCreated,
			},
		},
		"title_is_required": {
			input: input{
				reqBody: `{
					"description": "",
					"price":20000,
					"quantity": 10,
					"is_active": true,
					"user_id": 1
					}`,
			},
			expOutput: output{
				expErr:    ErrTitleCannotBeBlank,
				expStatus: http.StatusBadRequest,
			},
		},
		"invalid_price": {
			input: input{
				reqBody: `{
					"title": "test",
					"description": "",
					"price":-1,
					"quantity": 10,
					"is_active": true,
					"user_id": 1
					}`,
			},
			expOutput: output{
				expErr:    ErrInvalidPrice,
				expStatus: http.StatusBadRequest,
			},
		},
		"invalid_user_id": {
			input: input{
				reqBody: `{
					"title": "test",
					"description": "",
					"price":100000,
					"quantity": 10,
					"is_active": true,
					"user_id": -11
					}`,
			},
			expOutput: output{
				expErr:    ErrInvalidUserID,
				expStatus: http.StatusBadRequest,
			},
		},
		"invalid_quantity": {
			input: input{
				reqBody: `{
					"title": "test",
					"description": "",
					"price":100000,
					"quantity": -10,
					"is_active": true,
					"user_id": 1
					}`,
			},
			expOutput: output{
				expErr:    ErrInvalidQuantity,
				expStatus: http.StatusBadRequest,
			},
		},
		"user_is_not_exist": {
			input: input{
				reqBody: `{
					"title": "test",
					"description": "",
					"price":100000,
					"quantity": 10,
					"is_active": true,
					"user_id": 100
					}`,
				mockInput: productService.ProductInput{
					Title:       "test",
					Description: "",
					Price:       100000,
					Quantity:    10,
					IsActive:    true,
					UserID:      100,
				},
				mockResultError: productService.ErrUserNotExist,
			},
			expOutput: output{
				expErr:    ErrUserNotExist,
				expStatus: http.StatusNotFound,
			},
		},
		"invalid_request_body": {
			input: input{
				reqBody: `this is a text`,
			},
			expOutput: output{
				expErr:    ErrInvalidBodyRequest,
				expStatus: http.StatusBadRequest,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			serviceMock := new(productService.Mock)
			serviceMock.On("CreateProduct", context.Background(), tc.input.mockInput).Return(tc.input.mockResultProduct, tc.input.mockResultError)
			handler := NewHandler(nil, serviceMock, nil)

			fixJson := strings.NewReplacer(
				"\n", "",
				"\t", "")
			fixedBody := fixJson.Replace(tc.input.reqBody)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(fixedBody))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			//WHEN
			handler.CreateProduct(w, r)

			//THEN
			if tc.expOutput.expErr != nil {
				// must be error
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				require.EqualError(t, tc.expOutput.expErr, w.Body.String())
			} else {
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				var result model.Product
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				require.Equal(t, tc.expOutput.expBody, result)
			}

		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	type mockData struct {
		id    int
		input productService.ProductInput
		err   error
	}

	type givenData struct {
		id           int
		reqBody      string
		mock         mockData
		isCallToServ bool
	}

	type expectedData struct {
		expResult  string
		statusCode int
	}

	tcs := map[string]struct {
		given    givenData
		expected expectedData
		expErr   error
	}{
		"success": {
			given: givenData{
				id: 1,
				reqBody: `{
					"title": "LED",
					"price": 1000,
					"quantity": 10,
					"is_active": true,
					"user_id": 3
				}`,
				mock: mockData{
					id: 1,
					input: productService.ProductInput{
						Title:       "LED",
						Description: "",
						Price:       1000,
						Quantity:    10,
						IsActive:    true,
						UserID:      3,
					},
					err: nil,
				},
				isCallToServ: true,
			},
			expected: expectedData{
				statusCode: http.StatusOK,
				expResult:  `{"success":true,"msg":"Product updated successfully"}`,
			},
		},
		"error_missing_title": {
			given: givenData{
				id: 1,
				reqBody: `{
					"price": 1000,
					"quantity": 10,
					"is_active": true,
					"user_id": 3
				}`,
			},
			expected: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrTitleCannotBeBlank,
		},
		"error_invalid_price": {
			given: givenData{
				id: 1,
				reqBody: `{
					"title": "LED",
					"price": -9,
					"quantity": 10,
					"is_active": true,
					"user_id": 3
				}`,
			},
			expected: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidPrice,
		},
		"error_invalid_quantity": {
			given: givenData{
				id: 1,
				reqBody: `{
					"title": "LED",
					"price": 1000,
					"quantity": -5,
					"is_active": true,
					"user_id": 3
				}`,
			},
			expected: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidQuantity,
		},
		"error_product_not_found": {
			given: givenData{
				id: 1,
				reqBody: `{
					"title": "LED",
					"price": 1000,
					"quantity": 10,
					"is_active": true,
					"user_id": 3
				}`,
				mock: mockData{
					id: 1,
					input: productService.ProductInput{
						Title:       "LED",
						Description: "",
						Price:       1000,
						Quantity:    10,
						IsActive:    true,
						UserID:      3,
					},
					err: productService.ErrProductNotFound,
				},
				isCallToServ: true,
			},
			expected: expectedData{
				statusCode: http.StatusNotFound,
			},
			expErr: ErrProductNotFound,
		},
		"error_user_not_found": {
			given: givenData{
				id: 1,
				reqBody: `{
					"title": "LED",
					"price": 1000,
					"quantity": 10,
					"is_active": true,
					"user_id": 3
				}`,
				mock: mockData{
					id: 1,
					input: productService.ProductInput{
						Title:       "LED",
						Description: "",
						Price:       1000,
						Quantity:    10,
						IsActive:    true,
						UserID:      3,
					},
					err: productService.ErrUserNotExist,
				},
				isCallToServ: true,
			},
			expected: expectedData{
				statusCode: http.StatusNotFound,
			},
			expErr: ErrUserNotExist,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			r := httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(tc.given.id), strings.NewReader(tc.given.reqBody))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tc.given.id))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			serviceMock := new(productService.Mock)
			if tc.given.isCallToServ {
				serviceMock.On("UpdateProduct", r.Context(), tc.given.mock.id, tc.given.mock.input).Return(tc.given.mock.err)
			}

			handler := NewHandler(nil, serviceMock, nil)

			// When
			handler.UpdateProduct(w, r)

			// Then
			if tc.expErr != nil {
				// must be error
				require.Equal(t, tc.expected.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				require.Equal(t, tc.expected.statusCode, w.Code)
				require.Equal(t, tc.expected.expResult, w.Body.String())
			}

			if tc.given.isCallToServ {
				serviceMock.AssertExpectations(t)
			}
		})
	}

}

func TestProductHandler_DeleteProduct(t *testing.T) {
	type input struct {
		productID       string // url param
		mockInputID     int
		mockResultError error
	}
	type output struct {
		expResult string
		expStatus int
		expErr    error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				productID:   "1",
				mockInputID: 1,
			},
			expOutput: output{
				expResult: "Delete product successfully",
				expStatus: http.StatusOK,
			},
		},
		"invalid_id": {
			input: input{
				productID: "-2dsfasf",
			},
			expOutput: output{
				expErr:    ErrInvalidID,
				expStatus: http.StatusBadRequest,
			},
		},
		"not_found": {
			input: input{
				productID: "2",

				mockInputID:     2,
				mockResultError: productService.ErrProductNotFound,
			},
			expOutput: output{
				expErr:    ErrProductNotFound,
				expStatus: http.StatusNotFound,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {

			//GIVEN
			// 1. Define new request and response
			// Add "id" into request context
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+tc.input.productID, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.input.productID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			// 2. Define mock and handler
			serviceMock := new(productService.Mock)
			serviceMock.On("DeleteProduct", r.Context(), tc.input.mockInputID).Return(tc.input.mockResultError)
			handler := NewHandler(nil, serviceMock, nil)

			//WHEN
			handler.DeleteProduct(w, r)

			//THEN
			if tc.expOutput.expErr != nil {
				// must be error
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				require.EqualError(t, tc.expOutput.expErr, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expOutput.expStatus, w.Code)
				require.Equal(t, tc.expOutput.expResult, w.Body.String())
			}
		})
	}
}

func TestProductHandler_GetProducts(t *testing.T) {
	type input struct {
		reqBody              string // json
		mockInput            productService.GetProductsInput
		mockResultProducts   []productService.ProductItem
		mockResultTotalCount int64
		mockResultError      error
	}
	type output struct {
		result     getProductsResponse
		statusCode int
		err        error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				reqBody: `{
					"title":"test",
					"price_range":{"min_price":100, "max_price":3000},
					"is_active":true,
					"user_id":1,
					"order_by":{
						"title":"desc",
						"created_at":"desc"
					}
				}`,
				mockInput: productService.GetProductsInput{
					Title:      "test",
					PriceRange: productService.PriceRange{MinPrice: 100, MaxPrice: 3000},
					IsActive:   null.NewBool(true, true),
					UserID:     1,
					OrderBy: productService.OrderInput{
						Title:     "desc",
						CreatedAt: "desc",
					},
				},
				mockResultProducts: []productService.ProductItem{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						IsActive:    true,
						User: productService.CreatedBy{
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
						User: productService.CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
				},
				mockResultTotalCount: 2,
			},
			expOutput: output{
				result: getProductsResponse{
					Products: []productItemResponse{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: createdByResponse{
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
							User: createdByResponse{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						Limit:       20,
						TotalCount:  2,
					},
				},
				statusCode: http.StatusOK,
			},
		},
		"invalid_id": {
			input: input{
				reqBody: `{"id": -1}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidID,
			},
		},
		"invalid_price": {
			input: input{
				reqBody: `{
					"price_range":{
						"min_price":100,
						"max_price":-3000}
					}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidPriceRange,
			},
		},
		"invalid_user_id": {
			input: input{
				reqBody: `{
					"user_id": -1
					}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidUserID,
			},
		},
		"invalid_order_by": {
			input: input{
				reqBody: `{
					"order_by":{
							"title":"desc",
							"created_at":"dkjsfkheaf"
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
				mockResultProducts: []productService.ProductItem{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						IsActive:    true,
						User: productService.CreatedBy{
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
						User: productService.CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
				},
				mockResultTotalCount: 2,
			},
			expOutput: output{
				result: getProductsResponse{
					Products: []productItemResponse{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: createdByResponse{
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
							User: createdByResponse{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						Limit:       20,
						TotalCount:  2,
					},
				},
				statusCode: http.StatusOK,
			},
		},
		"paging_result": {
			input: input{
				reqBody: `{"pagination":{"page": 2,"limit":2}}`, //
				mockResultProducts: []productService.ProductItem{
					{
						ID:          1,
						Title:       "test",
						Description: "test",
						Price:       10000,
						Quantity:    10,
						IsActive:    true,
						User: productService.CreatedBy{
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
						User: productService.CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
				},
				mockInput: productService.GetProductsInput{
					Pagination: productService.Pagination{
						Page:  2,
						Limit: 2,
					},
				},
				mockResultTotalCount: 4,
			},
			expOutput: output{
				result: getProductsResponse{
					Products: []productItemResponse{
						{
							ID:          1,
							Title:       "test",
							Description: "test",
							Price:       10000,
							Quantity:    10,
							IsActive:    true,
							User: createdByResponse{
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
							User: createdByResponse{
								ID:    1,
								Name:  "admin",
								Email: "admin@example.com",
								Phone: "0987654321",
							},
						},
					},
					Pagination: pagination{
						CurrentPage: 2,
						Limit:       2,
						TotalCount:  4,
					},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {

			//GIVEN
			var reqBody io.Reader = nil
			if tc.input.reqBody != `` {
				reqBody = strings.NewReader(tc.input.reqBody)
			}

			r := httptest.NewRequest(http.MethodGet, "/api/v1/products/", reqBody)
			w := httptest.NewRecorder()

			// 2. Define mock and handler
			serviceMock := new(productService.Mock)
			serviceMock.On("GetProducts", r.Context(), tc.input.mockInput).Return(tc.input.mockResultProducts, tc.input.mockResultTotalCount, tc.input.mockResultError)
			handler := NewHandler(nil, serviceMock, nil)

			//WHEN
			handler.GetProducts(w, r)

			//THEN
			if tc.expOutput.err != nil {
				// must be error
				require.Equal(t, tc.expOutput.statusCode, w.Code)
				require.EqualError(t, tc.expOutput.err, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expOutput.statusCode, w.Code)

				var result getProductsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Error(err)
				}

				require.Equal(t, tc.expOutput.result, result)
			}
		})
	}
}

func TestProductHandler_ImportProductCSV(t *testing.T) {
	type givenData struct {
		fileName string
		fileData string
		mockErr  error
	}
	type expectedData struct {
		statusCode int
		msg        string
	}

	tcs := map[string]struct {
		given  givenData
		exp    expectedData
		expErr error
	}{
		"success": {
			given: givenData{
				fileName: "products.csv",
				fileData: "Title,Description,Price,Quantity,Is_Active,User_ID\nMahalia,rBIsqaccTvxl,462.25,537,1,1\nBerget,jVLHPRgDbMHNeds,910.59,913,1,1\nDede, vNlZDZd YtaB,535.69,356,0,1",
				mockErr:  nil,
			},
			exp: expectedData{
				statusCode: http.StatusOK,
				msg:        "Starting import product from csv: products.csv",
			},
		},
		"error": {
			given: givenData{
				fileName: "products",
				fileData: "",
			},
			exp: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidFileType,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			fileReader := strings.NewReader(tc.given.fileData)

			reqBody := new(bytes.Buffer)
			mw := multipart.NewWriter(reqBody)
			formWriter, err := mw.CreateFormFile("file", tc.given.fileName)
			require.NoError(t, err)

			if _, err := io.Copy(formWriter, fileReader); err != nil {
				t.Fatal(err)
			}

			mw.Close()

			r := httptest.NewRequest(http.MethodPost, "/api/v1/products/import-csv", reqBody)
			r.Header.Set("Content-Type", mw.FormDataContentType())
			w := httptest.NewRecorder()

			file, _, err := r.FormFile("file")
			require.NoError(t, err)
			defer file.Close()

			serviceMock := new(productService.Mock)
			serviceMock.On("ImportProductCSV", r.Context(), tc.given.fileName, file).Return(tc.given.mockErr)
			handler := NewHandler(nil, serviceMock, nil)

			// When
			handler.ImportProductCSV(w, r)

			// Then
			if tc.expErr != nil {
				require.Equal(t, tc.exp.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				require.Equal(t, tc.exp.statusCode, w.Code)
				require.Equal(t, tc.exp.msg, w.Body.String())
			}
		})
	}
}

func TestProductHandler_ExportProductsCSV(t *testing.T) {
	type input struct {
		reqBody         string // json
		mockInput       productService.GetProductsInput
		mockResultURL   string
		mockResultError error
	}
	type output struct {
		body       productCSVResponse
		statusCode int
		err        error
	}

	appURL := os.Getenv("APP_URL")

	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				reqBody: `{
					"title":"test",
					"price_range":{"min_price":100, "max_price":3000},
					"is_active":true,
					"user_id":1,
					"order_by":{
						"title":"desc",
						"created_at":"desc"
					}
				}`,
				mockInput: productService.GetProductsInput{
					Title:      "test",
					PriceRange: productService.PriceRange{MinPrice: 100, MaxPrice: 3000},
					IsActive:   null.NewBool(true, true),
					UserID:     1,
					OrderBy: productService.OrderInput{
						Title:     "desc",
						CreatedAt: "desc",
					},
				},
				mockResultURL: "file_20222807.csv",
			},
			expOutput: output{
				body:       productCSVResponse{ProductCSVURL: appURL + "/api/v1/files/file_20222807.csv"},
				statusCode: http.StatusOK,
			},
		},
		"invalid_id": {
			input: input{
				reqBody: `{"id": -1}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidID,
			},
		},
		"invalid_price": {
			input: input{
				reqBody: `{
					"price_range":{
						"min_price":100,
						"max_price":-3000}
					}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidPriceRange,
			},
		},
		"invalid_user_id": {
			input: input{
				reqBody: `{
					"user_id": -1
					}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidUserID,
			},
		},
		"invalid_order_by": {
			input: input{
				reqBody: `{
					"order_by":{
							"title":"desc",
							"created_at":"dkjsfkheaf"
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
				mockInput:     productService.GetProductsInput{},
				mockResultURL: "file_20222807.csv",
			},
			expOutput: output{
				body:       productCSVResponse{ProductCSVURL: appURL + "/api/v1/files/file_20222807.csv"},
				statusCode: http.StatusOK,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			var reqBody io.Reader = nil
			if tc.input.reqBody != `` {
				reqBody = strings.NewReader(tc.input.reqBody)
			}

			r := httptest.NewRequest(http.MethodGet, "/api/v1/products/export/csv", reqBody)
			w := httptest.NewRecorder()

			// 2. Define mock and handler
			serviceMock := new(productService.Mock)
			serviceMock.On("ExportProductsCSV", r.Context(), tc.input.mockInput, "docs/csvfiles").Return(tc.input.mockResultURL, tc.input.mockResultError)
			handler := NewHandler(nil, serviceMock, nil)

			//WHEN
			handler.ExportProductsCSV(w, r)

			//THEN
			if tc.expOutput.err != nil {
				// must be error
				require.Equal(t, tc.expOutput.statusCode, w.Code)
				require.EqualError(t, tc.expOutput.err, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expOutput.statusCode, w.Code)

				var result productCSVResponse
				err := json.Unmarshal(w.Body.Bytes(), &result)
				require.NoError(t, err)

				require.Equal(t, tc.expOutput.body, result)
			}
		})
	}
}

func TestProductHandler_DownloadCSVFile(t *testing.T) {
	type input struct {
		fileName          string
		mockInputFilePath string
		mockResultFile    []byte
		mockResultErr     error
	}
	type output struct {
		statusCode int
		fileStr    string
		err        error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				fileName:          "products_20220729.csv",
				mockInputFilePath: "products_20220729.csv",
				mockResultFile:    []byte("1"),
			},
			expOutput: output{
				statusCode: http.StatusOK,
				fileStr:    "1",
			},
		},
		"invalid_file_name": {
			input: input{
				fileName: "products_20220729.exe",
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidFileName,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			r := httptest.NewRequest(http.MethodGet, "/api/v1/files/"+tc.input.fileName, nil)
			w := httptest.NewRecorder()
			// Set fileName to chi route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("filename", tc.input.fileName)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// 2. Define mock and handler
			serviceMock := new(productService.Mock)
			serviceMock.On("DownloadCSV", r.Context(), "docs/csvfiles/"+tc.input.mockInputFilePath).Return(tc.input.mockResultFile, tc.input.mockResultErr)
			handler := NewHandler(nil, serviceMock, nil)

			//WHEN
			handler.DownloadCSVFile(w, r)

			//THEN
			if tc.expOutput.err != nil {
				// must be error
				require.Equal(t, tc.expOutput.statusCode, w.Code)
				require.EqualError(t, tc.expOutput.err, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expOutput.statusCode, w.Code)

				require.Equal(t, tc.expOutput.fileStr, w.Body.String())
			}
		})
	}
}
