package data

import (
	"ticketing-api/types"
)

type AccountSocket interface {
	Create(*types.Account) (*types.Account, error)
	Get() ([]*types.Account, error)
	GetByID(int) (*types.Account, error)
	GetByUsername(string) (*types.Account, error)
	Update(*types.Account) (*types.Account, error)
	Delete(int) error
}

type TicketSocket interface {
	Create(*types.Ticket) (*types.Ticket, error)
	Get() ([]*types.Ticket, error)
	GetByAssigneeIDs([]int) ([]*types.Ticket, error)
	GetByAuthorID(int) ([]*types.Ticket, error)
	GetByAuthorIDAssigneeIDs(int, []int) ([]*types.Ticket, error)
	GetByID(int) (*types.Ticket, error)
	Update(*types.Ticket) (*types.Ticket, error)
	Delete(int) error
}

type MessageSocket interface {
	Create(*types.Message) (*types.Message, error)
	Get(int) ([]*types.Message, error)
	GetByID(string) (*types.Message, error)
	Update(*types.Message) (*types.Message, error)
	Delete(string) error
}

type DataAdapter struct {
	Account AccountSocket
	Ticket  TicketSocket
	Message MessageSocket
}

func CreateDataAdapter(account AccountSocket, ticket TicketSocket, message MessageSocket) *DataAdapter {
	return &DataAdapter{
		Account: account,
		Ticket:  ticket,
		Message: message,
	}
}
