package v1

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/null/v8"

	userServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

const (
	UserRoleAdmin = "ADMIN"
	UserRoleGuest = "GUEST"
)

type userRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Role     string `json:"role" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

func isValidRole(role string) bool {
	return role == UserRoleAdmin || role == UserRoleGuest
}

func validateUserID(id string) (int, error) {
	userID, err := strconv.Atoi(id)
	if err != nil || userID < 0 {
		return 0, ErrInvalidUserID
	}
	return userID, nil
}

func validateUserInput(req userRequest) (userServ.InputUser, error) {
	if _, err := mail.ParseAddress(req.Email); err != nil { // parsed without error means valid email
		return userServ.InputUser{}, ErrInvalidEmail
	}
	if !isValidRole(req.Role) {
		return userServ.InputUser{}, ErrInvalidRole
	}
	if strings.TrimSpace(req.Name) == "" {
		return userServ.InputUser{}, ErrNameCannotBeBlank
	}
	if strings.TrimSpace(req.Email) == "" {
		return userServ.InputUser{}, ErrEmailCannotBeBlank
	}
	if strings.TrimSpace(req.Password) == "" {
		return userServ.InputUser{}, ErrPasswordCannotBeBlank
	}
	if strings.TrimSpace(req.Phone) == "" {
		return userServ.InputUser{}, ErrPhoneCannotBeBlank
	}
	if strings.TrimSpace(req.Role) == "" {
		return userServ.InputUser{}, ErrRoleCannotBeBlank
	}

	return userServ.InputUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
		IsActive: req.IsActive,
	}, nil
}

// CreateUser creates new user with given input
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decode
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleUserError(w, ErrInvalidBodyRequest)
		return
	}

	// 2. Validate request
	userInput, err := validateUserInput(req)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 3. Create User
	result, err := h.userServ.CreateUser(r.Context(), userInput)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 4. Return result
	utils.WriteJSONResponse(w, http.StatusCreated, result)
}

const (
	OrderTypeASC  string = "asc"
	OrderTypeDESC        = "desc"
)

type sortParams struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type getUserRequest struct {
	ID         int             `json:"id"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	IsActive   null.Bool       `json:"is_active"`
	Role       string          `json:"role"`
	Sort       sortParams      `json:"sort"`
	Pagination paginationInput `json:"pagination"`
}

type usersResponse struct {
	Users      []model.User `json:"users"`
	Pagination pagination   `json:"pagination"`
}

func validateGetUserInput(req getUserRequest) (userServ.InputGetUser, error) {
	if req.ID < 0 {
		return userServ.InputGetUser{}, ErrInvalidID
	}

	filterEmail := strings.TrimSpace(req.Email)
	if filterEmail != "" {
		if _, err := mail.ParseAddress(filterEmail); err != nil {
			return userServ.InputGetUser{}, ErrInvalidEmail
		}
	}
	filterRole := strings.TrimSpace(req.Role)
	if filterRole != "" && filterRole != UserRoleAdmin && filterRole != UserRoleGuest {
		return userServ.InputGetUser{}, ErrInvalidRole
	}
	filterName := strings.TrimSpace(req.Name)

	sortName := strings.TrimSpace(req.Sort.Name)
	if sortName != "" && sortName != OrderTypeASC && sortName != OrderTypeDESC {
		return userServ.InputGetUser{}, ErrInvalidSortType
	}
	sortEmail := strings.TrimSpace(req.Sort.Email)
	if sortEmail != "" && sortEmail != OrderTypeASC && sortEmail != OrderTypeDESC {
		return userServ.InputGetUser{}, ErrInvalidSortType
	}
	sortCreatedAt := strings.TrimSpace(req.Sort.CreatedAt)
	if sortCreatedAt != "" && sortCreatedAt != OrderTypeASC && sortCreatedAt != OrderTypeDESC {
		return userServ.InputGetUser{}, ErrInvalidSortType
	}

	pageArgs, err := validatePagination(req.Pagination)
	if err != nil {
		return userServ.InputGetUser{}, err
	}

	return userServ.InputGetUser{
		ID:       req.ID,
		Email:    filterEmail,
		Name:     filterName,
		IsActive: req.IsActive,
		Role:     filterRole,
		Sort: userServ.SortArgs{
			Name:      sortName,
			Email:     sortEmail,
			CreatedAt: sortCreatedAt,
		},
		Pagination: pageArgs,
	}, nil
}

// GetUsers returns list of users with given input.
func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Decode request
	var req getUserRequest
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleUserError(w, ErrInvalidBodyRequest)
			return
		}
	}

	// 2. Validate request
	getInput, err := validateGetUserInput(req)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 2. Get users
	result, totalCount, err := h.userServ.GetUsers(r.Context(), getInput)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 3. Return result
	utils.WriteJSONResponse(w, http.StatusOK, usersResponse{
		Users: result,
		Pagination: pagination{
			CurrentPage: getInput.Pagination.Page,
			Limit:       getInput.Pagination.Limit,
			TotalCount:  totalCount,
		},
	})
}

// UpdateUser new v1 handler, update one user by id
func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get id from url param
	id := chi.URLParam(r, "id")

	// 2. Get update fields from request body
	var userReq userRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		utils.WriteErrorResponse(w, ErrInvalidBodyRequest)
		return
	}

	// 3. validate ID
	// validate update fields
	// Create new input user from validated data
	userID, err := validateUserID(id)
	if err != nil {
		handleUserError(w, err)
		return
	}

	inputUser, err := validateUserInput(userReq)
	if err != nil {
		handleUserError(w, err)
		return
	}
	inputUser.ID = userID

	// 4. Call service to update user
	if err = h.userServ.UpdateUser(r.Context(), inputUser); err != nil {
		handleUserError(w, err)
		return
	}

	// 5. response updated User to client
	utils.WriteJSONResponse(w, http.StatusOK, utils.SuccessResponse{
		Success: true, Msg: "Update user successfully",
	})
}

func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get user ID from url param
	id := chi.URLParam(r, "id")

	// 2. Validate ID
	userID, err := validateUserID(id)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 3. Get user using "id"
	result, err := h.userServ.GetUser(r.Context(), userID)
	if err != nil {
		handleUserError(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, result)
}

const (
	MsgDeleteUserSuccess = "Delete user successfully"
)

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get user ID from url param
	id := chi.URLParam(r, "id")

	// 2. Validate ID
	userID, err := validateUserID(id)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// 3. Delete user using "id"
	if err := h.userServ.DeleteUser(r.Context(), userID); err != nil {
		handleUserError(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, utils.SuccessResponse{
		Success: true, Msg: MsgDeleteUserSuccess,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateLoginReq(loginReq LoginRequest) (userServ.LoginInput, error) {
	email := strings.TrimSpace(loginReq.Email)
	if email == "" {
		return userServ.LoginInput{}, ErrEmailCannotBeBlank
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return userServ.LoginInput{}, ErrInvalidEmail
	}
	if loginReq.Password == "" {
		return userServ.LoginInput{}, ErrPasswordCannotBeBlank
	}
	return userServ.LoginInput{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	}, nil
}

// Login handle login request
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Get request body
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		handleUserError(w, ErrInvalidBodyRequest)
		return
	}

	// Validate login request
	loginInput, err := validateLoginReq(loginReq)
	if err != nil {
		handleUserError(w, err)
		return
	}

	// Call login func of service
	result, err := h.userServ.Login(r.Context(), loginInput)
	if err != nil {
		handleUserError(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, result)
}
