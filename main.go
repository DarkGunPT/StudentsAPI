package main

import (
	"context"
	"example/user/restapi/config"
	"example/user/restapi/routes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Create context with timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the MongoDB connection
	client, err := config.ConnectDB(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Error connecting to mongodb:", err)
		return
	}

	// Initialize Gorilla Mux router
	router := mux.NewRouter()

	// Register students-related routes
	routes.StudentsRoutes(client, "SIPresentation", "Students", router)

	log.Println("Starting the http server at port :8080")
	// Start the HTTP server on port 8080
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting the http server:", err)
		return
	}
}
