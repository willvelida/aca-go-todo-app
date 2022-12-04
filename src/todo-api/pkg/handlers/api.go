package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/aca-go-todo/pkg/dal"
	"github.com/willvelida/aca-go-todo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = dal.ConnectDB()

func AddTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo models.Todo
	todo.ID = primitive.NewObjectID()

	_ = json.NewDecoder(r.Body).Decode(&todo)

	result, err := collection.InsertOne(context.TODO(), todo)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	var todo models.Todo

	filter := bson.M{"_id": id}
	err := collection.FindOne(r.Context(), filter).Decode(&todo)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	// Send a 200 response
	w.Header().Add("Content-Type", "application/json")

	var todos []models.Todo

	cur, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var todo models.Todo
		err := cur.Decode(&todo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var todo models.Todo

	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&todo)

	updatedTodo := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: todo.Title},
			{Key: "isCompleted", Value: todo.IsCompleted},
		},
		},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, updatedTodo).Decode(&todo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo.ID = id

	json.NewEncoder(w).Encode(todo)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}
