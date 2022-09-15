package v1

import (
	"strings"

	"github.com/volatiletech/null/v8"

	productServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
)

type orderRequest struct {
	Title     string `json:"title"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	CreatedAt string `json:"created_at"`
}

type priceRangeRequest struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type getProductsRequest struct {
	ID         int               `json:"id"`
	Title      string            `json:"title"`
	PriceRange priceRangeRequest `json:"price_range"`
	IsActive   null.Bool         `json:"is_active"`
	UserID     int               `json:"user_id"`
	OrderBy    orderRequest      `json:"order_by"`
	Pagination paginationInput   `json:"pagination"`
}

// validateGetProductsInput validate get products request
func validGetProductsInput(req getProductsRequest) (productServ.GetProductsInput, error) {
	// 1. Validate ID filter field if any
	if req.ID != 0 && req.ID < 0 {
		return productServ.GetProductsInput{}, ErrInvalidID
	}

	// 2. Validate price range filter field if any
	if (req.PriceRange.MinPrice != 0 || req.PriceRange.MaxPrice != 0) &&
		(req.PriceRange.MinPrice < 0 || req.PriceRange.MaxPrice < 0) {
		return productServ.GetProductsInput{}, ErrInvalidPriceRange
	}

	// 3. Validate user ID if any
	if req.UserID != 0 && req.UserID < 0 {
		return productServ.GetProductsInput{}, ErrInvalidUserID
	}

	// 4. Validate order by if any
	orderByTitle := strings.TrimSpace(req.OrderBy.Title)
	if orderByTitle != "" && orderByTitle != OrderTypeASC && orderByTitle != OrderTypeDESC {
		return productServ.GetProductsInput{}, ErrInvalidOrderBy
	}
	orderByCreatedAt := strings.TrimSpace(req.OrderBy.CreatedAt)
	if orderByCreatedAt != "" && orderByCreatedAt != OrderTypeASC && orderByCreatedAt != OrderTypeDESC {
		return productServ.GetProductsInput{}, ErrInvalidOrderBy
	}
	orderByPrice := strings.TrimSpace(req.OrderBy.Price)
	if orderByPrice != "" && orderByPrice != OrderTypeASC && orderByPrice != OrderTypeDESC {
		return productServ.GetProductsInput{}, ErrInvalidOrderBy
	}
	orderByQuantity := strings.TrimSpace(req.OrderBy.Quantity)
	if orderByQuantity != "" && orderByQuantity != OrderTypeASC && orderByQuantity != OrderTypeDESC {
		return productServ.GetProductsInput{}, ErrInvalidOrderBy
	}

	return productServ.GetProductsInput{
		ID:       req.ID,
		Title:    strings.TrimSpace(req.Title),
		IsActive: req.IsActive,
		PriceRange: productServ.PriceRange{
			MinPrice: req.PriceRange.MinPrice,
			MaxPrice: req.PriceRange.MaxPrice,
		},
		UserID: req.UserID,
		OrderBy: productServ.OrderInput{
			CreatedAt: orderByCreatedAt,
			Title:     orderByTitle,
			Price:     orderByPrice,
			Quantity:  orderByQuantity,
		},
		Pagination: productServ.Pagination{
			Limit: req.Pagination.Limit,
			Page:  req.Pagination.Page,
		},
	}, nil
}
