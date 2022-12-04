package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/todo", handlers.GetAllTodos).Methods(http.MethodGet)
	router.HandleFunc("/todo/{id}", handlers.GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/todo/{id}", handlers.DeleteTodo).Methods(http.MethodDelete)
	router.HandleFunc("/todo", handlers.AddTodo).Methods(http.MethodPost)
	router.HandleFunc("/todo/{id}", handlers.UpdateTodo).Methods(http.MethodPut)
	router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)

	log.Println("API is running")
	http.ListenAndServe(":8080", router)
}
