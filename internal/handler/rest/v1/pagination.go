package v1

import userService "github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"

const (
	DefaultLimit = 20
	MaxLimit     = 1000
)

type paginationInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type pagination struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalCount  int64 `json:"total_count"`
}

func validatePagination(pagination paginationInput) (userService.Pagination, error) {
	if pagination.Page < 0 {
		return userService.Pagination{}, ErrInvalidPaginationPage
	}
	if pagination.Limit < 0 || pagination.Limit > MaxLimit {
		return userService.Pagination{}, ErrInvalidPaginationLimit
	}

	pageInput := userService.Pagination{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}
	if pagination.Page == 0 {
		pageInput.Page = 1
	}
	if pagination.Limit == 0 {
		pageInput.Limit = DefaultLimit
	}

	return pageInput, nil
}
