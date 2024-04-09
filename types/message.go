package types

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	TicketID  int       `json:"ticket_id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
