package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/willvelida/aca-go-todo/pkg/mocks"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	// Send a 200 response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mocks.Todos)
}
