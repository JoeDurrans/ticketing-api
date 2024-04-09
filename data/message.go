package data

import (
	"fmt"
	"ticketing-api/types"

	"github.com/gocql/gocql"
)

type MessageAdapter struct {
	db *gocql.Session
}

func CreateMessageAdapter(db *gocql.Session) *MessageAdapter {
	return &MessageAdapter{
		db: db,
	}
}

func (m *MessageAdapter) Create(message *types.Message) (*types.Message, error) {
	err := m.db.Query("INSERT INTO message (id, ticket_id, author_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", message.ID, message.TicketID, message.AuthorID, message.Content, message.CreatedAt, message.UpdatedAt).Exec()
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	return message, nil
}

func (m *MessageAdapter) Delete(id string) error {
	err := m.db.Query("DELETE FROM message WHERE id = ?", id).Exec()
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}

	return nil
}

func (m *MessageAdapter) GetByID(id string) (*types.Message, error) {
	scanner := m.db.Query("SELECT * FROM message WHERE id = ?", id).Iter().Scanner()

	for scanner.Next() {
		return scanIntoMessage(scanner)
	}

	return nil, fmt.Errorf("mesage %s not found", id)
}

func (m *MessageAdapter) Update(message *types.Message) (*types.Message, error) {
	err := m.db.Query("UPDATE message SET content = ?, updated_at = ? WHERE id = ?", message.Content, message.UpdatedAt, message.ID).Exec()
	if err != nil {
		return nil, fmt.Errorf("error updating message: %w", err)
	}

	return nil, nil
}

func scanIntoMessage(scanner gocql.Scanner) (*types.Message, error) {

	msg := &types.Message{}

	err := scanner.Scan(&msg.ID, &msg.AuthorID, &msg.Content, &msg.CreatedAt, &msg.TicketID, &msg.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning message: %w", err)
	}

	return msg, nil
}
