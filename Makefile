define setup_env
    $(eval ENV_FILE := $(1))
    $(eval include $(ENV_FILE))
    $(eval export)
endef

.PHONY: setup db db-migration docker-build-go-image docker-run-go-image run down test gql-gen

setup: db db-migration docker-build-go-image docker-run-go-image

docker-build-go-image:
	@docker build -f Dockerfile -t s3corp-golang-fresher-go-dev .

run:
	$(call setup_env,.env.dev)
	@go run cmd/serverd/main.go

vendor:
	@go mod tidy
	@go mod vendor

db:
	@docker-compose up -d db

db-migration:
	@docker-compose up -d db-migrate

docker-run-go-image:
	@docker-compose up -d app

down:
	@docker-compose down

test:
	$(call setup_env,.env.dev)
	@go test -v ./...

gql-gen:
	@cd internal/handler/gql && go run github.com/99designs/gqlgen
