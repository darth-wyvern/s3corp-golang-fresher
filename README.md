# s3corp-golang-fresher

Get started with Golang and Postgresql

## Introduction

## How to run?

Use `make setup` to run the application. This will start postgresql, migrate the database, build the app to docker image, and start the app.

Besides, we provide the following commands:

```Bash
make setup # to setup the application and start it

make db # to run postgresql

make db-migration # to run migrate for postgresql

make run # to run the application in local

make docker-build-go-image # to build the application to docker image

make docker-run-go-image # to run only the go application on docker container

make down # to stop all the application

make test # to run tests

make gql-gen # to generate the graphql models
```

## User APIs

Create user: POST /api/v1/users 

Request body:

```json
{
  "name": "Test",
  "email": "test@test.com",
  "password": "123456789",
  "phone": "123456789",
  "role": "GUEST",
  "is_active": true
}
```

Update user: POST /api/v1/users/{id}

Request body:

```json
{
  "name": "Test",
  "email": "test@test.com",
  "password": "123456789",
  "phone": "123456789",
  "role": "GUEST",
  "is_active": true
}
```

Get users: GET /api/v1/users

Request body:

```json
{
  "name": "test",
  "email": "test@example.com",
  "role": "GUEST",
  "is_active": true,
  "sort": {
    "name": "asc",
    "email": "asc",
    "created_at": "desc"
  },
  "pagination": {
    "page": 1,
    "limit": 10
  }
}
```

Get user: GET /api/v1/users/{id}

Request body: none

Delete user: DELETE /api/v1/users/{id}

Request body: none

Login: POST /api/v1/users/login

Request body: 
```json
  {
  "email":"mai@example.com",
  "password":"123456789"
}
```

## Product APIs

Update product: PUT /api/v1/products/{id}

Request body:

```json
{
  "title": "Test",
  "description": "Test",
  "price": 6.5,
  "quantity": 10,
  "is_active": true,
  "user_id": 1
}
```

Create product: POST /api/v1/products

Request body:

```json
{
  "title": "Test",
  "description": "Test",
  "price": 6.5,
  "quantity": 10,
  "is_active": true,
  "user_id": 1
}
```

Get product: GET /api/v1/products/{id}

Request body: none

Get product: GET /api/v1/products

Request body:

```json
{
    "id":1, 
    "title":"", 
    "price_range":
    {
      "from":100, 
      "to":3000
    }, 
    "is_active":false, 
    "user_id":1,
    "order_by":{
      "title":"desc",
      "quantity":"desc",
      "price":"asc",
      "created_at":"asc"
    },
    "pagination":{
      "limit":2,
      "page":1
    }
}
```
Export product to csv file: GET /api/v1/products/export/csv/

Request body:

```json
{
    "id":1, 
    "title":"", 
    "price_range":
    {
      "from":100, 
      "to":3000
    }, 
    "is_active":false, 
    "user_id":1,
    "order_by":{
      "title":"desc",
      "quantity":"desc",
      "price":"asc",
      "created_at":"asc"
    }
}
```
Create product using GraphQL: POST /api/v1/graphql

Request body:

```graphql
mutation {
    CreateProduct(
        input:{title:"test",description:"test",price: 100000, isActive:NO, quantity:10,userID:1}
    ){
        id,title,description,price,quantity,isActive, userID
    }
}
```

Get products using GraphQL: GET /api/v1/graphql

Request body:

```graphql
query {
    GetProducts(
        input:{title: "test", pagination: {page: 1, limit: 500}}
    ){
        pagination {
            currentPage
            limit
            totalCount
        }
        products {
            id
            title
        }
    }
}
```

Import product csv: POST /api/v1/products/import-csv

Request body: form-data key `file` with value is csv file

Create order: POST /api/v1/orders

Request body:
```json
{
    "note": "New order",
    "user_id": 6,
    "items": [
        {
            "product_id": 1008,
            "quantity": 10,
            "discount": 0,
            "note":"item 1"
        },
        {
            "product_id": 1009,
            "quantity": 20,
            "discount": 0.6,
            "note": "item 2"
        },
        {
            "product_id": 1010,
            "quantity": 20,
            "discount": 0.3,
            "note": "item 6"
        }
    ]
}
```

Get orders : GET /api/v1/orders

Request body:
```json
{
  "sort_by":{
    "created_at":"asc"
  },
  "filter":{
        "order_number":"ORDER_NUMBER_2",
        "status":"NEW"
	},
	"pagination":{
		"page":2,
		"limit":2
	}
}
```

Summary statistic: GET /api/v1/statistics

Response:

- success

```json
{
"users": {
  "total": 100,
  "total_inactive": 10
},
"products": {
  "total": 100,
  "total_inactive": 10
},
"orders": {
  "total_new": 100,
  "total_pending": 5,
  "total_success": 10,
  "total_failed": 6
},
"latestOrders": [
  {
    "order_id": 1,
    "order_number": "101",
    "order_date": "2020-01-01",
    "status": "NEW",
    "user_id": "2",
    "total": 100000
  },
  {
    "order_id": 2,
    "order_number": "102",
    "order_date": "2020-01-01",
    "status": "NEW",
    "user_id": "2",
    "total": 100000
  }
]
}
```