package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	productServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

type productRequest struct {
	Title       string  `json:"title"`       // required
	Description string  `json:"description"` // default ""
	Price       float64 `json:"price"`       // required
	Quantity    int     `json:"quantity"`    // default 0
	IsActive    bool    `json:"is_active"`   // default true
	UserID      int     `json:"user_id"`     // required
}

func validateProductInput(r productRequest) (productServ.ProductInput, error) {
	title := strings.TrimSpace(r.Title)
	if title == "" {
		return productServ.ProductInput{}, ErrTitleCannotBeBlank
	}
	if r.Price <= 0 {
		return productServ.ProductInput{}, ErrInvalidPrice
	}
	if r.Quantity < 0 {
		return productServ.ProductInput{}, ErrInvalidQuantity
	}
	if r.UserID < 0 {
		return productServ.ProductInput{}, ErrInvalidUserID
	}
	return productServ.ProductInput{
		Title:       title,
		Description: r.Description,
		Price:       r.Price,
		Quantity:    r.Quantity,
		IsActive:    r.IsActive,
		UserID:      r.UserID,
	}, nil
}

func validateProductID(id string) (int, error) {
	result, err := strconv.Atoi(id)
	if err != nil || result < 0 {
		return 0, ErrInvalidID
	}
	return result, nil
}

// GetProduct response product
func (h Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get id from url param
	// Parse id to int
	id := chi.URLParam(r, "id")
	parsedID, err := validateProductID(id)
	if err != nil {
		utils.WriteErrorResponse(w, err)
		return
	}
	// Call SERVICE to get data
	product, err := h.productServ.GetProduct(r.Context(), parsedID)
	if err != nil {
		if err == sql.ErrNoRows { // If no row
			utils.WriteErrorResponse(w, ErrProductNotFound)
		} else { // anything else
			utils.WriteErrorResponse(w, ErrInternalServerError)
		}
		return
	}
	// response data to client
	utils.WriteJSONResponse(w, http.StatusOK, product)
}

func (h Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get Request body
	var product productRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.WriteErrorResponse(w, ErrInvalidBodyRequest)
		return
	}
	// Validate Input
	validatedProduct, err := validateProductInput(product)
	if err != nil {
		utils.WriteErrorResponse(w, err)
		return
	}
	// Create product
	res, err := h.productServ.CreateProduct(r.Context(), validatedProduct)
	if err != nil {
		handleProductError(w, err)
		return
	}
	utils.WriteJSONResponse(w, http.StatusCreated, res)
}

const (
	MsgProductUpdatedSuccess = "Product updated successfully"
)

func (h Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Get input id from url param
	// Parse id
	id := chi.URLParam(r, "id")
	productID, pErr := validateProductID(id)
	if pErr != nil {
		handleProductError(w, pErr)
		return
	}

	// 2. Get Request body
	var product productRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		handleProductError(w, ErrInvalidBodyRequest)
		return
	}

	// 3. Validate Input
	validatedProduct, vErr := validateProductInput(product)
	if vErr != nil {
		handleProductError(w, vErr)
		return
	}

	// 4. Call update product function from product service
	if err := h.productServ.UpdateProduct(r.Context(), productID, validatedProduct); err != nil {
		handleProductError(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, utils.SuccessResponse{
		Success: true,
		Msg:     MsgProductUpdatedSuccess,
	})
}

//DeleteProduct handle delete product request by ID
func (h Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Get input id from url param
	// Parse id
	id := chi.URLParam(r, "id")
	parsedID, err := validateProductID(id)
	if err != nil {
		utils.WriteErrorResponse(w, err)
		return
	}
	// 2. Call delete product function from product service
	if err := h.productServ.DeleteProduct(r.Context(), parsedID); err != nil {
		handleProductError(w, err)
		return
	}
	// 3. if it's successful, respone successful message and status ok
	message := "Delete product successfully"
	utils.WriteJSONResponse(w, http.StatusOK, message)
}

type productItemResponse struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Quantity    int               `json:"quantity"`
	IsActive    bool              `json:"is_active"`
	User        createdByResponse `json:"user"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
type createdByResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type getProductsResponse struct {
	Products   []productItemResponse `json:"products"`
	Pagination pagination            `json:"pagination"`
}

// GetProducts handle get products request
func (h Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// 1. Get request body
	var getProductsReq getProductsRequest
	if r.ContentLength != 0 { // If body != none
		if err := json.NewDecoder(r.Body).Decode(&getProductsReq); err != nil {
			handleProductError(w, ErrInvalidBodyRequest)
			return
		}
	}
	// 2. Validate request body
	getProductsInput, err := validGetProductsInput(getProductsReq)
	if err != nil {
		handleProductError(w, err)
		return
	}

	// 3. Call service to get products
	products, totalCount, err := h.productServ.GetProducts(r.Context(), getProductsInput)
	if err != nil {
		handleProductError(w, err)
		return
	}
	result := make([]productItemResponse, len(products))
	for i, p := range products {
		result[i] = productItemResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			IsActive:    p.IsActive,
			User: createdByResponse{
				ID:    p.User.ID,
				Name:  p.User.Name,
				Email: p.User.Email,
				Phone: p.User.Phone,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	if getProductsInput.Pagination.Page == 0 && getProductsInput.Pagination.Limit == 0 {
		getProductsInput.Pagination.Page = 1
		getProductsInput.Pagination.Limit = 20
	}

	utils.WriteJSONResponse(w, http.StatusOK, getProductsResponse{
		Products: result,
		Pagination: pagination{
			CurrentPage: getProductsInput.Pagination.Page,
			Limit:       getProductsInput.Pagination.Limit,
			TotalCount:  totalCount,
		},
	})
}

const (
	maxSizeUploadCSVFile = 1024 * 1024 // 1 MB
)

// ImportProductCSV imports product from csv file
func (h Handler) ImportProductCSV(w http.ResponseWriter, r *http.Request) {
	// Get csv file from request body
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		handleProductError(w, ErrInvalidBodyRequest)
		return
	}
	defer file.Close()

	// Check max size upload csv file (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, maxSizeUploadCSVFile)
	if err := r.ParseMultipartForm(maxSizeUploadCSVFile); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, ErrFileSizeTooLarge)
		return
	}

	// Validate csv file extension
	if filepath.Ext(fileHeader.Filename) != ".csv" {
		handleProductError(w, ErrInvalidFileType)
		return
	}

	// Processing upload csv file
	go func() {
		if err := h.productServ.ImportProductCSV(r.Context(), fileHeader.Filename, file); err != nil {
			log.Println("Error when import CSV file: ", err.Error())
			return
		}
		log.Println("Import CSV file successfully")
	}()

	utils.WriteJSONResponse(w, http.StatusOK, fmt.Sprintf("Starting import product from csv: %s", fileHeader.Filename))
}

type productCSVResponse struct {
	ProductCSVURL string `json:"product_csv_url"`
}

func (h Handler) ExportProductsCSV(w http.ResponseWriter, r *http.Request) {
	// 1. define get products request
	var getProductsReq getProductsRequest
	if r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&getProductsReq); err != nil {
			handleProductError(w, ErrInvalidBodyRequest)
			return
		}
	}

	// 2. validate the request data
	getProductsInput, err := validGetProductsInput(getProductsReq)
	if err != nil {
		handleProductError(w, err)
		return
	}

	// 3. Call export func and send mail
	result, err := h.productServ.ExportProductsCSV(r.Context(), getProductsInput, "docs/csvfiles")
	if err != nil {
		handleProductError(w, err)
		return
	}
	// Define url to download FILE
	appURL := os.Getenv("APP_URL")
	urlToDownload := appURL + "/api/v1/files/" + result
	utils.WriteJSONResponse(w, http.StatusOK, productCSVResponse{ProductCSVURL: urlToDownload})
}

func validateFileName(fileName string) (string, error) {
	slices := strings.Split(fileName, ".")
	// Check file is csv file
	if slices[len(slices)-1] != "csv" {
		return "", ErrInvalidFileName
	}
	return fileName, nil
}

func (h Handler) DownloadCSVFile(w http.ResponseWriter, r *http.Request) {
	// Get filename from url param
	fileNameReq := chi.URLParam(r, "filename")

	//Validate filename
	fileName, err := validateFileName(fileNameReq)
	if err != nil {
		handleProductError(w, err)
		return
	}

	result, err := h.productServ.DownloadCSV(r.Context(), "docs/csvfiles/"+fileName)
	if err != nil {
		handleProductError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=file.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
