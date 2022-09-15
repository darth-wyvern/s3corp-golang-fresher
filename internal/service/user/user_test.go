package user

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
)

func TestUserService_CreateUser(t *testing.T) {
	type createUserData struct {
		input  mock.AnythingOfTypeArgument
		result model.User
		err    error
	}

	type existUserData struct {
		input  string
		result bool
		err    error
	}

	type givenData struct {
		input      InputUser
		createUser createUserData
		existUser  existUserData
	}

	type expectedData struct {
		data model.User
	}

	tcs := map[string]struct {
		given     givenData
		expResult expectedData
		expErr    error
	}{
		"success": {
			given: givenData{
				input: InputUser{
					Name:     "guest",
					Email:    "guest@example.com",
					Password: "abcd",
					Phone:    "0987654321",
					Role:     "GUEST",
					IsActive: true,
				},
				createUser: createUserData{
					input: mock.AnythingOfType("User"),
					result: model.User{
						ID:       1,
						Name:     "guest",
						Email:    "guest@example.com",
						Password: mock.Anything,
						Phone:    "0987654321",
						Role:     "GUEST",
						IsActive: true,
					},
				},
				existUser: existUserData{
					input:  "guest@example.com",
					result: false,
				},
			},
			expResult: expectedData{
				data: model.User{
					ID:       1,
					Name:     "guest",
					Email:    "guest@example.com",
					Password: mock.Anything,
					Phone:    "0987654321",
					Role:     "GUEST",
					IsActive: true,
				},
			},
		},
		"error_email_duplicate": {
			given: givenData{
				input: InputUser{
					Name:     "guest",
					Email:    "guest@example.com",
					Password: "abcd",
					Phone:    "0987654321",
					Role:     "GUEST",
					IsActive: true,
				},
				createUser: createUserData{
					input: mock.AnythingOfType("User"),
					result: model.User{
						ID:       1,
						Name:     "guest",
						Email:    "guest@example.com",
						Password: mock.Anything,
						Phone:    "0987654321",
						Role:     "GUEST",
						IsActive: true,
					},
				},
				existUser: existUserData{
					input:  "guest@example.com",
					result: true,
				},
			},
			expErr: ErrEmailExisted,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			userRepoMock := new(user.Mock)

			userRepoMock.On("CreateUser", context.Background(), tc.given.createUser.input).Return(tc.given.createUser.result, tc.given.createUser.err)
			userRepoMock.On("ExistsUserByEmail", context.Background(), tc.given.existUser.input).Return(tc.given.existUser.result, tc.given.existUser.err)

			repoMock := new(repository.Mock)
			repoMock.On("User").Return(userRepoMock)

			service := New(repoMock)

			// When
			result, err := service.CreateUser(context.Background(), tc.given.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, tc.expErr, "Should not be error")
				require.Equal(t, tc.expResult.data, result)
			}
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	type mockData struct {
		input      user.Filter
		result     model.UserSlice
		totalCount int64
		err        error
	}

	type givenData struct {
		input        InputGetUser
		mock         mockData
		isCallToRepo bool
	}

	tcs := map[string]struct {
		given     givenData
		expResult []model.User
		expErr    error
	}{
		"success_field_with_sort": {
			given: givenData{
				input: InputGetUser{
					Name: "test",
					Sort: SortArgs{
						Email: "asc",
					},
				},
				mock: mockData{
					input: user.Filter{
						Name: "test",
						Sort: user.SortParams{
							Email: "asc",
						},
					},
					result: model.UserSlice{
						{
							ID:       2,
							Name:     "test",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       1,
							Name:     "test",
							Email:    "test3@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					totalCount: 2,
				},
				isCallToRepo: true,
			},
			expResult: []model.User{
				{
					ID:       2,
					Name:     "test",
					Email:    "test2@exam.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
				{
					ID:       1,
					Name:     "test",
					Email:    "test3@exam.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
			},
		},
		"success_multi_field_with_sort": {
			given: givenData{
				input: InputGetUser{
					Role:     "ADMIN",
					IsActive: null.Bool{Valid: true, Bool: false},
					Sort: SortArgs{
						Email:     "asc",
						CreatedAt: "desc",
					},
				},
				mock: mockData{
					input: user.Filter{
						Role:     "ADMIN",
						IsActive: null.Bool{Valid: true, Bool: false},
						Sort: user.SortParams{
							Email:     "asc",
							CreatedAt: "desc",
						},
					},
					result: model.UserSlice{
						{
							ID:       2,
							Name:     "test",
							Email:    "test2@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
						{
							ID:       1,
							Name:     "test",
							Email:    "test3@exam.com",
							Password: "123456",
							Phone:    "123456",
							Role:     "ADMIN",
							IsActive: true,
						},
					},
					totalCount: 2,
				},
				isCallToRepo: true,
			},
			expResult: []model.User{
				{
					ID:       2,
					Name:     "test",
					Email:    "test2@exam.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
				{
					ID:       1,
					Name:     "test",
					Email:    "test3@exam.com",
					Password: "123456",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			userRepoMock := new(user.Mock)
			if tc.given.isCallToRepo {
				userRepoMock.On("GetUsers", context.Background(), tc.given.mock.input).Return(tc.given.mock.result, tc.given.mock.totalCount, tc.given.mock.err)
			}
			repoMock := new(repository.Mock)
			repoMock.On("User").Return(userRepoMock)

			service := New(repoMock)

			// When
			result, _, err := service.GetUsers(context.Background(), tc.given.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, tc.expErr, "Should not be error")
				require.Equal(t, tc.expResult, result)
			}

			if tc.given.isCallToRepo {
				userRepoMock.AssertExpectations(t)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type input struct {
		updateUser         InputUser
		ctx                context.Context
		mockInputUser      mock.AnythingOfTypeArgument
		mockInputCTX       context.Context
		mockResultAffected int64
		mockResultError    error
	}
	tcs := map[string]struct {
		input  input
		expErr error
	}{
		"success": {
			input: input{
				updateUser: InputUser{
					ID:       1,
					Name:     "TEST",
					Email:    "test@example",
					Password: "test",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
				ctx: context.Background(),

				mockInputUser:      mock.AnythingOfType("User"),
				mockInputCTX:       context.Background(),
				mockResultAffected: 1,
			},
			expErr: nil,
		},
		"not_found": {
			input: input{
				updateUser: InputUser{
					ID:       2,
					Name:     "TEST",
					Email:    "test@example",
					Password: "test",
					Phone:    "123456",
					Role:     "ADMIN",
					IsActive: true,
				},
				ctx:                context.Background(),
				mockInputUser:      mock.AnythingOfType("User"),
				mockInputCTX:       context.Background(),
				mockResultAffected: 0,
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			userRepoMock := new(user.Mock)
			userRepoMock.On("UpdateUser", tc.input.mockInputCTX, tc.input.mockInputUser).Return(tc.input.mockResultAffected, tc.input.mockResultError)
			repoMock := new(repository.Mock)
			repoMock.On("User").Return(userRepoMock)

			service := New(repoMock)

			// When
			err := service.UpdateUser(tc.input.ctx, tc.input.updateUser)

			// Then
			if tc.expErr != nil {
				//must be error
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				//must be success
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	type mockData struct {
		userID int
		result model.User
		err    error
	}

	type givenData struct {
		userID int
		mock   mockData
	}

	tcs := map[string]struct {
		given     givenData
		expResult model.User
		expErr    error
	}{
		"success": {
			given: givenData{
				userID: 1,
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
			expResult: model.User{
				ID:       1,
				Name:     "test",
				Email:    "test@example.com",
				Password: "123456",
				Phone:    "123456",
				Role:     "GUEST",
				IsActive: true,
			},
		},
		"error": {
			given: givenData{
				userID: 1,
				mock: mockData{
					userID: 1,
					result: model.User{},
					err:    sql.ErrNoRows,
				},
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			userRepoMock := new(user.Mock)
			userRepoMock.On("GetUser", context.Background(), tc.given.mock.userID).Return(tc.given.mock.result, tc.given.mock.err)
			repoMock.On("User").Return(userRepoMock)

			userServ := New(repoMock)

			// When
			result, err := userServ.GetUser(context.Background(), tc.given.userID)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
			userRepoMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type mockData struct {
		userID  int
		rowsAff int64
		err     error
	}

	type givenData struct {
		userID int
		mock   mockData
	}

	tcs := map[string]struct {
		given  givenData
		expErr error
	}{
		"success": {
			given: givenData{
				userID: 1,
				mock: mockData{
					userID:  1,
					rowsAff: 1,
				},
			},
		},
		"error_no_rows_affected": {
			given: givenData{
				userID: 1,
				mock: mockData{
					userID:  1,
					rowsAff: 0,
				},
			},
			expErr: ErrUserNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			userRepoMock := new(user.Mock)
			userRepoMock.On("DeleteUser", context.Background(), tc.given.mock.userID).Return(tc.given.mock.rowsAff, tc.given.mock.err)
			repoMock.On("User").Return(userRepoMock)

			userServ := New(repoMock)

			// When
			err := userServ.DeleteUser(context.Background(), tc.given.userID)

			// Then
			if tc.expErr != nil {
				//must be error
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				//must be success
				require.NoError(t, err)
			}
			userRepoMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type input struct {
		ctx             context.Context
		loginInput      LoginInput
		mockInputCTX    context.Context
		mockInputEmail  string
		mockResultUser  model.User
		mockResultError error
	}
	type output struct {
		result LoginResponse
		err    error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx: context.Background(),
				loginInput: LoginInput{
					Email:    "example@example.com",
					Password: "123456789",
				},
				mockInputCTX:   context.Background(),
				mockInputEmail: "example@example.com",
				mockResultUser: model.User{
					ID:       1,
					Name:     "Guest",
					Email:    "example@example.com",
					Password: "$2a$14$R9cbWpV2ZjDxjvtWSiZ12OxKxJgpVePfeP8MpumxWr0yq614nKPeK",
					Phone:    "0987654321",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				result: LoginResponse{
					Scope:     "GUEST",
					ExpiresIn: tokenExpireTime,
					TokenType: "Bearer",
				},
			},
		},
		"email_is_not_registered": {
			input: input{
				ctx: context.Background(),
				loginInput: LoginInput{
					Email:    "example2@example.com",
					Password: "123456789",
				},
				mockInputCTX:    context.Background(),
				mockInputEmail:  "example2@example.com",
				mockResultError: sql.ErrNoRows,
			},
			expOutput: output{
				err: ErrEmailNotExist,
			},
		},
		"password_is_incorrect": {
			input: input{
				ctx: context.Background(),
				loginInput: LoginInput{
					Email:    "example@example.com",
					Password: "12345678910",
				},
				mockInputCTX:   context.Background(),
				mockInputEmail: "example@example.com",
				mockResultUser: model.User{
					ID:       1,
					Name:     "Guest",
					Email:    "example@example.com",
					Password: "$2a$14$R9cbWpV2ZjDxjvtWSiZ12OxKxJgpVePfeP8MpumxWr0yq614nKPeK",
					Phone:    "0987654321",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				err: ErrPasswordIncorrect,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			repoMock := new(repository.Mock)
			userRepoMock := new(user.Mock)
			userRepoMock.On("GetUserByEmail", tc.input.mockInputCTX, tc.input.mockInputEmail).Return(tc.input.mockResultUser, tc.input.mockResultError)
			repoMock.On("User").Return(userRepoMock)

			userServ := New(repoMock)

			// WHEN
			result, err := userServ.Login(tc.input.ctx, tc.input.loginInput)

			// THEN
			if err != nil {
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				tc.expOutput.result.AccessToken = result.AccessToken
				require.Equal(t, tc.expOutput.result, result)
			}
		})
	}
}

func TestStatisticsService_GetStatistics(t *testing.T) {
	type mockData struct {
		userSummary     user.SummaryStatistics
		userErr         error
		productSummary  product.SummaryStatistics
		productErr      error
		orderSummary    []order.Statistics
		orderErr        error
		latestOrders    []order.OrderInfo
		latestOrdersErr error
	}

	tcs := map[string]struct {
		mock      mockData
		expResult SummaryStatistics
		expErr    error
	}{
		"success": {
			mock: mockData{
				userSummary: user.SummaryStatistics{
					Total:         2,
					TotalInactive: 1,
				},
				productSummary: product.SummaryStatistics{
					Total:         2,
					TotalInactive: 0,
				},
				orderSummary: []order.Statistics{
					{
						Status: "NEW",
						Count:  2,
					},
					{
						Status: "PENDING",
						Count:  1,
					},
					{
						Status: "SUCCESS",
						Count:  1,
					},
					{
						Status: "FAILED",
						Count:  0,
					},
				},
				latestOrders: []order.OrderInfo{
					{
						OrderID:     1,
						OrderNumber: "TEST",
						OrderDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Status:      "NEW",
						UserID:      10,
						Total:       15000,
					},
				},
			},
			expResult: SummaryStatistics{
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
				LatestOrders: []OrderInfo{
					{
						OrderID:     1,
						OrderNumber: "TEST",
						OrderDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Status:      "NEW",
						UserID:      10,
						Total:       15000,
					},
				},
			},
		},
		"error": {
			mock: mockData{
				userSummary:    user.SummaryStatistics{},
				userErr:        fmt.Errorf("UNKNOWN ERROR"),
				productSummary: product.SummaryStatistics{},
				orderSummary:   []order.Statistics{},
				latestOrders:   []order.OrderInfo{},
			},
			expResult: SummaryStatistics{},
			expErr:    fmt.Errorf("UNKNOWN ERROR"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			userRepoMock := new(user.Mock)
			userRepoMock.On("GetStatistics", context.Background()).Return(tc.mock.userSummary, tc.mock.userErr)
			productRepoMock := new(product.Mock)
			productRepoMock.On("GetStatistics", context.Background()).Return(tc.mock.productSummary, tc.mock.productErr)
			orderRepoMock := new(order.Mock)
			orderRepoMock.On("GetStatistics", context.Background()).Return(tc.mock.orderSummary, tc.mock.orderErr)
			orderRepoMock.On("GetLatestOrder", context.Background(), 10).Return(tc.mock.latestOrders, tc.mock.latestOrdersErr)

			repoMock := new(repository.Mock)
			repoMock.On("User").Return(userRepoMock)
			repoMock.On("Product").Return(productRepoMock)
			repoMock.On("Order").Return(orderRepoMock)

			service := New(repoMock)

			// When
			result, err := service.GetStatistics(context.Background(), 10)

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
