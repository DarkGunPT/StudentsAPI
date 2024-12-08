package api

import (
	"bytes"
	"context"
	"encoding/json"
	"example/user/restapi/config"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func setupRouter() (*http.ServeMux, *mongo.Client) {
	router := http.NewServeMux()
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := config.ConnectDB(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Error connecting to mongodb:", err)
	}

	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			DeleteStudent(client, "SIPresentation", "Students", w, r)
		} else if r.Method == http.MethodPost {
			CreateStudent(client, "SIPresentation", "Students", w, r)
		}
	})

	return router, client
}

func TestCreateStudent(t *testing.T) {
	router, _ := setupRouter()

	// TEST 1
	student1 := map[string]string{
		"name": "Nuno",
		"age":  "24",
	}

	createJSON, _ := json.Marshal(student1)

	req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, got %v", w.Code)
	}

	// TEST 2...

}

func TestDeleteStudent(t *testing.T) {
	router, client := setupRouter()
	collection := client.Database("SIPresentation").Collection("Students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Insert the document into the collection
	result, _ := collection.InsertOne(ctx, map[string]interface{}{
		"name": "TestDelete",
		"age":  16,
	})

	// TEST 1
	incorrectId := map[string]string{
		"id": "67342a29b808c1f1f6c4059x",
	}

	deleteJSON, _ := json.Marshal(incorrectId)
	req := httptest.NewRequest(http.MethodDelete, "/students", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, got %v", w.Code)
	}

	// TEST 2
	notFoundId := map[string]string{
		"id": "67342a29b808c1f1f6c4059e",
	}

	deleteJSON, _ = json.Marshal(notFoundId)
	req = httptest.NewRequest(http.MethodDelete, "/students", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %v", w.Code)
	}

	// TEST 3
	correctId := map[string]interface{}{
		"id": result.InsertedID,
	}

	deleteJSON, _ = json.Marshal(correctId)
	req = httptest.NewRequest(http.MethodDelete, "/students", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}
