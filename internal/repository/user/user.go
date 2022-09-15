package user

import (
	"context"
	"fmt"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

// CreateUser creates a new user.
func (r impl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	if err := user.Insert(ctx, r.db, boil.Whitelist("name", "email", "password", "phone", "role", "is_active", "created_at", "updated_at")); err != nil {
		return model.User{}, err
	}
	return user, nil
}

// ExistsUserByEmail checks if a user exists by email.
func (r impl) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	return model.Users(model.UserWhere.Email.EQ(email)).Exists(ctx, r.db)
}

func (r impl) ExistsUserByID(ctx context.Context, id int) (bool, error) {
	return model.Users(model.UserWhere.ID.EQ(id)).Exists(ctx, r.db)
}

type SortParams struct {
	Name      string
	Email     string
	CreatedAt string
}

type Pagination struct {
	Page  int
	Limit int
}

type Filter struct {
	ID         int
	Email      string
	Name       string
	IsActive   null.Bool
	Role       string
	Sort       SortParams
	Pagination Pagination
}

// GetUsers returns a list of users by filter.
func (r impl) GetUsers(ctx context.Context, input Filter) (model.UserSlice, int64, error) {
	// 1. Init query mods slice.
	var qms []qm.QueryMod

	// 2. Add filter condition.
	if input.ID > 0 {
		qms = append(qms, model.UserWhere.ID.EQ(input.ID))
	}
	if input.Email != "" {
		qms = append(qms, model.UserWhere.Email.EQ(input.Email))
	}
	if input.Name != "" {
		qms = append(qms, qm.Where("name LIKE ?", "%"+input.Name+"%"))
	}
	if input.IsActive.Valid {
		qms = append(qms, model.UserWhere.IsActive.EQ(input.IsActive.Bool))
	}
	if input.Role != "" {
		qms = append(qms, model.UserWhere.Role.EQ(input.Role))
	}

	// 3. Calculate total rows of filtered users list.
	totalCount, err := model.Users(qms...).Count(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	// 4. Add sort condition.
	if input.Sort != (SortParams{}) {
		if input.Sort.Name != "" {
			qms = append(qms, qm.OrderBy(fmt.Sprintf("%s %s", model.UserColumns.Name, input.Sort.Name)))
		}
		if input.Sort.Email != "" {
			qms = append(qms, qm.OrderBy(fmt.Sprintf("%s %s", model.UserColumns.Email, input.Sort.Email)))
		}
		if input.Sort.CreatedAt != "" {
			qms = append(qms, qm.OrderBy(fmt.Sprintf("%s %s", model.UserColumns.CreatedAt, input.Sort.CreatedAt)))
		}
	} else {
		qms = append(qms, qm.OrderBy(fmt.Sprintf("%s %s", model.UserColumns.UpdatedAt, "desc")))
	}

	// 5. Add pagination condition.
	if input.Pagination != (Pagination{}) {
		qms = append(
			qms,
			qm.Offset(input.Pagination.Limit*(input.Pagination.Page-1)), // Example: Pagination is 2, offset is 10. So offset is 10.
			qm.Limit(input.Pagination.Limit),
		)
	} else {
		qms = append(qms, qm.Offset(0), qm.Limit(20)) // Default pagination is 0, limit is 20.
	}

	// 6. Get users by query mods.
	users, err := model.Users(qms...).All(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// UpdateUser updates the user profile and returns the updated user profile
func (r impl) UpdateUser(ctx context.Context, updateUser model.User) (int64, error) {
	// Update the user profile
	result, err := updateUser.Update(context.Background(), r.db, boil.Infer())

	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r impl) GetUser(ctx context.Context, id int) (model.User, error) {
	user, err := model.Users(model.UserWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}

func (r impl) DeleteUser(ctx context.Context, id int) (int64, error) {
	result, err := model.Users(model.UserWhere.ID.EQ(id)).DeleteAll(ctx, r.db)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetUserByEmail returns the user with the given email
func (r impl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	// Get the user by email
	result, err := model.Users(model.UserWhere.Email.EQ(email)).One(ctx, r.db)
	if err != nil {
		return model.User{}, err
	}
	return *result, nil
}

type SummaryStatistics struct {
	Total         int64
	TotalInactive int64
}

func (r impl) GetStatistics(ctx context.Context) (SummaryStatistics, error) {
	// Get the total number of users
	total, err := model.Users().Count(ctx, r.db)
	if err != nil {
		return SummaryStatistics{}, err
	}

	// Get the total number of inactive users
	totalInactive, err := model.Users(model.UserWhere.IsActive.EQ(false)).Count(ctx, r.db)
	if err != nil {
		return SummaryStatistics{}, err
	}
	return SummaryStatistics{
		Total:         total,
		TotalInactive: totalInactive,
	}, nil
}
