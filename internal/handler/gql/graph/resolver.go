package graph

import (
	"github.com/vinhnv1/s3corp-golang-fresher/internal/handler/gql/graph/generated"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userServ    user.IService
	productServ product.IProduct
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

// NewResolver creates a new Resolver
func NewResolver(userServ user.IService, productServ product.IProduct) Resolver {
	return Resolver{userServ: userServ, productServ: productServ}
}
