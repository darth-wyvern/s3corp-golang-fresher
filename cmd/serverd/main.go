package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vinhnv1/s3corp-golang-fresher/cmd/serverd/router"
	"github.com/vinhnv1/s3corp-golang-fresher/pkg/db"
)

func main() {
	// Create DB connection
	dbConn, err := db.DBConnect(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("DB connection error ", err)
	}
	fmt.Println("DB connection success")

	// Create routes
	r := router.InitRouter(dbConn)

	fmt.Println("Server is running on port 5000")

	// Start server
	if err := http.ListenAndServe(":5000", r); err != nil {
		fmt.Println("Server error ", err)
	}
}
