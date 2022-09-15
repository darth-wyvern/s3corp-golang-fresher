package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/null/v8"

	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	userServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
)

func TestHandler_CreateUser(t *testing.T) {
	type MockData struct {
		userInput userServ.InputUser
		result    model.User
		err       error
	}

	type GivenData struct {
		reqBody      string
		mockData     MockData
		isCallToServ bool
	}

	type expectedData struct {
		statusCode int
		data       model.User
	}

	tcs := map[string]struct {
		given     GivenData
		expResult expectedData
		expErr    error
	}{
		"success": {
			given: GivenData{
				reqBody: `{
						"name": "guest",
						"email": "guest@example.com",
						"password": "123456",
						"phone": "123456",
						"role": "GUEST",
						"is_active": true
					}`,
				mockData: MockData{
					userInput: userServ.InputUser{
						Name:     "guest",
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
					result: model.User{
						ID:       1,
						Name:     "guest",
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: 201,
				data: model.User{
					ID:       1,
					Name:     "guest",
					Email:    "guest@example.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
		},
		"error-invalid-email": {
			given: GivenData{
				reqBody: `{
						"name": "guest",
						"email": "guestexamplecom",
						"password": "123456",
						"phone": "123456",
						"role": "GUEST",
						"is_active": true
					}`,
				mockData: MockData{
					userInput: userServ.InputUser{
						Name:     "guest",
						Email:    "guestexamplecom",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
					result: model.User{
						ID:       1,
						Name:     "guest",
						Email:    "guestexamplecom",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
				},
			},
			expResult: expectedData{
				statusCode: 400,
			},
			expErr: ErrInvalidEmail,
		},
		"error-invalid-role": {
			given: GivenData{
				reqBody: `{
						"name": "guest",
						"email": "guest@example.com",
						"password": "123456",
						"phone": "123456",
						"role": "invalid",
						"is_active": true
					}`,
				mockData: MockData{
					userInput: userServ.InputUser{
						Name:     "guest",
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "invalid",
						IsActive: true,
					},
					result: model.User{
						ID:       1,
						Name:     "guest",
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "invalid",
						IsActive: true,
					},
				},
			},
			expResult: expectedData{
				statusCode: 400,
			},
			expErr: ErrInvalidRole,
		},
		"error-missing-name-field": {
			given: GivenData{
				reqBody: `{
						"email": "guest@example.com",
						"password": "123456",
						"phone": "123456",
						"role": "GUEST",
						"is_active": true
					}`,
				mockData: MockData{
					userInput: userServ.InputUser{
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "ADMIN",
						IsActive: true,
					},
					result: model.User{
						ID:       1,
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "ADMIN",
						IsActive: true,
					},
				},
			},
			expResult: expectedData{
				statusCode: 400,
			},
			expErr: ErrNameCannotBeBlank,
		},
		"error-email-duplicate": {
			given: GivenData{
				reqBody: `{
						"name": "guest",
						"email": "guest@example.com",
						"password": "123456",
						"phone": "123456",
						"role": "GUEST",
						"is_active": true
					}`,
				mockData: MockData{
					userInput: userServ.InputUser{
						Name:     "guest",
						Email:    "guest@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
					result: model.User{},
					err:    userServ.ErrEmailExisted,
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: 400,
			},
			expErr: ErrEmailExisted,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			serviceMock := new(userServ.Mock)
			serviceMock.On("CreateUser", context.Background(), tc.given.mockData.userInput).Return(tc.given.mockData.result, tc.given.mockData.err)
			handler := NewHandler(serviceMock, nil, nil)

			r := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(tc.given.reqBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			// When
			handler.CreateUser(w, r)

			// Then
			if tc.expErr != nil {
				// Should be error
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				// Should be success
				require.Equal(t, tc.expResult.statusCode, w.Code)

				var actualResult model.User
				if err := json.Unmarshal(w.Body.Bytes(), &actualResult); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.expResult.data, actualResult, "Should be equal expected result")
			}

			if tc.given.isCallToServ {
				serviceMock.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_GetUsers(t *testing.T) {
	type mockData struct {
		input      userServ.InputGetUser
		totalCount int64
		output     []model.User
		err        error
	}

	type givenData struct {
		reqBody      string
		mock         mockData
		isCallToServ bool
	}

	type expectedData struct {
		statusCode int
		data       usersResponse
	}

	tcs := map[string]struct {
		given     givenData
		expResult expectedData
		expErr    error
	}{
		"success_one_field": {
			given: givenData{
				reqBody: `{
					"id": 1
				}`,
				mock: mockData{
					input: userServ.InputGetUser{
						ID:         1,
						Pagination: userServ.Pagination{Page: 1, Limit: 20},
					},
					totalCount: 1,
					output: []model.User{
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: http.StatusOK,
				data: usersResponse{
					Users: []model.User{
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						Limit:       20,
						TotalCount:  1,
					},
				},
			},
		},
		"success_multi_field": {
			given: givenData{
				reqBody: `{
					"is_active": true,
					"role": "ADMIN"
				}`,
				mock: mockData{
					input: userServ.InputGetUser{
						IsActive:   null.Bool{Valid: true, Bool: true},
						Role:       "ADMIN",
						Pagination: userServ.Pagination{Page: 1, Limit: 20},
					},
					output: []model.User{
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       2,
							Name:     "test2",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					totalCount: 2,
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: http.StatusOK,
				data: usersResponse{
					Users: []model.User{
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       2,
							Name:     "test2",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						Limit:       20,
						TotalCount:  2,
					},
				},
			},
		},
		"success_multi_field_with_sort": {
			given: givenData{
				reqBody: `{
					"is_active": true,
					"role": "ADMIN",
					"sort": {
						"name": "desc",
						"created_at": "asc"
					}
				}`,
				mock: mockData{
					input: userServ.InputGetUser{
						IsActive: null.Bool{Bool: true, Valid: true},
						Role:     "ADMIN",
						Sort: userServ.SortArgs{
							Name:      OrderTypeDESC,
							CreatedAt: OrderTypeASC,
						},
						Pagination: userServ.Pagination{Page: 1, Limit: 20},
					},
					output: []model.User{
						{
							ID:       2,
							Name:     "test2",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					totalCount: 2,
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: http.StatusOK,
				data: usersResponse{
					Users: []model.User{
						{
							ID:       2,
							Name:     "test2",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       1,
							Name:     "test",
							Email:    "test@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					Pagination: pagination{
						CurrentPage: 1,
						TotalCount:  2,
						Limit:       20,
					},
				},
			},
		},
		"error_invalid_id": {
			given: givenData{
				reqBody: `{
					"id": -2
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidID,
		},
		"error_invalid_email": {
			given: givenData{
				reqBody: `{
					"email": "exam.com"
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidEmail,
		},
		"error_invalid_role": {
			given: givenData{
				reqBody: `{
					"role": "SUPERMAN"
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidRole,
		},
		"error_invalid_sort_order_type": {
			given: givenData{
				reqBody: `{
					"sort": {
						"name": "a -> z"
					}
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidSortType,
		},
		"error_invalid_multi_sort_order_type": {
			given: givenData{
				reqBody: `{
					"sort": {
						"name": "asc",
						"created_at": "z -> a"
					}
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidSortType,
		},
		"error_invalid_pagination_page": {
			given: givenData{
				reqBody: `{
					"pagination": {
						"page": -1,
						"limit": 20
					}
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidPaginationPage,
		},
		"error_invalid_pagination_limit": {
			given: givenData{
				reqBody: `{
					"pagination": {
						"page": 1,
						"limit": -1
					}
				}`,
				isCallToServ: false,
			},
			expResult: expectedData{
				statusCode: http.StatusBadRequest,
			},
			expErr: ErrInvalidPaginationLimit,
		},
		"error_not_found": {
			given: givenData{
				reqBody: `{
					"pagination": {
						"page": 3,
						"limit": 2
					}
				}`,
				mock: mockData{
					input: userServ.InputGetUser{
						Pagination: userServ.Pagination{
							Page:  3,
							Limit: 2,
						},
					},
					output:     []model.User{},
					totalCount: 4,
					err:        userServ.ErrUserNotFound,
				},
				isCallToServ: true,
			},
			expResult: expectedData{
				statusCode: http.StatusNotFound,
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			serviceMock := new(userServ.Mock)
			if tc.given.isCallToServ {
				serviceMock.On("GetUsers", context.Background(), tc.given.mock.input).Return(tc.given.mock.output, tc.given.mock.totalCount, tc.given.mock.err)
			}
			handler := NewHandler(serviceMock, nil, nil)

			r := httptest.NewRequest("GET", "/api/v1/users", strings.NewReader(tc.given.reqBody))

			w := httptest.NewRecorder()

			// When
			handler.GetUsers(w, r)

			// Then
			if tc.expErr != nil {
				// Should be error
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				// Should be success
				require.Equal(t, tc.expResult.statusCode, w.Code)

				var actualResult usersResponse
				if err := json.Unmarshal(w.Body.Bytes(), &actualResult); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.expResult.data, actualResult, "Should be equal expected result")
			}

			if tc.given.isCallToServ {
				serviceMock.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	type input struct {
		userID     string
		reqBody    string // json
		mockInput  userServ.InputUser
		mockResult error
	}
	type output struct {
		expResult     string // json
		expErr        error
		expStatusCode int
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				userID: "1",
				reqBody: `{
					"name": "guest",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
				mockInput: userServ.InputUser{
					ID:       1,
					Name:     "guest",
					Email:    "guest@example.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				expStatusCode: http.StatusOK,
				expResult:     `{"success":true,"msg":"Update user successfully"}`,
			},
		},
		"not_found": {
			input: input{
				userID: "2",
				reqBody: `{
					"name": "guest",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
				mockInput: userServ.InputUser{
					ID:       2,
					Name:     "guest",
					Email:    "guest@example.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
				mockResult: userServ.ErrUserNotFound,
			},
			expOutput: output{
				expStatusCode: http.StatusNotFound,
				expErr:        ErrUserNotFound,
			},
		},
		"invalid_id": {
			input: input{
				userID: "d-2",
				reqBody: `{
					"name": "guest",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrInvalidUserID,
			},
		},
		"invalid_email": {
			input: input{
				userID: "2",
				reqBody: `{
					"name": "guest",
					"email": "guestexample.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrInvalidEmail,
			},
		},
		"invalid_role": {
			input: input{
				userID: "2",
				reqBody: `{
					"name": "guest",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "Nope",
					"is_active": true
				}`,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrInvalidRole,
			},
		},
		"missing_name_field": {
			input: input{
				userID: "2",
				reqBody: `{
					"name": "",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrNameCannotBeBlank,
			},
		},
		"missing_password_field": {
			input: input{
				userID: "1",
				reqBody: `{
					"name": "sdfaf",
					"email": "guest@example.com",
					"password": "",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrPasswordCannotBeBlank,
			},
		},
		"email_duplicated": {
			input: input{
				userID: "2",
				reqBody: `{
					"name": "sdfaf",
					"email": "guest@example.com",
					"password": "123456",
					"phone": "123456",
					"role": "GUEST",
					"is_active": true
				}`,
				mockInput: userServ.InputUser{
					ID:       2,
					Name:     "sdfaf",
					Email:    "guest@example.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
				mockResult: userServ.ErrEmailExisted,
			},
			expOutput: output{
				expStatusCode: http.StatusBadRequest,
				expErr:        ErrEmailExisted,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			r := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tc.input.userID, strings.NewReader(tc.input.reqBody))
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.input.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			serviceMock := new(userServ.Mock)
			serviceMock.On("UpdateUser", r.Context(), tc.input.mockInput).Return(tc.input.mockResult)
			handler := NewHandler(serviceMock, nil, nil)
			// When
			handler.UpdateUser(w, r)
			//THEN
			if tc.expOutput.expErr != nil {
				//must be error
				require.Equal(t, tc.expOutput.expStatusCode, w.Code)
				require.EqualError(t, tc.expOutput.expErr, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expOutput.expStatusCode, w.Code)
				require.Equal(t, tc.expOutput.expResult, w.Body.String())
			}
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	type mockData struct {
		userID int
		result model.User
		err    error
	}

	type givenData struct {
		userID string
		mock   mockData
	}

	type expectedData struct {
		statusCode int
		result     model.User
	}

	tcs := map[string]struct {
		given     givenData
		expResult expectedData
		expErr    error
	}{
		"success": {
			given: givenData{
				userID: "1",
				mock: mockData{
					userID: 1,
					result: model.User{
						ID:       1,
						Name:     "test",
						Email:    "test@example.com",
						Password: "123456",
						Phone:    "123456",
						Role:     "GUEST",
						IsActive: true,
					},
				},
			},
			expResult: expectedData{
				statusCode: http.StatusOK,
				result: model.User{
					ID:       1,
					Name:     "test",
					Email:    "test@example.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
		},
		"error": {
			given: givenData{
				userID: "1",
				mock: mockData{
					userID: 1,
					result: model.User{},
					err:    userServ.ErrUserNotFound,
				},
			},
			expResult: expectedData{
				statusCode: http.StatusNotFound,
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			r := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tc.given.userID, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.given.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			serviceMock := new(userServ.Mock)
			serviceMock.On("GetUser", r.Context(), tc.given.mock.userID).Return(tc.given.mock.result, tc.given.mock.err)

			handler := NewHandler(serviceMock, nil, nil)

			// When
			handler.GetUser(w, r)

			// Then
			if tc.expErr != nil {
				//must be error
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expResult.statusCode, w.Code)

				var actualResult model.User
				if err := json.Unmarshal(w.Body.Bytes(), &actualResult); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.expResult.result, actualResult)
			}
			serviceMock.AssertExpectations(t)

		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	type mockData struct {
		userID int
		err    error
	}

	type givenData struct {
		userID string
		mock   mockData
	}

	type expectedData struct {
		statusCode int
		result     string
	}

	tcs := map[string]struct {
		given     givenData
		expResult expectedData
		expErr    error
	}{
		"success": {
			given: givenData{
				userID: "1",
				mock: mockData{
					userID: 1,
					err:    nil,
				},
			},
			expResult: expectedData{
				statusCode: http.StatusOK,
				result:     "{\"success\":true,\"msg\":\"Delete user successfully\"}",
			},
		},
		"error": {
			given: givenData{
				userID: "1",
				mock: mockData{
					userID: 1,
					err:    userServ.ErrUserNotFound,
				},
			},
			expResult: expectedData{
				statusCode: http.StatusNotFound,
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tc.given.userID, nil)
			w := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.given.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			serviceMock := new(userServ.Mock)
			serviceMock.On("DeleteUser", r.Context(), tc.given.mock.userID).Return(tc.given.mock.err)

			handler := NewHandler(serviceMock, nil, nil)

			// When
			handler.DeleteUser(w, r)

			// Then
			if tc.expErr != nil {
				//must be error
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.EqualError(t, tc.expErr, w.Body.String())
			} else {
				//must be success
				require.Equal(t, tc.expResult.statusCode, w.Code)
				require.Equal(t, tc.expResult.result, w.Body.String())
			}
			serviceMock.AssertExpectations(t)

		})
	}
}

func TestHandler_Login(t *testing.T) {
	type input struct {
		reqBody       string
		mockInput     userServ.LoginInput
		mockResult    userServ.LoginResponse
		mockResultErr error
	}
	type output struct {
		body       userServ.LoginResponse
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
				   "email":"example@example.com",
				   "password":"123456789"
				}`,
				mockInput: userServ.LoginInput{
					Email:    "example@example.com",
					Password: "123456789",
				},
				mockResult: userServ.LoginResponse{
					AccessToken: "fbhaffiuerfweifewiuffgefhjgfiuyfr",
					Scope:       "GUEST",
					ExpiresIn:   1800,
					TokenType:   "Bearer",
				},
			},
			expOutput: output{
				statusCode: http.StatusOK,
				body: userServ.LoginResponse{
					AccessToken: "fbhaffiuerfweifewiuffgefhjgfiuyfr",
					Scope:       "GUEST",
					ExpiresIn:   1800,
					TokenType:   "Bearer",
				},
			},
		},
		"invalid_email": {
			input: input{
				reqBody: `{
				   "email":"exampleexample.com",
				   "password":"123456789"
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidEmail,
			},
		},
		"password_can_not_be_blank": {
			input: input{
				reqBody: `{
				   "email":"example@example.com",
				   "password":""
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrPasswordCannotBeBlank,
			},
		},
		"invalid_request_body": {
			input: input{
				reqBody: `{
				   "email":"example@example.com",
				   "password":"123456",
				}`,
			},
			expOutput: output{
				statusCode: http.StatusBadRequest,
				err:        ErrInvalidBodyRequest,
			},
		},
		"password_is_incorrect": {
			input: input{
				reqBody: `{
				   "email":"example@example.com",
				   "password":"123456789"
				}`,
				mockInput: userServ.LoginInput{
					Email:    "example@example.com",
					Password: "123456789",
				},
				mockResultErr: userServ.ErrPasswordIncorrect,
			},
			expOutput: output{
				statusCode: http.StatusUnauthorized,
				err:        ErrPasswordIncorrect,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			r := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", strings.NewReader(tc.input.reqBody))
			w := httptest.NewRecorder()

			serviceMock := new(userServ.Mock)
			serviceMock.On("Login", r.Context(), tc.input.mockInput).Return(tc.input.mockResult, tc.input.mockResultErr)

			handler := NewHandler(serviceMock, nil, nil)

			// WHEN
			handler.Login(w, r)

			//THEN
			if tc.expOutput.err != nil {
				// must be error
				require.Equal(t, tc.expOutput.statusCode, w.Code)
				require.EqualError(t, tc.expOutput.err, w.Body.String())
			} else {
				require.Equal(t, tc.expOutput.statusCode, w.Code)

				var result userServ.LoginResponse
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.expOutput.body, result)
			}
		})
	}
}
