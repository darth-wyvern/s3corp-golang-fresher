package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/handler/gql/graph/mod"
	productServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
)

func stringToPtr(s string) *string {
	return &s
}

func int64ToPtr(i int64) *int64 {
	return &i
}

func intToPtr(i int) *int {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

// CreateProduct is the resolver for the CreateProduct field.
func (m *mutationResolver) CreateProduct(ctx context.Context, input mod.CreateProductInput) (*mod.Product, error) {
	// 1. Validate product input
	newProduct, err := validProductInput(input)
	if err != nil {
		return nil, err
	}

	// 2. Call to service to create new product
	result, err := m.productServ.CreateProduct(ctx, newProduct)
	if err != nil {
		switch err {
		case productServ.ErrUserNotExist:
			return nil, errUserNotExist
		default:
			return nil, errInternalServerError
		}
	}

	return &mod.Product{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		Price:       result.Price,
		Quantity:    result.Quantity,
		IsActive:    result.IsActive,
		UserID:      result.UserID,
	}, nil

}

func validProductInput(input mod.CreateProductInput) (productServ.ProductInput, error) {
	var result productServ.ProductInput
	// Validate title of the product input
	result.Title = strings.TrimSpace(input.Title)
	if result.Title == "" {
		return productServ.ProductInput{}, errTitleCannotBeBlank
	}

	// Validate description of the product input
	if input.Description != nil {
		result.Description = strings.TrimSpace(*input.Description)
	}

	// Validate product price
	if input.Price <= 0 {
		return productServ.ProductInput{}, errInvalidPrice
	}
	result.Price = input.Price

	// Validate quantity
	if input.Quantity < 0 {
		return productServ.ProductInput{}, errInvalidQuantity
	}
	result.Quantity = input.Quantity

	// Validate user_id
	if input.UserID < 0 {
		return productServ.ProductInput{}, errInvalidUserID
	}
	result.UserID = input.UserID

	// Validate isActive
	if input.IsActive == mod.ActiveTypeYes {
		result.IsActive = true
	}

	return result, nil
}

const (
	DefaultLimit = 20
	MaxLimit     = 1000
)

func validateProductPagination(pagination mod.PaginationInput) (productServ.Pagination, error) {
	if *pagination.Page < 1 {
		return productServ.Pagination{}, errInvalidPage
	}
	if *pagination.Limit < 0 || *pagination.Limit > MaxLimit {
		return productServ.Pagination{}, errInvalidLimit
	}

	pageInput := productServ.Pagination{
		Page:  *pagination.Page,
		Limit: *pagination.Limit,
	}

	return pageInput, nil
}

const (
	orderTypeASC  = "asc"
	orderTypeDESC = "desc"
)

func validGetProductsInput(input mod.GetProductsInput) (productServ.GetProductsInput, error) {
	servInput := productServ.GetProductsInput{}

	// Validate ID filter field if any
	if input.ID != nil {
		if *input.ID < 0 {
			return productServ.GetProductsInput{}, errInvalidID
		} else {
			servInput.ID = *input.ID
		}
	}

	if input.Title != nil {
		servInput.Title = strings.TrimSpace(*input.Title)
	}

	// Validate price range filter field if any
	if input.PriceRange != nil {
		if (input.PriceRange.MinPrice != 0 || input.PriceRange.MaxPrice != 0) &&
			(input.PriceRange.MinPrice < 0 || input.PriceRange.MaxPrice < 0) {
			return productServ.GetProductsInput{}, errInvalidPriceRange
		} else {
			servInput.PriceRange = productServ.PriceRange{
				MinPrice: input.PriceRange.MinPrice,
				MaxPrice: input.PriceRange.MaxPrice,
			}
		}
	}

	if input.IsActive != nil {
		servInput.IsActive = null.NewBool(*input.IsActive, true)
	}

	if input.UserID != nil {
		if *input.UserID < 0 {
			return productServ.GetProductsInput{}, errInvalidUserID
		} else {
			servInput.UserID = *input.UserID
		}
	}

	// Validate order by if any
	if input.OrderBy != nil {
		if input.OrderBy.Title != nil {
			orderByTitle := strings.TrimSpace(*input.OrderBy.Title)
			if orderByTitle != "" && orderByTitle != orderTypeASC && orderByTitle != orderTypeDESC {
				return productServ.GetProductsInput{}, errInvalidOrderBy
			} else {
				servInput.OrderBy.Title = orderByTitle
			}
		}
		if input.OrderBy.CreatedAt != nil {
			orderByCreatedAt := strings.TrimSpace(*input.OrderBy.CreatedAt)
			if orderByCreatedAt != "" && orderByCreatedAt != orderTypeASC && orderByCreatedAt != orderTypeDESC {
				return productServ.GetProductsInput{}, errInvalidOrderBy
			} else {
				servInput.OrderBy.CreatedAt = orderByCreatedAt
			}
		}
		if input.OrderBy.Price != nil {
			orderByPrice := strings.TrimSpace(*input.OrderBy.Price)
			if orderByPrice != "" && orderByPrice != orderTypeASC && orderByPrice != orderTypeDESC {
				return productServ.GetProductsInput{}, errInvalidOrderBy
			} else {
				servInput.OrderBy.Price = orderByPrice
			}
		}
		if input.OrderBy.Quantity != nil {
			orderByQuantity := strings.TrimSpace(*input.OrderBy.Quantity)
			if orderByQuantity != "" && orderByQuantity != orderTypeASC && orderByQuantity != orderTypeDESC {
				return productServ.GetProductsInput{}, errInvalidOrderBy
			} else {
				servInput.OrderBy.Quantity = orderByQuantity
			}
		}
	}

	// Validate pagination
	if input.Pagination != nil {
		pageArgs, err := validateProductPagination(*input.Pagination)
		if err != nil {
			return productServ.GetProductsInput{}, err
		} else {
			servInput.Pagination = pageArgs
		}
	} else {
		servInput.Pagination = productServ.Pagination{
			Page:  1,
			Limit: DefaultLimit,
		}
	}

	return servInput, nil
}

func (q *queryResolver) GetProducts(ctx context.Context, input mod.GetProductsInput) (*mod.GetProductsOutput, error) {
	// Validate request body
	getProductsInput, err := validGetProductsInput(input)
	if err != nil {
		return &mod.GetProductsOutput{}, err
	}

	// Call service to get products
	products, totalCount, err := q.productServ.GetProducts(ctx, getProductsInput)
	if err != nil {
		return &mod.GetProductsOutput{}, err
	}

	productsOutput := make([]*mod.Product, len(products))
	for i, p := range products {
		productsOutput[i] = &mod.Product{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			IsActive:    p.IsActive,
			UserID:      p.User.ID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}
	return &mod.GetProductsOutput{
		Products: productsOutput,
		Pagination: &mod.Pagination{
			CurrentPage: &getProductsInput.Pagination.Page,
			Limit:       &getProductsInput.Pagination.Limit,
			TotalCount:  &totalCount,
		},
	}, nil
}
