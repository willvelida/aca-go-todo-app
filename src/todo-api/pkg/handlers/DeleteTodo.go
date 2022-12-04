package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/mocks"
)

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// read the id parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Iterate over all the mock todos
	for i, todo := range mocks.Todos {
		if todo.ID == id {
			// Delete the todo
			mocks.Todos = append(mocks.Todos[:i], mocks.Todos[i+1:]...)

			// Send a 200 response
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Deleted")
			break
		}
	}
}
