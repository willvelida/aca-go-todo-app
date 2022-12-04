package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/handlers"
)

func TestGetAllTodos(t *testing.T) {
	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetAllTodos)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetTodo(t *testing.T) {
	testTodo := []struct {
		ID          string
		Title       string
		IsCompleted bool
	}{
		{
			ID:          "1",
			Title:       "Write presentation",
			IsCompleted: true,
		},
	}

	for _, tc := range testTodo {
		path := fmt.Sprintf("/todo/%s", tc.ID)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/todo/{id}", handlers.GetTodo).Methods("GET")
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}
}

func TestGetTodoNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "2")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(handlers.GetTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAddTodo(t *testing.T) {
	var requestBody = []byte(`{"title":"Write presentation", "isCompleted":true}`)
	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddTodo)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateTodo(t *testing.T) {
	var requestBody = []byte(`{"id":"1","title":"Write presentation", "isCompleted":true}`)
	req, err := http.NewRequest("PUT", "/todo", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UpdateTodo)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteTodo(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "1")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(handlers.DeleteTodo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthCheck)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
