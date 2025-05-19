package entities

import "time"

type Task struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Deadline    time.Time `json:"deadline"`
	StatusID    int       `db:"status_id" json:"status_id"`
}

type CreateTaskDto struct {
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}
