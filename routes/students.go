package routes

import (
	"example/user/restapi/api"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func StudentsRoutes(client *mongo.Client, dbName, userCollection string, router *mux.Router) {
	// Define route for creating a new student
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		api.CreateStudent(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to get all students
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		api.GetStudents(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to delete a service by id
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteStudent(client, dbName, userCollection, w, r)
	}).Methods("DELETE")

	// Define route to update a student by id
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateStudent(client, dbName, userCollection, w, r)
	}).Methods("PUT")
}
