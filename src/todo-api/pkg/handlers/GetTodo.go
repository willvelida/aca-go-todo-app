package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/mocks"
)

func GetTodo(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Iterate over all the mock todos
	for _, todo := range mocks.Todos {
		if todo.ID == id {
			// Send a 200 response
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(todo)
			break
		} else {
			// Send a 404 response
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Not Found")
		}
	}
}
