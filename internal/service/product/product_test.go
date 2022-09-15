package product

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
)

func TestProductService_GetProduct(t *testing.T) {
	type input struct {
		productID         int
		ctx               context.Context
		mockInputID       int
		mockResultProduct model.Product
		mockResultError   error
	}
	type output struct {
		product model.Product
		err     error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				productID: 1,
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

				ctx: context.Background(),
			},
			expOutput: output{
				product: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			productMock := new(product.Mock)
			productMock.On("GetProduct", context.Background(), tc.input.productID).Return(tc.input.mockResultProduct, tc.input.mockResultError)
			repoMock := new(repository.Mock)
			repoMock.On("Product").Return(productMock)
			productService := New(repoMock)

			// WHEN
			result, err := productService.GetProduct(tc.input.ctx, tc.input.productID)

			//THEN
			if tc.expOutput.err != nil {
				//must be error
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				// must be success
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.product, result)
			}
			productMock.AssertExpectations(t)
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	type mockCreateProduct struct {
		ctx     context.Context
		product model.Product
		result  model.Product
		err     error
	}
	type mockExistUser struct {
		ctx    context.Context
		userID int
		result bool
		err    error
	}

	type input struct {
		product           ProductInput
		ctx               context.Context
		mockCreateProduct mockCreateProduct
		mockExistUser     mockExistUser
	}

	type output struct {
		product model.Product
		err     error
	}

	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				product: ProductInput{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				ctx: context.Background(),
				mockCreateProduct: mockCreateProduct{
					ctx: context.Background(),
					product: model.Product{
						Title:       "test",
						Description: "",
						Price:       20000,
						Quantity:    10,
						IsActive:    true,
						UserID:      1,
					},
					result: model.Product{
						ID:          1,
						Title:       "test",
						Description: "",
						Price:       20000,
						Quantity:    10,
						IsActive:    true,
						UserID:      1,
					},
				},
				mockExistUser: mockExistUser{
					ctx:    context.Background(),
					userID: 1,
					result: true,
				},
			},
			expOutput: output{
				product: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
			},
		},
		"user_is_not_exist": {
			input: input{
				product: ProductInput{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      100,
				},
				ctx: context.Background(),
				mockExistUser: mockExistUser{
					ctx:    context.Background(),
					userID: 100,
					result: false,
				},
			},
			expOutput: output{
				err: ErrUserNotExist,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			productMock := new(product.Mock)
			userMock := new(user.Mock)
			productMock.On("CreateProduct", tc.input.mockCreateProduct.ctx, tc.input.mockCreateProduct.product).Return(tc.input.mockCreateProduct.result, tc.input.mockCreateProduct.err)
			userMock.On("ExistsUserByID", tc.input.mockExistUser.ctx, tc.input.mockExistUser.userID).Return(tc.input.mockExistUser.result, tc.input.mockExistUser.err)
			repoMock := new(repository.Mock)
			repoMock.On("Product").Return(productMock)
			repoMock.On("User").Return(userMock)

			productService := New(repoMock)

			// WHEN
			result, err := productService.CreateProduct(tc.input.ctx, tc.input.product)

			//THEN
			if tc.expOutput.err != nil {
				//must be error
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				// must be success
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.product, result)
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	type mockData struct {
		input       model.Product
		affectedRow int64
		err         error
	}

	type givenData struct {
		id    int
		input ProductInput
		mock  mockData
	}

	tcs := map[string]struct {
		given  givenData
		expErr error
	}{
		"success": {
			given: givenData{
				id: 1,
				input: ProductInput{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				mock: mockData{
					input: model.Product{
						ID:          1,
						Title:       "test",
						Description: "",
						Price:       20000,
						Quantity:    10,
						IsActive:    true,
						UserID:      1,
					},
					affectedRow: 1,
				},
			},
			expErr: nil,
		},
		"error": {
			given: givenData{
				id: 1,
				input: ProductInput{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				mock: mockData{
					input: model.Product{
						ID:          1,
						Title:       "test",
						Description: "",
						Price:       20000,
						Quantity:    10,
						IsActive:    true,
						UserID:      1,
					},
					affectedRow: 0,
				},
			},
			expErr: ErrProductNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			productMock := new(product.Mock)
			productMock.On("UpdateProduct", context.Background(), tc.given.mock.input).Return(tc.given.mock.affectedRow, tc.given.mock.err)
			productMock.On("ExistsProductByID", context.Background(), tc.given.id).Return(true, nil)
			repoMock.On("Product").Return(productMock)

			productService := New(repoMock)

			// When
			err := productService.UpdateProduct(context.Background(), tc.given.id, tc.given.input)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	// type struct Input for DeleteProduct func, input/output for mock func
	type input struct {
		ctx                context.Context
		productID          int
		mockInputID        int
		mockInputCTX       context.Context
		mockOutputError    error
		mockOutputAffected int
	}
	type output struct {
		product model.Product
		err     error
	}
	tcs := map[string]struct {
		input    input
		expError error // output
	}{
		"success": {
			input: input{
				ctx:       context.Background(),
				productID: 1,

				mockInputID:        1,
				mockInputCTX:       context.Background(),
				mockOutputAffected: 1,
			},
			expError: nil,
		},
		"not_found": {
			input: input{
				ctx:       context.Background(),
				productID: 2,

				mockInputID:        2,
				mockInputCTX:       context.Background(),
				mockOutputAffected: 0,
			},
			expError: ErrProductNotFound,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			productMock := new(product.Mock)
			productMock.On("DeleteProduct", tc.input.mockInputCTX, tc.input.productID).Return(tc.input.mockOutputAffected, tc.input.mockOutputError)
			repoMock := new(repository.Mock)
			repoMock.On("Product").Return(productMock)

			productService := New(repoMock)

			//WHEN
			err := productService.DeleteProduct(tc.input.ctx, tc.input.productID)

			//THEN
			if tc.expError != nil {
				//must be error
				require.EqualError(t, err, tc.expError.Error())
			} else {
				//must be success
				require.NoError(t, err)
			}
			productMock.AssertExpectations(t)
		})
	}
}

func TestProductService_GetProducts(t *testing.T) {
	type input struct {
		ctx                  context.Context
		getProductsInput     GetProductsInput
		mockInputCTX         context.Context
		mockInputFilter      product.Filter
		mockResultProducts   []product.ProductItem
		mockResultError      error
		mockResultTotalCount int64
	}
	type output struct {
		err        error
		result     []ProductItem
		totalCount int64
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx: context.Background(),
				getProductsInput: GetProductsInput{
					Title: "success",
					PriceRange: PriceRange{
						MinPrice: 0,
						MaxPrice: 100000,
					},
					IsActive: null.NewBool(true, true),
					UserID:   2,
					OrderBy: OrderInput{
						Title:     "desc",
						CreatedAt: "desc",
					},
					Pagination: Pagination{
						Limit: 20,
						Page:  1,
					},
				},
				mockInputCTX: context.Background(),
				mockInputFilter: product.Filter{
					Title: "success",
					PriceRange: product.PriceRange{
						MinPrice: 0,
						MaxPrice: 100000,
					},
					IsActive: null.NewBool(true, true),
					UserID:   2,
					OrderBy: product.OrderBy{
						Title:     "desc",
						CreatedAt: "desc",
					},
					Pagination: product.Pagination{
						Limit: 20,
						Page:  1,
					},
				},
				mockResultProducts: []product.ProductItem{
					{
						ID:          1,
						Title:       "success",
						Description: "success",
						Price:       100,
						Quantity:    15,
						IsActive:    true,
						User: product.CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          2,
						Title:       "success 2",
						Description: "success",
						Price:       10000,
						Quantity:    150,
						IsActive:    true,
						User: product.CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
				},
				mockResultTotalCount: 2,
			},
			expOutput: output{
				result: []ProductItem{
					{
						ID:          1,
						Title:       "success",
						Description: "success",
						Price:       100,
						Quantity:    15,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          2,
						Title:       "success 2",
						Description: "success",
						Price:       10000,
						Quantity:    150,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
				},
				totalCount: 2,
			},
		},
		"empty_result": {
			input: input{
				ctx: context.Background(),
				getProductsInput: GetProductsInput{
					Title: "empty result",
					PriceRange: PriceRange{
						MinPrice: 10000,
						MaxPrice: 100000,
					},
					IsActive: null.NewBool(false, true),
					UserID:   1,
					Pagination: Pagination{
						Limit: 20,
						Page:  1,
					},
				},
				mockInputCTX: context.Background(),
				mockInputFilter: product.Filter{
					Title: "empty result",
					PriceRange: product.PriceRange{
						MinPrice: 10000,
						MaxPrice: 100000,
					},
					IsActive: null.NewBool(false, true),
					UserID:   1,
					Pagination: product.Pagination{
						Limit: 20,
						Page:  1,
					},
				},
				mockResultProducts: []product.ProductItem{},
			},
			expOutput: output{
				result:     []ProductItem{},
				totalCount: 0,
			},
		},
		"paging_result": {
			input: input{
				ctx: context.Background(),
				getProductsInput: GetProductsInput{
					Pagination: Pagination{
						Limit: 2,
						Page:  2,
					},
				},
				mockInputCTX: context.Background(),
				mockInputFilter: product.Filter{
					Pagination: product.Pagination{
						Limit: 2,
						Page:  2,
					},
				},
				mockResultProducts: []product.ProductItem{
					{
						ID:          3,
						Title:       "Macbook 11",
						Description: "None",
						Price:       2500000,
						Quantity:    150,
						IsActive:    true,
						User: product.CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          4,
						Title:       "Macbook 12",
						Description: "None",
						Price:       3000000,
						Quantity:    100,
						IsActive:    true,
						User: product.CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
				},
				mockResultTotalCount: 4,
			},
			expOutput: output{
				result: []ProductItem{
					{
						ID:          3,
						Title:       "Macbook 11",
						Description: "None",
						Price:       2500000,
						Quantity:    150,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          4,
						Title:       "Macbook 12",
						Description: "None",
						Price:       3000000,
						Quantity:    100,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
				},
				totalCount: 4,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			productMock := new(product.Mock)
			productMock.On("GetProducts", tc.input.mockInputCTX, tc.input.mockInputFilter).Return(tc.input.mockResultProducts, tc.input.mockResultTotalCount, tc.input.mockResultError)
			repoMock := new(repository.Mock)
			repoMock.On("Product").Return(productMock)

			productService := New(repoMock)

			// WHEN
			result, totalCount, err := productService.GetProducts(tc.input.ctx, tc.input.getProductsInput)

			// THEN
			if err != nil {
				// must be error
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				// must be success
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.totalCount, totalCount)
				require.Equal(t, tc.expOutput.result, result)
			}
		})
	}
}

func TestProductService_ImportProductCSV(t *testing.T) {
	type mockData struct {
		fn    mock.AnythingOfTypeArgument
		txErr error
	}

	type givenData struct {
		csvData string
		mock    mockData
	}

	tcs := map[string]struct {
		given  givenData
		expErr error
	}{
		"success": {
			given: givenData{
				csvData: "Title,Description,Price,Quantity,Is_Active,User_ID\nMahalia,rBIsqaccTvxl,462.25,537,1,1\nBerget,jVLHPRgDbMHNeds,910.59,913,1,1\nDede, vNlZDZd YtaB,535.69,356,0,1",
				mock: mockData{
					fn:    mock.AnythingOfType("func(*sql.Tx) error"),
					txErr: nil,
				},
			},
			expErr: nil,
		},
		"error_when_send_null_file": {
			given: givenData{
				csvData: "",
				mock: mockData{
					fn:    mock.AnythingOfType("func(*sql.Tx) error"),
					txErr: nil,
				},
			},
			expErr: io.EOF,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			repoMock := new(repository.Mock)
			repoMock.On("Tx", context.Background(), tc.given.mock.fn).Return(tc.given.mock.txErr)
			productMock := new(product.Mock)
			repoMock.On("Product").Return(productMock)

			productServ := New(repoMock)
			csvReader := strings.NewReader(tc.given.csvData)

			// When
			err := productServ.ImportProductCSV(context.Background(), "product.csv", csvReader)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProductService_ExportProductsToCSV(t *testing.T) {
	type input struct {
		ctx                context.Context
		getProductsInput   GetProductsInput
		dir                string
		mockInputCTX       context.Context
		mockInputFilter    product.Filter
		mockResultErr      error
		mockResultProducts []product.ProductItem
	}
	type output struct {
		data string
		err  error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx:              context.Background(),
				getProductsInput: GetProductsInput{},
				dir:              "",
				mockInputCTX:     context.Background(),
				mockInputFilter: product.Filter{
					Pagination: product.Pagination{
						Page:  1,
						Limit: 1000,
					},
				},
				mockResultProducts: []product.ProductItem{
					{
						ID:          1,
						Title:       "Test",
						Description: "test",
						Price:       100000,
						Quantity:    20,
						IsActive:    true,
						User: product.CreatedBy{
							Name: "Mai",
						},
					},
				},
			},
			expOutput: output{
				data: `ID,Title,Description,Price,Quantity,Activated,Created By,Created Date,Updated Date
				1,Test,test,100000.00,20,true,Mai,0001-01-01 00:00:00,0001-01-01 00:00:00
				`,
			},
		},
		"cannot_create_file": {
			input: input{
				ctx:              context.Background(),
				getProductsInput: GetProductsInput{},
				dir:              "files",
				mockInputCTX:     context.Background(),
				mockInputFilter: product.Filter{
					Pagination: product.Pagination{
						Page:  1,
						Limit: 1000,
					},
				},
				mockResultProducts: []product.ProductItem{
					{
						ID:          1,
						Title:       "Test",
						Description: "test",
						Price:       100000,
						Quantity:    20,
						IsActive:    true,
						User: product.CreatedBy{
							Name: "Mai",
						},
					},
				},
			},
			expOutput: output{
				data: `ID,Title,Description,Price,Quantity,Activated,Created By,Created Date,Updated Date
				1,Test,test,100000.00,20,true,Mai,0001-01-01 00:00:00,0001-01-01 00:00:00
				`,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			productRepoMock := new(product.Mock)
			productRepoMock.On("GetProducts", tc.input.mockInputCTX, tc.input.mockInputFilter).Return(tc.input.mockResultProducts, int64(0), tc.input.mockResultErr)
			repoMock := new(repository.Mock)
			repoMock.On("Product").Return(productRepoMock)

			productService := New(repoMock)

			//WHEN
			result, err := productService.ExportProductsCSV(tc.input.mockInputCTX, tc.input.getProductsInput, "")
			defer func() {
				if _, err := os.Stat(result); err != os.ErrNotExist {
					os.Remove(result)
				}
			}()
			//THEN
			if err != nil {
				//must be error
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				require.NoError(t, err)
				require.FileExists(t, result)
				result, err := os.ReadFile(result)
				require.NoError(t, err)

				expResult := strings.ReplaceAll(tc.expOutput.data, "\t", "")
				require.Equal(t, expResult, string(result))
			}
		})
	}
}

func TestProductService_GetCSVFile(t *testing.T) {
	type input struct {
		ctx       context.Context
		filePath  string
		givenFile string
	}
	type output struct {
		file string
		err  error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx:       context.Background(),
				filePath:  "file_test.csv",
				givenFile: "ID,Title,Description,Price,Quantity,Activated,Created By,Created Date,Updated Date\n1,Test,test,100000.00,20,true,Mai,0001-01-01 00:00:00,0001-01-01 00:00:00",
			},
			expOutput: output{
				file: "ID,Title,Description,Price,Quantity,Activated,Created By,Created Date,Updated Date\n1,Test,test,100000.00,20,true,Mai,0001-01-01 00:00:00,0001-01-01 00:00:00",
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			productService := New(nil)
			file, err := os.Create(tc.input.filePath) // Create file to test
			require.NoError(t, err)
			defer func() {
				if _, err := os.Stat(tc.input.filePath); err != os.ErrNotExist {
					os.Remove(tc.input.filePath)
				}
			}()
			file.Write([]byte(tc.input.givenFile))

			//WHEN
			result, err := productService.DownloadCSV(tc.input.ctx, tc.input.filePath)

			//THEN
			if err != nil {
				//must be error
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.file, string(result))
			}
		})
	}
}
