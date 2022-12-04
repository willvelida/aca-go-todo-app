package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/mocks"
	"github.com/willvelida/aca-go-todo/pkg/models"
)

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// read the id parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Read the request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedTodo models.Todo
	json.Unmarshal(body, &updatedTodo)

	// Iterate over all the mock todos
	for i, todo := range mocks.Todos {
		if todo.ID == id {
			// Update the todo
			mocks.Todos[i] = updatedTodo

			// Send a 200 response
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Updated")
			break
		}
	}
}
