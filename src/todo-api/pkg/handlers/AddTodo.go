package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/willvelida/aca-go-todo/pkg/mocks"
	"github.com/willvelida/aca-go-todo/pkg/models"
)

func AddTodo(w http.ResponseWriter, r *http.Request) {
	// read to request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var todo models.Todo
	json.Unmarshal(body, &todo)

	// Append to Todo mocks
	todo.ID = uuid.New().String()
	mocks.Todos = append(mocks.Todos, todo)

	// Send a 201 response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}
