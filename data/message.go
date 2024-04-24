package data

import (
	"fmt"
	"ticketing-api/types"
	"time"

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

func (m *MessageAdapter) Get(id int) ([]*types.Message, error) {
	scanner := m.db.Query("SELECT id, ticket_id, author_id, content, created_at, updated_at FROM message WHERE ticket_id = ? ORDER BY created_at DESC", id).Iter().Scanner()

	messages := []*types.Message{}

	for scanner.Next() {
		msg, err := scanIntoMessage(scanner)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (m *MessageAdapter) Create(message *types.Message) (*types.Message, error) {
	err := m.db.Query("INSERT INTO message (id, ticket_id, author_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", message.ID, message.TicketID, message.AuthorID, message.Content, message.CreatedAt, message.UpdatedAt).Exec()
	if err != nil {
		return nil, fmt.Errorf("error creating message")
	}

	return message, nil
}

func (m *MessageAdapter) Delete(id string, created_at time.Time, ticket_id int) error {
	err := m.db.Query("DELETE FROM message WHERE id = ? AND created_at = ? AND ticket_id = ?", id, created_at, ticket_id).Exec()
	if err != nil {
		return fmt.Errorf("error deleting message")
	}

	return nil
}

func (m *MessageAdapter) GetByID(id string, created_at time.Time, ticket_id int) (*types.Message, error) {
	scanner := m.db.Query("SELECT id, ticket_id, author_id, content, created_at, updated_at FROM message WHERE id = ? AND created_at = ? AND ticket_id = ?", id, created_at, ticket_id).Iter().Scanner()

	for scanner.Next() {
		return scanIntoMessage(scanner)
	}

	return nil, fmt.Errorf("mesage %s not found", id)
}

func (m *MessageAdapter) Update(message *types.Message) (*types.Message, error) {
	err := m.db.Query("UPDATE message SET author_id = ?, content = ?, updated_at = ? WHERE id = ? AND created_at = ? AND ticket_id = ?", message.AuthorID, message.Content, message.UpdatedAt, message.ID, message.CreatedAt, message.TicketID).Exec()
	if err != nil {
		return nil, fmt.Errorf("error updating message %w", err)
	}

	return message, nil
}

func scanIntoMessage(scanner gocql.Scanner) (*types.Message, error) {
	msg := &types.Message{}

	err := scanner.Scan(&msg.ID, &msg.TicketID, &msg.AuthorID, &msg.Content, &msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error reading message")
	}

	return msg, nil
}
