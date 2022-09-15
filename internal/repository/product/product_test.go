package product

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/db"
)

func TestProductRepo_GetProduct(t *testing.T) {
	type input struct {
		productID     int
		givenDataPath string
		ctx           context.Context
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

				givenDataPath: "test_data/get_product.sql",
				ctx:           context.Background(),
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
			//GIVEN
			dbConn, err := db.DBConnect(os.Getenv("DB_URL"))
			productRepo := New(dbConn)
			require.NoError(t, err)
			db.LoadSqlTestFile(t, dbConn, tc.input.givenDataPath) // execute sql file to add data for test
			defer dbConn.Exec("delete from products; delete from users ;")

			//WHEN
			result, err := productRepo.GetProduct(tc.input.ctx, tc.input.productID)

			//THEN
			if tc.expOutput.err != nil {
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				tc.expOutput.product.CreatedAt = result.CreatedAt
				tc.expOutput.product.UpdatedAt = result.UpdatedAt

				require.NoError(t, err)
				require.Equal(t, tc.expOutput.product, result)
			}
		})
	}
}

func TestProductRepo_CreateProduct(t *testing.T) {
	type input struct {
		newProduct    model.Product
		ctx           context.Context
		givenDataPath string //path to sql file
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
				newProduct: model.Product{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/create_product.sql",
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
		"is_active_is_false": {
			input: input{
				newProduct: model.Product{
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    false,
					UserID:      1,
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/create_product.sql",
			},

			expOutput: output{
				product: model.Product{
					ID:          1,
					Title:       "test",
					Description: "",
					Price:       20000,
					Quantity:    10,
					IsActive:    false,
					UserID:      1,
				},
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			dbConn, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)
			db.LoadSqlTestFile(t, dbConn, tc.input.givenDataPath)
			defer dbConn.Exec("delete from products; delete from users ;")
			productRepo := New(dbConn)

			//WHEN
			result, err := productRepo.CreateProduct(tc.input.ctx, tc.input.newProduct)

			//THEN
			if tc.expOutput.err != nil {
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				tc.expOutput.product.ID = result.ID
				tc.expOutput.product.CreatedAt = result.CreatedAt
				tc.expOutput.product.UpdatedAt = result.UpdatedAt

				require.NoError(t, err)
				require.Equal(t, tc.expOutput.product, result)
			}
		})
	}
}

func TestProductRepo_UpdateProduct(t *testing.T) {
	tcs := map[string]struct {
		input     model.Product
		expResult int64
		expErr    error
	}{
		"success": {
			input: model.Product{
				ID:          1,
				Title:       "test",
				Description: "this is description",
				Price:       15000,
				Quantity:    9,
				IsActive:    true,
				UserID:      1,
			},
			expResult: 1,
			expErr:    nil,
		},
		"success_no_affected": {
			input: model.Product{
				ID:          10,
				Title:       "test",
				Description: "this is description",
				Price:       15000,
				Quantity:    9,
				IsActive:    true,
				UserID:      1,
			},
			expResult: 0,
			expErr:    nil,
		},
		"error_user_not_found": {
			input: model.Product{
				ID:          1,
				Title:       "test",
				Description: "this is description",
				Price:       15000,
				Quantity:    9,
				IsActive:    true,
				UserID:      100,
			},
			expErr: errors.New("model: unable to update products row: pq: insert or update on table \"products\" violates foreign key constraint \"products_user_id_fkey\""),
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			productRepo := New(dbTest)

			db.LoadSqlTestFile(t, dbTest, "test_data/products.sql")
			defer dbTest.Exec("delete from products; delete from users;")

			// When
			result, err := productRepo.UpdateProduct(context.Background(), tc.input)

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

func TestProductRepo_DeleteProduct(t *testing.T) {
	type input struct {
		productID     int
		ctx           context.Context
		givenDataPath string //path to sql file
	}
	type output struct {
		affectedRows int64
		err          error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				productID:     1,
				ctx:           context.Background(),
				givenDataPath: "test_data/delete_product.sql",
			},

			expOutput: output{
				affectedRows: 1,
			},
		},
		"not_found": {
			input: input{
				productID:     2,
				ctx:           context.Background(),
				givenDataPath: "test_data/delete_product.sql",
			},

			expOutput: output{
				affectedRows: 0,
				err:          nil,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			dbConn, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)
			db.LoadSqlTestFile(t, dbConn, tc.input.givenDataPath)
			defer dbConn.Exec("delete from users ;")
			defer dbConn.Exec("delete from products;")
			productRepo := New(dbConn)

			//WHEN
			result, err := productRepo.DeleteProduct(tc.input.ctx, tc.input.productID)

			//THEN
			if tc.expOutput.err != nil {
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.affectedRows, result)
			}
		})
	}
}

func TestProductRepo_GetProducts(t *testing.T) {
	type input struct {
		filter        Filter
		ctx           context.Context
		givenDataPath string //path to sql file
	}
	type output struct {
		products   []ProductItem
		totalCount int64
		err        error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				filter: Filter{
					PriceRange: PriceRange{
						MinPrice: 1000000, MaxPrice: 10000000,
					},
					Title:  "Gucci",
					UserID: 2,
					OrderBy: OrderBy{
						Price: "desc",
					},
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/get_products.sql",
			},
			expOutput: output{
				products: []ProductItem{
					{
						ID:          11,
						Title:       `Gucci trousers`,
						Description: `To be a gentleman`,
						Price:       10000000,
						Quantity:    9,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          9,
						Title:       `Gucci TShirt`,
						Description: `Fashion TShirt`,
						Price:       1000000,
						Quantity:    10,
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
		"sorted_data": {
			input: input{
				filter: Filter{
					UserID: 2,
					OrderBy: OrderBy{
						Price: "desc",
					},
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/get_products.sql",
			},
			expOutput: output{
				products: []ProductItem{
					{
						ID:          11,
						Title:       `Gucci trousers`,
						Description: `To be a gentleman`,
						Price:       10000000,
						Quantity:    9,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          9,
						Title:       `Gucci TShirt`,
						Description: `Fashion TShirt`,
						Price:       1000000,
						Quantity:    10,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          10,
						Title:       `Book : "Cha giau cha ngheo"`,
						Description: `Nice to read every weekend`,
						Price:       200000,
						Quantity:    8,
						IsActive:    true,
						User: CreatedBy{
							ID:    2,
							Name:  "admin 2",
							Email: "admin2@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          3,
						Title:       `Thien long pencil`,
						Description: `Nice pen from Thien Long company`,
						Price:       5000,
						Quantity:    3,
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
		"empty_data_by_is_active": {
			input: input{
				filter: Filter{
					IsActive: null.NewBool(false, true),
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/get_products.sql",
			},
			expOutput: output{
				products:   []ProductItem{},
				totalCount: 0,
			},
		},
		"success_with_pagination": {
			input: input{
				filter: Filter{
					IsActive: null.NewBool(true, true),
					Pagination: Pagination{
						Limit: 3,
						Page:  2,
					},
				},
				ctx:           context.Background(),
				givenDataPath: "test_data/get_products.sql",
			},
			expOutput: output{
				products: []ProductItem{
					{
						ID:          4,
						Title:       `Book : "Dac nhan tam"`,
						Description: `The favious book of the world`,
						Price:       100000,
						Quantity:    4,
						IsActive:    true,
						User: CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          5,
						Title:       `Note book`,
						Description: `It is small and pretty`,
						Price:       20000,
						Quantity:    5,
						IsActive:    true,
						User: CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
					{
						ID:          6,
						Title:       `Sneaker`,
						Description: `Beutyfull shoe`,
						Price:       5000000,
						Quantity:    10,
						IsActive:    true,
						User: CreatedBy{
							ID:    1,
							Name:  "admin",
							Email: "admin@example.com",
							Phone: "0987654321",
						},
					},
				},
				totalCount: 11,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			//GIVEN
			dbConn, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)
			db.LoadSqlTestFile(t, dbConn, tc.input.givenDataPath)
			defer dbConn.Exec("delete from users ;")
			defer dbConn.Exec("delete from products;")
			productRepo := New(dbConn)

			//WHEN
			result, totalCount, err := productRepo.GetProducts(tc.input.ctx, tc.input.filter)

			//THEN
			if tc.expOutput.err != nil {
				//must be error
				require.EqualError(t, tc.expOutput.err, err.Error())
			} else {
				// must be success
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.totalCount, totalCount)
				require.Equal(t, len(tc.expOutput.products), len(result))
				for i, v := range result {
					tc.expOutput.products[i].CreatedAt = v.CreatedAt
					tc.expOutput.products[i].UpdatedAt = v.UpdatedAt
					require.Equal(t, tc.expOutput.products[i], v)
				}
			}
		})
	}
}

func TestProductRepo_InsertALl(t *testing.T) {
	tcs := map[string]struct {
		given  []model.Product
		expErr error
	}{
		"success": {
			given: []model.Product{
				{
					Title:       "test1",
					Description: "This is record test 1",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      1,
				},
				{
					Title:       "test2",
					Description: "This is record test 2",
					Price:       5000,
					Quantity:    20,
					IsActive:    false,
					UserID:      1,
				},
				{
					Title:       "test3",
					Description: "This is record test 3",
					Price:       5,
					Quantity:    99,
					IsActive:    true,
					UserID:      1,
				},
			},
			expErr: nil,
		},
		"error_violate_foreign_key": {
			given: []model.Product{
				{
					Title:       "test1",
					Description: "This is record test 1",
					Price:       20000,
					Quantity:    10,
					IsActive:    true,
					UserID:      2,
				},
				{
					Title:       "test2",
					Description: "This is record test 2",
					Price:       5000,
					Quantity:    20,
					IsActive:    false,
					UserID:      2,
				},
				{
					Title:       "test3",
					Description: "This is record test 3",
					Price:       5,
					Quantity:    99,
					IsActive:    true,
					UserID:      2,
				},
			},
			expErr: errors.New("pq: insert or update on table \"products\" violates foreign key constraint \"products_user_id_fkey\""),
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)
			db.LoadSqlTestFile(t, dbTest, "test_data/products.sql")
			defer dbTest.Exec("delete from products; delete from users ;")
			txTest, txErr := dbTest.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
			require.NoError(t, txErr)

			productRepo := New(dbTest)

			// When
			err := productRepo.InsertAll(context.Background(), txTest, tc.given)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}

			err = txTest.Rollback()
			require.NoError(t, err)
		})
	}
}

func TestOrderRepository_GetStatistics(t *testing.T) {
	tcs := map[string]struct {
		expResult SummaryStatistics
		expErr    error
	}{
		"success": {
			expResult: SummaryStatistics{
				Total:         4,
				TotalInactive: 1,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			defer dbTest.Close()

			repo := New(dbTest)
			db.LoadSqlTestFile(t, dbTest, "test_data/products.sql")
			defer dbTest.Exec("DELETE FROM products;")

			// When
			result, err := repo.GetStatistics(context.Background())

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
