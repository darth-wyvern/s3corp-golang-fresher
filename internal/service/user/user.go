package user

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository/user"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/bcrypt"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/jwt"
)

// InputUser is the input for v1
type InputUser struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone    string
	Role     string
	IsActive bool
}

// CreateUser creates a new user by InputUser param.
func (serv impl) CreateUser(ctx context.Context, input InputUser) (model.User, error) {
	// 1. Check exist user with this email
	existed, err := serv.repo.User().ExistsUserByEmail(ctx, input.Email)
	if err != nil {
		return model.User{}, err
	}

	if existed {
		return model.User{}, ErrEmailExisted
	}

	// 2. Hash user password by bcrypt
	hashedPass, err := bcrypt.HashPassword(input.Password)
	if err != nil {
		return model.User{}, ErrPasswordCannotBeHashed
	}

	// 3. Create user
	result, err := serv.repo.User().CreateUser(ctx, model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPass,
		Phone:    input.Phone,
		Role:     input.Role,
		IsActive: input.IsActive,
	})
	if err != nil {
		return model.User{}, err
	}

	return result, nil
}

type SortArgs struct {
	Name      string
	Email     string
	CreatedAt string
}

type Pagination struct {
	Page, Limit int
}

type InputGetUser struct {
	ID         int
	Email      string
	Name       string
	IsActive   null.Bool
	Role       string
	Sort       SortArgs
	Pagination Pagination
}

// toFilter converts InputGetUser to Filter
func toFilter(input InputGetUser) user.Filter {
	sortParams := user.SortParams{
		Name:      input.Sort.Name,
		Email:     input.Sort.Email,
		CreatedAt: input.Sort.CreatedAt,
	}
	pagination := user.Pagination{
		Page:  input.Pagination.Page,
		Limit: input.Pagination.Limit,
	}

	return user.Filter{
		ID:         input.ID,
		Email:      input.Email,
		Name:       input.Name,
		IsActive:   input.IsActive,
		Role:       input.Role,
		Sort:       sortParams,
		Pagination: pagination,
	}
}

// GetUsers returns list of user by InputGetUser param.
func (serv impl) GetUsers(ctx context.Context, input InputGetUser) ([]model.User, int64, error) {
	// 1. Init filter
	filterData := toFilter(input)

	// 2. Call repo to get list of users.
	userSlice, totalCount, err := serv.repo.User().GetUsers(ctx, filterData)
	if err != nil {
		return []model.User{}, 0, err
	}

	// 3. Convert userSlice to list of model.User
	users := make([]model.User, len(userSlice))
	for i := 0; i < len(userSlice); i++ {
		users[i] = *userSlice[i]
	}

	return users, totalCount, nil
}

func (serv impl) UpdateUser(ctx context.Context, input InputUser) error {
	// TODO: check if email is existed

	// Hash new password
	hashedNewPass, err := bcrypt.HashPassword(input.Password)
	if err != nil {
		return ErrPasswordCannotBeHashed
	}

	//Call the repository method
	result, err := serv.repo.User().UpdateUser(ctx,
		model.User{
			ID:       input.ID,
			Name:     input.Name,
			Email:    input.Email,
			Password: hashedNewPass,
			Phone:    input.Phone,
			Role:     input.Role,
			IsActive: input.IsActive,
		},
	)
	if err != nil {
		return err
	}
	// Return not_found error if the affectedRows =1
	if result < 1 {
		return ErrUserNotFound
	}
	return nil
}

func (serv impl) GetUser(ctx context.Context, id int) (model.User, error) {
	// 1. Call repo to get user.
	result, err := serv.repo.User().GetUser(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	} else if err != nil {
		return model.User{}, err
	}

	return result, nil
}

func (serv impl) DeleteUser(ctx context.Context, id int) error {
	//Call the repository method
	result, err := serv.repo.User().DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Return not_found error if the affectedRows =1
	if result < 1 {
		return ErrUserNotFound
	}
	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken string        `json:"access_token"`
	Scope       string        `json:"scope"`
	ExpiresIn   time.Duration `json:"expires_in"`
	TokenType   string        `json:"token_type"`
}

const (
	tokenExpireTime = 30 * time.Minute
)

// Login authenticate user data
func (serv impl) Login(ctx context.Context, input LoginInput) (LoginResponse, error) {
	// Get user with email
	user, err := serv.repo.User().GetUserByEmail(ctx, input.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return LoginResponse{}, ErrEmailNotExist
	} else if err != nil {
		return LoginResponse{}, err
	}

	// Verify password
	if !bcrypt.CheckPasswordHash(input.Password, user.Password) {
		return LoginResponse{}, ErrPasswordIncorrect
	}

	// Generate access_token
	secretKey := os.Getenv("ACCESS_TOKEN_KEY")

	jwtInput := jwt.JWTInput{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		SecretKey: secretKey,
		ExpiresIn: tokenExpireTime,
	}

	_, accessToken, err := jwt.GenerateJWTToken(jwtInput)
	if err != nil {
		return LoginResponse{}, ErrTokeCannotBeGenerated
	}

	return LoginResponse{
		AccessToken: accessToken,
		Scope:       user.Role,
		ExpiresIn:   tokenExpireTime,
		TokenType:   "Bearer",
	}, nil
}
