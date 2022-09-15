package user

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/db"
)

func TestRepository_CreateUser(t *testing.T) {
	tcs := map[string]struct {
		given     model.User
		expResult model.User
		expErr    error
	}{
		"success": {
			given: model.User{
				ID:       2,
				Name:     "test02",
				Email:    "test@example.com",
				Password: "123456",
				Phone:    "123456",
				Role:     "GUEST",
				IsActive: true,
			},
			expResult: model.User{
				ID:       2,
				Name:     "test02",
				Email:    "test@example.com",
				Password: "123456",
				Phone:    "123456",
				Role:     "GUEST",
				IsActive: true,
			},
		},
		"error_duplicate_email": {
			given: model.User{
				Name:     "test2",
				Email:    "test2@example.com",
				Password: "123456",
				Phone:    "123456",
				Role:     "GUEST",
				IsActive: true,
			},
			expErr: errors.New("model: unable to insert into users: pq: duplicate key value violates unique constraint \"email_on_users\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			db.LoadSqlTestFile(t, dbTest, "test_data/users.sql")
			defer dbTest.Exec("DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.CreateUser(context.Background(), tc.given)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)

				tc.expResult.ID = result.ID
				tc.expResult.CreatedAt = result.CreatedAt
				tc.expResult.UpdatedAt = result.UpdatedAt

				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestRepository_ExistsUserByEmail(t *testing.T) {
	tcs := map[string]struct {
		given     string
		expResult bool
		expErr    error
	}{
		"success01": {
			given:     "test1@example.com",
			expResult: true,
		},
		"success02": {
			given:     "test12@example.com",
			expResult: false,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)

			db.LoadSqlTestFile(t, dbTest, "test_data/users.sql")
			defer dbTest.Exec("DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.ExistsUserByEmail(context.Background(), tc.given)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestUserRepository_GetUsers(t *testing.T) {
	tcs := map[string]struct {
		given         Filter
		expResult     []model.User
		expTotalCount int64
		expErr        error
	}{
		"success_field": {
			given: Filter{
				Name: "test1",
			},
			expResult: []model.User{
				{
					ID:   10,
					Name: "test1", Email: "test1@example.com", Password: "test", Phone: "test", Role: "ADMIN", IsActive: true,
				},
			},
			expTotalCount: 1,
		},
		"success_field_with_sort": {
			given: Filter{
				Name: "test",
			},
			expResult: []model.User{
				{
					ID:   10,
					Name: "test1", Email: "test1@example.com", Password: "test", Phone: "test", Role: "ADMIN", IsActive: true,
				},
				{
					ID:   11,
					Name: "test2", Email: "test2@example.com", Password: "test", Phone: "test", Role: "ADMIN", IsActive: false,
				},
				{
					ID:   12,
					Name: "test3", Email: "test3@example.com", Password: "test", Phone: "test", Role: "GUEST", IsActive: true,
				},
				{
					ID:   13,
					Name: "test4", Email: "test4@example.com", Password: "test", Phone: "test", Role: "GUEST", IsActive: true,
				},
			},
			expTotalCount: 4,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, err := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, err)
			repo := New(dbTest)

			db.LoadSqlTestFile(t, dbTest, "test_data/users.sql")
			defer dbTest.Exec("DELETE FROM users;")

			// When
			result, totalCount, err := repo.GetUsers(context.Background(), tc.given)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)

				for i, user := range result {
					// Ignore id, created_at, updated_at fields
					tc.expResult[i].ID = user.ID
					tc.expResult[i].CreatedAt = user.CreatedAt
					tc.expResult[i].UpdatedAt = user.UpdatedAt

					require.Equal(t, tc.expResult[i], *user)
				}

				require.Equal(t, tc.expTotalCount, totalCount)
			}
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	type input struct {
		ctx  context.Context
		user model.User
	}
	type output struct {
		expErr    error
		expResult int64 // affected Rows
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx: context.Background(),
				user: model.User{
					ID:       10,
					Name:     "test2",
					Email:    "test3@example.com",
					Password: "test",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				expResult: 1,
			},
		},
		"not_found": {
			input: input{
				ctx: context.Background(),
				user: model.User{
					ID:       12,
					Name:     "test2",
					Email:    "test3@example.com",
					Password: "test",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				expResult: 0,
			},
		},
		"email_duplicated": {
			input: input{
				ctx: context.Background(),
				user: model.User{
					ID:       10,
					Name:     "test2",
					Email:    "test2@example.com",
					Password: "test",
					Phone:    "123456",
					Role:     "GUEST",
					IsActive: true,
				},
			},
			expOutput: output{
				expResult: 0,
				expErr:    errors.New("model: unable to update users row: pq: duplicate key value violates unique constraint \"email_on_users\""),
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			db.LoadSqlTestFile(t, dbTest, "test_data/update_user.sql")
			defer dbTest.Exec("DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.UpdateUser(context.Background(), tc.input.user)

			// Then
			if tc.expOutput.expErr != nil {
				require.EqualError(t, err, tc.expOutput.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expOutput.expResult, result)
			}
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	tcs := map[string]struct {
		given     int
		expResult model.User
		expErr    error
	}{
		"success": {
			given: 10,
			expResult: model.User{
				ID:       10,
				Name:     "test1",
				Email:    "test1@example.com",
				Password: "test",
				Phone:    "test",
				Role:     "ADMIN",
				IsActive: true,
			},
		},
		"error_not_found": {
			given:  15,
			expErr: sql.ErrNoRows,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			db.LoadSqlTestFile(t, dbTest, "test_data/users.sql")
			defer dbTest.Exec("DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.GetUser(context.Background(), tc.given)

			// Then
			if tc.expErr != nil {
				//must be error
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				//must be success
				require.NoError(t, err)
				tc.expResult.CreatedAt = result.CreatedAt
				tc.expResult.UpdatedAt = result.UpdatedAt
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	tcs := map[string]struct {
		given   int
		rowsAff int64
		expErr  error
	}{
		"success": {
			given:   1,
			rowsAff: 1,
		},
		"error": {
			given:  2,
			expErr: errors.New("model: unable to delete all from users: pq: update or delete on table \"users\" violates foreign key constraint \"products_user_id_fkey\" on table \"products\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			db.LoadSqlTestFile(t, dbTest, "test_data/delete_user.sql")
			defer dbTest.Exec("DELETE FROM products; DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.DeleteUser(context.Background(), tc.given)

			// Then
			if tc.expErr != nil {
				//must be error
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				//must be success
				require.NoError(t, err)
				require.Equal(t, tc.rowsAff, result)
			}
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	type input struct {
		ctx   context.Context
		email string
	}
	type output struct {
		user model.User
		err  error
	}
	tcs := map[string]struct {
		input     input
		expOutput output
	}{
		"success": {
			input: input{
				ctx:   context.Background(),
				email: "mai@example.com",
			},
			expOutput: output{
				user: model.User{
					ID:       10,
					Name:     "Mai",
					Email:    "mai@example.com",
					Password: "test",
					Phone:    "test",
					Role:     "ADMIN",
					IsActive: true,
				},
			},
		},
		"not_found": {
			input: input{
				ctx:   context.Background(),
				email: "mai2@example.com",
			},
			expOutput: output{
				err: sql.ErrNoRows,
			},
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			dbTest, dbErr := db.DBConnect(os.Getenv("DB_URL"))
			require.NoError(t, dbErr)

			db.LoadSqlTestFile(t, dbTest, "test_data/get_user_by_email.sql")
			defer dbTest.Exec("DELETE FROM products;DELETE FROM users;")

			repo := New(dbTest)

			// When
			result, err := repo.GetUserByEmail(tc.input.ctx, tc.input.email)

			// Then
			if tc.expOutput.err != nil {
				//must be error
				require.EqualError(t, err, tc.expOutput.err.Error())
			} else {
				//must be success
				require.NoError(t, err)
				// Ignore createdAt and updatedAt fields
				tc.expOutput.user.CreatedAt = result.CreatedAt
				tc.expOutput.user.UpdatedAt = result.UpdatedAt
				require.Equal(t, tc.expOutput.user, result)
			}
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
			db.LoadSqlTestFile(t, dbTest, "test_data/users.sql")
			defer dbTest.Exec("DELETE FROM users;")

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
