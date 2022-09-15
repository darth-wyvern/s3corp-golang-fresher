package v1

import (
	"net/http"

	productServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
	userServ "github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

var (
	ErrNameCannotBeBlank      = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "name cannot be blank"}
	ErrEmailCannotBeBlank     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "email cannot be blank"}
	ErrPhoneCannotBeBlank     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "phone cannot be blank"}
	ErrPasswordCannotBeBlank  = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "password cannot be blank"}
	ErrRoleCannotBeBlank      = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "role cannot be blank"}
	ErrInvalidEmail           = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_email", Desc: "email is invalid"}
	ErrInvalidRole            = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_role", Desc: "role is invalid"}
	ErrInvalidSortField       = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_sort_field", Desc: "sort field is invalid"}
	ErrInvalidSortType        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_sort_type", Desc: "sort type is invalid"}
	ErrUserIDExisted          = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "user_id_existed", Desc: "user id is already exists"}
	ErrEmailExisted           = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "email_existed", Desc: "email is already exists"}
	ErrInvalidBodyRequest     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_body_request", Desc: "body request is invalid"}
	ErrInvalidID              = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_id", Desc: "id is invalid"}
	ErrInvalidPrice           = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_price", Desc: "price is invalid"}
	ErrTitleCannotBeBlank     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "title cannot be blank"}
	ErrInvalidQuantity        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_quantity", Desc: "quantity is invalid"}
	ErrInvalidUserID          = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_user_id", Desc: "user id is invalid"}
	ErrInvalidPaginationPage  = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_page", Desc: "page is invalid"}
	ErrInvalidPaginationLimit = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_limit", Desc: "limit is invalid"}
	ErrInvalidOrderBy         = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_order_by", Desc: "order by is invalid"}
	ErrInvalidPriceRange      = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_price_range", Desc: "price range is invalid"}
	ErrFileSizeTooLarge       = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "file_size_too_large", Desc: "file size too large"}
	ErrInvalidFileType        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_file_type", Desc: "file type is invalid"}
	ErrItemsCannotBeBlank     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "items cannot be blank"}
	ErrInvalidProductID       = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_product_id", Desc: "product id is invalid"}
	ErrInvalidDiscount        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_discount", Desc: "discount is invalid"}
	ErrInvalidFileName        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_file_name", Desc: "file name is invalid"}
	ErrInvalidOrderID         = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_order_id", Desc: "order id is invalid"}
	ErrInvalidOrderStatus     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_order_status", Desc: "order status is invalid"}
	ErrPasswordIncorrect      = utils.ErrorResponse{Status: http.StatusUnauthorized, Code: "incorrect_password", Desc: "password is incorrect"}
	ErrEmailNotExist          = utils.ErrorResponse{Status: http.StatusUnauthorized, Code: "email_not_exist", Desc: "email does not exist"}
	ErrFileNotExist           = utils.ErrorResponse{Status: http.StatusNotFound, Code: "file_not_exist", Desc: "file does not exist"}
	ErrUserNotExist           = utils.ErrorResponse{Status: http.StatusNotFound, Code: "user_not_exist", Desc: "user does not exist"}
	ErrProductNotFound        = utils.ErrorResponse{Status: http.StatusNotFound, Code: "product_not_found", Desc: "product is not found"}
	ErrUserNotFound           = utils.ErrorResponse{Status: http.StatusNotFound, Code: "user_not_found", Desc: "user is not found"}
	ErrInternalServerError    = utils.ErrorResponse{Status: http.StatusInternalServerError, Code: "internal_error", Desc: "internal server error"}
	ErrFileCannotBeCreated    = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "file_cannot_be_created", Desc: "file cannot be created"}
)

// handleUserError handle error and write to response
func handleUserError(w http.ResponseWriter, err error) {
	var v, ok = err.(utils.ErrorResponse)
	if ok {
		utils.WriteJSONResponse(w, v.Status, v)
	} else {
		switch err {
		case userServ.ErrEmailNotExist:
			utils.WriteJSONResponse(w, ErrEmailNotExist.Status, ErrEmailNotExist)
		case userServ.ErrEmailExisted:
			utils.WriteJSONResponse(w, ErrEmailExisted.Status, ErrEmailExisted)
		case userServ.ErrUserIDExisted:
			utils.WriteJSONResponse(w, ErrUserIDExisted.Status, ErrUserIDExisted)
		case userServ.ErrUserNotFound:
			utils.WriteJSONResponse(w, ErrUserNotFound.Status, ErrUserNotFound)
		case userServ.ErrPasswordIncorrect:
			utils.WriteJSONResponse(w, ErrPasswordIncorrect.Status, ErrPasswordIncorrect)
		default:
			utils.WriteJSONResponse(w, ErrInternalServerError.Status, ErrInternalServerError)
		}
	}
}

// handleProductError handle error and write to response
func handleProductError(w http.ResponseWriter, err error) {
	var v, ok = err.(utils.ErrorResponse)
	if ok {
		utils.WriteJSONResponse(w, v.Status, v)
	} else {
		switch err {
		case productServ.ErrProductNotFound:
			utils.WriteJSONResponse(w, ErrProductNotFound.Status, ErrProductNotFound)
		case productServ.ErrUserNotExist:
			utils.WriteJSONResponse(w, ErrUserNotExist.Status, ErrUserNotExist)
		case productServ.ErrFileCannotBeCreated:
			utils.WriteJSONResponse(w, ErrFileCannotBeCreated.Status, ErrFileCannotBeCreated)
		case productServ.ErrFileCannotBeRead:
			utils.WriteJSONResponse(w, ErrFileNotExist.Status, ErrFileNotExist)
		default:
			utils.WriteJSONResponse(w, ErrInternalServerError.Status, ErrInternalServerError)
		}
	}
}

// handleStatisticsError handle error and write to response
func handleStatisticError(w http.ResponseWriter, err error) {
	var v, ok = err.(utils.ErrorResponse)
	if ok {
		utils.WriteJSONResponse(w, v.Status, v)
	} else {
		utils.WriteJSONResponse(w, ErrInternalServerError.Status, ErrInternalServerError)
	}
}
