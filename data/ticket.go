package data

import (
	"database/sql"
	"fmt"
	"ticketing-api/types"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type TicketAdapter struct {
	db *sql.DB
}

func CreateTicketAdapter(db *sql.DB) *TicketAdapter {
	return &TicketAdapter{
		db: db,
	}
}

func (t *TicketAdapter) Create(ticket *types.Ticket) (*types.Ticket, error) {
	id := 0
	err := t.db.QueryRow("INSERT INTO ticket (title, description, author_id, status) VALUES ($1, $2, $3, $4) RETURNING id", ticket.Title, ticket.Description, ticket.AuthorID, ticket.Status).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket")
	}

	ticket.ID = id

	for _, id := range ticket.AssigneeIDs {
		_, err := t.db.Exec("INSERT INTO assignee (ticket_id, account_id) VALUES ($1, $2)", ticket.ID, id)
		if err != nil {
			return nil, fmt.Errorf("error creating assignee")
		}
	}

	return ticket, nil
}

func (t *TicketAdapter) Get() ([]*types.Ticket, error) {
	return t.fetchTickets("SELECT ticket.id, ticket.title, ticket.description, ticket.status, ticket.author_id, ticket.created_at, ticket.updated_at, assignee.account_id FROM ticket LEFT JOIN assignee ON ticket.id = assignee.ticket_id")
}

func (t *TicketAdapter) GetByAuthorID(authorID int) ([]*types.Ticket, error) {
	return t.fetchTickets("SELECT ticket.id, ticket.title, ticket.description, ticket.status, ticket.author_id, ticket.created_at, ticket.updated_at, assignee.account_id FROM ticket LEFT JOIN assignee ON ticket.id = assignee.ticket_id WHERE author_id = $1", authorID)
}

func (t *TicketAdapter) GetByAssigneeIDs(assigneeIDs []int) ([]*types.Ticket, error) {
	return t.fetchTickets("SELECT ticket.id, ticket.title, ticket.description, ticket.status, ticket.author_id, ticket.created_at, ticket.updated_at, assignee.account_id FROM ticket LEFT JOIN assignee ON ticket.id = assignee.ticket_id WHERE assignee.account_id = $1", pq.Array(assigneeIDs))
}

func (t *TicketAdapter) GetByAuthorIDAssigneeIDs(authorID int, assigneeIDs []int) ([]*types.Ticket, error) {
	return t.fetchTickets("SELECT ticket.id, ticket.title, ticket.description, ticket.status, ticket.author_id, ticket.created_at, ticket.updated_at, assignee.account_id FROM ticket LEFT JOIN assignee ON ticket.id = assignee.ticket_id WHERE author_id = $1 AND assignee.account_id = ANY($2)", authorID, pq.Array(assigneeIDs))
}

func (t *TicketAdapter) GetByID(id int) (*types.Ticket, error) {
	tickets, err := t.fetchTickets("SELECT ticket.id, ticket.title, ticket.description, ticket.status, ticket.author_id, ticket.created_at, ticket.updated_at, assignee.account_id FROM ticket LEFT JOIN assignee ON ticket.id = assignee.ticket_id WHERE ticket.id = $1", id)
	if err != nil {
		return nil, err
	}

	if len(tickets) > 0 {
		return tickets[0], nil
	}

	return nil, fmt.Errorf("ticket %d not found", id)
}

func (t *TicketAdapter) Update(ticket *types.Ticket) (*types.Ticket, error) {
	_, err := t.db.Exec("UPDATE ticket SET title = $1, description = $2, status = $3 WHERE id = $4", ticket.Title, ticket.Description, ticket.Status, ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating ticket")
	}

	_, err = t.db.Exec("DELETE FROM assignee WHERE ticket_id = $1", ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("error deleting assignee")
	}

	for _, assigneeID := range ticket.AssigneeIDs {
		_, err := t.db.Exec("INSERT INTO assignee (ticket_id, account_id) VALUES ($1, $2)", ticket.ID, assigneeID)
		if err != nil {
			return nil, fmt.Errorf("error creating assignee")
		}

	}

	return ticket, nil
}

func (t *TicketAdapter) Delete(id int) error {
	_, err := t.db.Exec("DELETE FROM ticket WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting ticket")
	}

	return nil
}

func (t *TicketAdapter) fetchTickets(query string, args ...any) ([]*types.Ticket, error) {
	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickets")
	}

	ticketMap := make(map[int]*types.Ticket)

	for rows.Next() {
		ticket, err := scanIntoTicket(rows)
		if err != nil {
			return nil, err
		}

		if existingTicket, exists := ticketMap[ticket.ID]; exists {
			existingTicket.AssigneeIDs = append(existingTicket.AssigneeIDs, ticket.AssigneeIDs...)
		} else {
			ticketMap[ticket.ID] = ticket
		}
	}

	tickets := []*types.Ticket{}

	for _, ticket := range ticketMap {
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func scanIntoTicket(rows *sql.Rows) (*types.Ticket, error) {
	assigneeID := sql.NullInt64{}
	ticket := &types.Ticket{
		AssigneeIDs: []int{},
	}

	err := rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.AuthorID, &ticket.CreatedAt, &ticket.UpdatedAt, &assigneeID)
	if err != nil {
		return nil, fmt.Errorf("error reading ticket")
	}

	if assigneeID.Valid {
		ticket.AssigneeIDs = append(ticket.AssigneeIDs, int(assigneeID.Int64))
	}

	return ticket, nil
}
