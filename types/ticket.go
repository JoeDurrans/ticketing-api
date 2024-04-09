package types

import "time"

type Status string

const (
	StatusOpen     Status = "open"
	StatusPending  Status = "pending"
	StutusActive   Status = "active"
	StatusResolved Status = "resolved"
	StatusClosed   Status = "closed"
)

type Ticket struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	AuthorID    int       `json:"author_id"`
	AssigneeIDs []int     `json:"assignee_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
