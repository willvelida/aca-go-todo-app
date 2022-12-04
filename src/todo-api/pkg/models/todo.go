package models

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}
