package product

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/vinhnv1/s3corp-golang-fresher/internal/model"
	productRepo "github.com/vinhnv1/s3corp-golang-fresher/internal/repository/product"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/mail"
)

type ProductInput struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	IsActive    bool
	UserID      int
}

//GetProduct get product by id
func (serv impl) GetProduct(ctx context.Context, id int) (model.Product, error) {
	//Get data from repository
	result, err := serv.repo.Product().GetProduct(ctx, id)
	if err != nil {
		return model.Product{}, err
	}
	return result, nil

}

// CreateProduct create new product from product input
func (serv impl) CreateProduct(ctx context.Context, newProduct ProductInput) (model.Product, error) {
	// 1. Check exists user by product user_id
	existed, err := serv.repo.User().ExistsUserByID(ctx, newProduct.UserID)
	if err != nil {
		return model.Product{}, err
	}
	if !existed {
		return model.Product{}, ErrUserNotExist
	}

	result, err := serv.repo.Product().CreateProduct(ctx,
		model.Product{
			Title:       newProduct.Title,
			Description: newProduct.Description,
			Price:       newProduct.Price,
			Quantity:    newProduct.Quantity,
			IsActive:    newProduct.IsActive,
			UserID:      newProduct.UserID,
		})
	if err != nil {
		return model.Product{}, err
	}
	return result, nil
}

//UpdateProduct updates a product with the specified product
func (serv impl) UpdateProduct(ctx context.Context, id int, product ProductInput) error {
	// 1. Check exists user by product user_id
	existed, err := serv.repo.Product().ExistsProductByID(ctx, id)
	if err != nil {
		return err
	}

	if !existed {
		return ErrUserNotExist
	}

	// 2. Call repo func to update product, get result and error
	// result: number of rows are updated and any error
	result, err := serv.repo.Product().UpdateProduct(ctx, model.Product{
		ID:          id,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		IsActive:    product.IsActive,
		UserID:      product.UserID,
	})
	if err != nil {
		return err
	}

	// 2. If number of rows are updated equal zero, return not found error
	if result == 0 {
		return ErrProductNotFound
	}

	return nil
}

//DeleteProduct delete product by id
func (serv impl) DeleteProduct(ctx context.Context, id int) error {
	// 1. Call repo func to delete product, get result and error
	// affected rows: number of rows are deleted and any error
	affectedRows, err := serv.repo.Product().DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	// 2. If number of rows are deleted equal zero, return not found error
	if affectedRows == 0 {
		return ErrProductNotFound
	}
	// 3. Return nil if everything is successful
	return nil
}

type OrderInput struct {
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
	Page, Limit int
}

type GetProductsInput struct {
	ID         int
	Title      string
	PriceRange PriceRange
	IsActive   null.Bool
	UserID     int
	OrderBy    OrderInput
	Pagination Pagination
}

// toFilter returns a filter which be converted by getProductsInput
func toFilter(input GetProductsInput) productRepo.Filter {

	return productRepo.Filter{
		ID:    input.ID,
		Title: input.Title,
		PriceRange: productRepo.PriceRange{
			MinPrice: input.PriceRange.MinPrice,
			MaxPrice: input.PriceRange.MaxPrice,
		},
		IsActive: input.IsActive,
		UserID:   input.UserID,
		OrderBy: productRepo.OrderBy{
			Title:     input.OrderBy.Title,
			Price:     input.OrderBy.Price,
			Quantity:  input.OrderBy.Quantity,
			CreatedAt: input.OrderBy.CreatedAt,
		},
		Pagination: productRepo.Pagination{
			Limit: input.Pagination.Limit,
			Page:  input.Pagination.Page,
		},
	}
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

//GetProducts returns a list of products
func (serv impl) GetProducts(ctx context.Context, input GetProductsInput) ([]ProductItem, int64, error) {
	// Set pagination is default value if it is empty
	if input.Pagination == (Pagination{}) {
		input.Pagination.Page = 1
		input.Pagination.Limit = 20
	}
	// 1. Convert input to repository filter
	filter := toFilter(input)

	// 2. Get products using the filter
	products, totalCount, err := serv.repo.Product().GetProducts(ctx, filter)
	if err != nil {
		return []ProductItem{}, 0, err
	}
	var result = make([]ProductItem, len(products))
	for i, p := range products {
		result[i] = ProductItem{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			IsActive:    p.IsActive,
			User: CreatedBy{
				ID:    p.User.ID,
				Name:  p.User.Name,
				Email: p.User.Email,
				Phone: p.User.Phone,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return result, totalCount, err
}

var productColumns = []string{"id", "title", "description", "price", "quantity", "is_active", "user_id", "created_at", "updated_at"}

const (
	minProductCSVField = 6
)

func isValidProductFieldName(fieldName string) bool {
	for _, column := range productColumns {
		if strings.EqualFold(fieldName, column) {
			return true
		}
	}
	return false
}

func csvProductHeaders(fields []string) map[string]int {
	headers := make(map[string]int)
	for pos, field := range fields {
		headers[strings.ToLower(strings.TrimSpace(field))] = pos // use lowercase field name as key
	}
	return headers
}

func (serv impl) ImportProductCSV(ctx context.Context, fileName string, csvFile io.Reader) error {
	// Init csv reader
	reader := csv.NewReader(csvFile)

	// Read header row
	firstRow, err := reader.Read()
	if err != nil {
		log.Println("Error when read header row: ", err.Error())
		return err
	}
	for _, field := range firstRow {
		if !isValidProductFieldName(field) {
			return fmt.Errorf("invalid field name: %s", field)
		}
	}
	headers := csvProductHeaders(firstRow)

	var productList []model.Product
	rowIndex := 1 // skip header row

	// Loop all rows in csv file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error when read content csv at row: %d with error: %s\n", rowIndex, err.Error())
			return err
		}

		// Received row values
		rowMap := map[string]string{}
		for fieldName, pos := range headers {
			rowMap[fieldName] = strings.TrimSpace(record[pos])
		}

		if rowMap["title"] == "" {
			log.Printf("Skipping row (%d) because required field is empty (title: %v)\n", rowIndex, rowMap["title"])
			continue
		}

		price, err := strconv.ParseFloat(record[headers["price"]], 64)
		if price <= 0 {
			log.Printf("Skipping row (%d) because required field is invalid (price: %v)\n", rowIndex, price)
			continue
		}
		if err != nil {
			log.Printf("Skipping row (%d) because error when parse price: %s\n", rowIndex, err.Error())
			continue
		}

		quantity, err := strconv.Atoi(rowMap["quantity"])
		if quantity <= 0 {
			log.Printf("Skipping row (%d) because required field is invalid (quantity: %v)\n", rowIndex, quantity)
			continue
		}
		if err != nil {
			log.Printf("Skipping row (%d) because error when parse quantity: %s\n", rowIndex, err.Error())
			continue
		}

		isActive, err := strconv.ParseBool(rowMap["is_active"])
		if err != nil {
			log.Printf("Skipping row (%d) because error when parse is_active: %s\n", rowIndex, err.Error())
			continue
		}

		userId, err := strconv.Atoi(rowMap["user_id"])
		if err != nil {
			log.Printf("Skipping row (%d) because error when parse user_id: %s\n", rowIndex, err.Error())
			continue
		}

		productList = append(productList, model.Product{
			Title:       rowMap["title"],
			Description: rowMap["description"],
			Price:       price,
			Quantity:    quantity,
			IsActive:    isActive,
			UserID:      userId,
		})

		rowIndex++ // next row
	}

	// Insert all product data into database
	if err = serv.repo.Tx(context.Background(), func(tx *sql.Tx) error {
		return serv.repo.Product().InsertAll(ctx, tx, productList)
	}); err != nil {
		log.Printf("Failed processing product data from %s: %s\n", fileName, err.Error())
		return err
	}

	return nil
}

var exportProductColumns = []string{"ID", "Title", "Description", "Price", "Quantity", "Activated", "Created By", "Created Date", "Updated Date"}

// ExportProductsCSV returns name of the file which is exported
func (serv impl) ExportProductsCSV(ctx context.Context, input GetProductsInput, dir string) (string, error) {
	// TODO: handle data if number of products is greater than 1000
	input.Pagination.Limit = 1000
	input.Pagination.Page = 1

	// Call repository to get filled product list
	products, _, err := serv.repo.Product().GetProducts(ctx, toFilter(input))
	if err != nil {
		return "", err
	}
	// Define file bytes
	var fileBytes []byte

	// Define header rows, and add them to the fileBytes
	headerRow := []byte(strings.Join(exportProductColumns, ",") + "\n")
	fileBytes = append(fileBytes, headerRow...)

	// Add the rows to the fileBytes
	for _, p := range products {
		row := []byte(fmt.Sprintf("%d,%s,%s,%.2f,%d,%v,%s,%v,%v\n", p.ID, p.Title, p.Description, p.Price, p.Quantity, p.IsActive, p.User.Name, p.CreatedAt.Format("2006-01-02 15:04:05"), p.UpdatedAt.Format("2006-01-02 15:04:05")))
		fileBytes = append(fileBytes, row...)
	}

	// Create a file and check for errors
	now := time.Now()
	filename := fmt.Sprintf("products_%s.csv", now.Format("20060102"))
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", ErrFileCannotBeCreated
	}
	defer file.Close()
	// Write bytes to the file
	file.Write(fileBytes)

	// Send Email
	if err := mail.SendEmail(mail.EmailInput{
		To:         []string{os.Getenv("MAIL_TO")},
		Subject:    "Export product list to CSV file",
		Message:    "Export product list to CSV file successfully",
		Attachment: []string{filePath},
	}); err != nil {
		return "", fmt.Errorf("failed sending email: %v", err)
	}

	return filename, nil
}

func (serv impl) DownloadCSV(ctx context.Context, filePath string) ([]byte, error) {
	// Read the file
	result, err := os.ReadFile(filePath)
	if err != nil {
		return nil, ErrFileCannotBeRead
	}
	return result, nil
}

type SummaryStatistics struct {
	Total         int64
	TotalInactive int64
}
