package graph

import (
	"net/http"

	"github.com/vinhnv1/s3corp-golang-fresher/pkg/utils"
)

var (
	errInvalidPrice        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_price", Desc: "price is invalid"}
	errTitleCannotBeBlank  = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_input", Desc: "title cannot be blank"}
	errInvalidQuantity     = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_quantity", Desc: "quantity is invalid"}
	errInvalidUserID       = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_user_id", Desc: "use id is invalid"}
	errUserNotExist        = utils.ErrorResponse{Status: http.StatusNotFound, Code: "user_not_exist", Desc: "user does not exist"}
	errInvalidPage         = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_page", Desc: "page is invalid"}
	errInvalidLimit        = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_limit", Desc: "limit is invalid"}
	errInvalidID           = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_id", Desc: "id is invalid"}
	errInvalidPriceRange   = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_price_range", Desc: "price range is invalid"}
	errInvalidOrderBy      = utils.ErrorResponse{Status: http.StatusBadRequest, Code: "invalid_order_by", Desc: "order by is invalid"}
	errInternalServerError = utils.ErrorResponse{Status: http.StatusInternalServerError, Code: "internal_server_error", Desc: "internal server error"}
)
