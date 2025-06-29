package main

import "time"

const (
	TODO        = "todo"
	IN_PROGRESS = "in progress"
	DONE        = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
