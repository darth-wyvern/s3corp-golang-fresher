package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
)

func (r impl) GetProduct(ctx context.Context, id int) (model.Product, error) {
	product, err := model.Products(qm.Where("id=?", id)).One(context.Background(), r.db)
	if product == nil {
		return model.Product{}, err
	}
	return *product, err
}

func (r impl) ExistsProductByID(ctx context.Context, id int) (bool, error) {
	return model.Products(model.ProductWhere.ID.EQ(id)).Exists(context.Background(), r.db)
}

func (r impl) CreateProduct(ctx context.Context, newProduct model.Product) (model.Product, error) {
	err := newProduct.Insert(context.Background(), r.db, boil.Whitelist("title", "description", "price", "quantity", "is_active", "user_id", "created_at", "updated_at"))
	if err != nil {
		return model.Product{}, err
	}
	return newProduct, nil
}

func (r impl) UpdateProduct(ctx context.Context, product model.Product) (int64, error) {
	affected, err := product.Update(context.Background(), r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r impl) DeleteProduct(ctx context.Context, id int) (int64, error) {
	affected, err := model.Products(qm.Where("id=?", id)).DeleteAll(context.Background(), r.db)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

type OrderBy struct {
	Title     string
	Price     string
	Quantity  string
	CreatedAt string
}

type PriceRange struct {
	MinPrice float64
	MaxPrice float64
}

type Pagination struct {
	Page  int
	Limit int
}

type Filter struct {
	ID         int
	Title      string
	PriceRange PriceRange
	IsActive   null.Bool
	UserID     int
	OrderBy    OrderBy
	Pagination Pagination
}

type ProductItem struct {
	ID          int
	Title       string
	Description string
	Price       float64
	Quantity    int
	IsActive    bool
	User        CreatedBy
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type CreatedBy struct {
	ID    int
	Name  string
	Email string
	Phone string
}

func (r impl) GetProducts(ctx context.Context, filter Filter) ([]ProductItem, int64, error) {
	// 1. Init query mods slice.
	var qms []qm.QueryMod

	// 2. Add filter condition.
	if filter.ID > 0 {
		qms = append(qms, model.ProductWhere.ID.EQ(filter.ID))
	}
	if filter.Title != "" {
		qms = append(qms, qm.Where("title LIKE ?", "%"+filter.Title+"%"))
	}
	if filter.UserID > 0 {
		qms = append(qms, model.ProductWhere.UserID.EQ(filter.UserID))
	}
	if filter.PriceRange.MinPrice > 0 || filter.PriceRange.MaxPrice > 0 {
		qms = append(qms, model.ProductWhere.Price.GTE(filter.PriceRange.MinPrice))
		qms = append(qms, model.ProductWhere.Price.LTE(filter.PriceRange.MaxPrice))
	}
	if filter.IsActive.Valid {
		qms = append(qms, model.ProductWhere.IsActive.EQ(filter.IsActive.Bool))
	}

	// 3. Calculate total rows of filtered products list.
	totalCount, err := model.Products(qms...).Count(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	//4. Sorting
	if filter.OrderBy != (OrderBy{}) {
		if filter.OrderBy.Title != "" {
			qms = append(qms, qm.OrderBy(model.ProductColumns.Title+" "+filter.OrderBy.Title))
		}
		if filter.OrderBy.CreatedAt != "" {
			qms = append(qms, qm.OrderBy(model.ProductColumns.CreatedAt+" "+filter.OrderBy.CreatedAt))
		}
		if filter.OrderBy.Price != "" {
			qms = append(qms, qm.OrderBy(model.ProductColumns.Price+" "+filter.OrderBy.Price))
		}
		if filter.OrderBy.Quantity != "" {
			qms = append(qms, qm.OrderBy(model.ProductColumns.Quantity+" "+filter.OrderBy.Quantity))
		}
	} else {
		qms = append(qms, qm.OrderBy(model.ProductColumns.UpdatedAt+" desc"))
	}

	// 5. Load relationships
	qms = append(qms, qm.Load(model.ProductRels.User))

	// 6. Add pagination condition.
	if filter.Pagination != (Pagination{}) {
		qms = append(
			qms,
			qm.Offset(filter.Pagination.Limit*(filter.Pagination.Page-1)), // Example: Pagination is 2, offset is 10. So offset is 10.
			qm.Limit(filter.Pagination.Limit))
	}

	// 7. Get the products with the queries
	productSlice, err := model.Products(qms...).All(context.Background(), r.db)
	if err != nil {
		return []ProductItem{}, 0, err
	}

	// 8. Map the productSlice to []productItem
	var result = make([]ProductItem, len(productSlice))
	for i, p := range productSlice {
		result[i] = ProductItem{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			IsActive:    p.IsActive,
			User: CreatedBy{
				ID:    p.R.User.ID,
				Name:  p.R.User.Name,
				Email: p.R.User.Email,
				Phone: p.R.User.Phone,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return result, totalCount, nil
}

func (r impl) InsertAll(ctx context.Context, tx *sql.Tx, products []model.Product) error {
	// Init query string
	queryStr := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES ",
		model.TableNames.Products,
		model.ProductColumns.Title, model.ProductColumns.Description, model.ProductColumns.Price,
		model.ProductColumns.Quantity, model.ProductColumns.IsActive, model.ProductColumns.UserID,
	)

	var values []interface{}
	totalField := 6
	idx := 0

	// Loop through the products and append the query string and values
	for _, p := range products {
		queryStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6)
		values = append(values, p.Title, p.Description, p.Price, p.Quantity, p.IsActive, p.UserID)
		idx += totalField
	}

	// Trim the last ',' from the query string
	queryStr = strings.TrimSuffix(queryStr, ",")

	// Execute the query
	if _, err := queries.Raw(queryStr, values...).Exec(tx); err != nil {
		return err
	}

	return nil
}

type SummaryStatistics struct {
	Total         int64
	TotalInactive int64
}

func (r impl) GetStatistics(ctx context.Context) (SummaryStatistics, error) {
	// Get the total number of products
	total, err := model.Products().Count(ctx, r.db)
	if err != nil {
		return SummaryStatistics{}, err
	}

	// Get the total number of inactive products
	totalInactive, err := model.Products(model.ProductWhere.IsActive.EQ(false)).Count(ctx, r.db)
	if err != nil {
		return SummaryStatistics{}, err
	}
	return SummaryStatistics{
		Total:         total,
		TotalInactive: totalInactive,
	}, nil
}
