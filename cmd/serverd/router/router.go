package router

import (
	"database/sql"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/handler/gql/graph"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/handler/gql/graph/generated"
	v1 "github.com/vinhnv1/s3corp-golang-fresher/internal/handler/rest/v1"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/repository"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/order"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/product"
	"github.com/vinhnv1/s3corp-golang-fresher/internal/service/user"
)

// InitRouter return all handler
func InitRouter(db *sql.DB) *chi.Mux {
	repo := repository.New(db)
	userServ := user.New(repo)
	productServ := product.New(repo)
	orderServ := order.New(repo)
	h := v1.NewHandler(userServ, productServ, orderServ)

	resolver := graph.NewResolver(userServ, productServ)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Handle("/api/v1/graphql", graphRouter(resolver))
	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/products", productRouter(h))
		api.Route("/users", userRouter(h))
		api.Route("/orders", orderRouter(h))
		api.Route("/files", fileRouter(h))
		api.Get("/statistics", h.GetStatistics)
	})
	return r
}

func productRouter(h v1.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/export/csv", h.ExportProductsCSV)
		r.Get("/{id}", h.GetProduct)
		r.Get("/", h.GetProducts)
		r.Post("/", h.CreateProduct)
		r.Post("/import-csv", h.ImportProductCSV)
		r.Put("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
	}
}

func userRouter(h v1.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/", h.CreateUser)
		r.Get("/", h.GetUsers)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
	}
}
func fileRouter(h v1.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/{filename}", h.DownloadCSVFile)
	}
}

func graphRouter(resolver graph.Resolver) *handler.Server {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))
}

func orderRouter(h v1.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.CreateOrder)
		r.Get("/", h.GetOrders)
	}
}
