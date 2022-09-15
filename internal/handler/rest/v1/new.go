package v1

import (
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
)

type Handler struct {
	userServ    user.IService
	productServ product.IProduct
	orderServ   order.IOrder
}

func NewHandler(userServ user.IService, productServ product.IProduct, orderServ order.IOrder) Handler {
	return Handler{userServ: userServ, productServ: productServ, orderServ: orderServ}
}
