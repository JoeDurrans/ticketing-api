package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"ticketing-api/auth"
	"ticketing-api/types"
)

func (s *APIServer) handleCreateTicket(w http.ResponseWriter, r *http.Request) error {
	req := &CreateTicketRequest{}

	err := DecodeRequest(r, req)
	if err != nil {
		return err
	}

	if req.AuthorID > 0 {
		if ok := auth.AccountIDAuth(r, req.AuthorID, types.RoleAdmin); !ok {
			return fmt.Errorf("permission denied")
		}
	}

	if req.Status != "" {
		if ok := auth.AccountIDAuth(r, req.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
			return fmt.Errorf("permission denied")
		}
	}

	if len(req.AssigneeIDs) > 0 {
		if ok := auth.AccountIDAuth(r, req.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
			return fmt.Errorf("permission denied")
		}
	}

	ticket := CreateTicket(req.Title, req.Description, req.AuthorID, req.Status, req.AssigneeIDs)

	ticket, err = s.db.Ticket.Create(ticket)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "ticket created", Data: ticket})
}

func (s *APIServer) handleGetTickets(w http.ResponseWriter, r *http.Request) error {

	authorID, err := getAuthorID(r)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, authorID, types.RoleAdmin, types.RoleEditor); !ok {
		return fmt.Errorf("permission denied")
	}

	assigneeIDs, err := getAssigneeIDs(r)
	if err != nil {
		return err
	}

	tickets := []*types.Ticket{}

	if authorID != 0 && len(assigneeIDs) > 0 {
		tickets, err = s.db.Ticket.GetByAuthorIDAssigneeIDs(authorID, assigneeIDs)
		if err != nil {
			return err
		}
	} else if authorID != 0 {
		tickets, err = s.db.Ticket.GetByAuthorID(authorID)
		if err != nil {
			return err
		}
	} else if len(assigneeIDs) > 0 {
		tickets, err = s.db.Ticket.GetByAssigneeIDs(assigneeIDs)
		if err != nil {
			return err
		}
	} else {
		tickets, err = s.db.Ticket.Get()
		if err != nil {
			return err
		}
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "tickets found", Data: tickets})
}

func (s *APIServer) handleGetTicketByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	ticket, err := s.db.Ticket.GetByID(id)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, ticket.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
		return fmt.Errorf("permission denied")
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "ticket found", Data: ticket})
}

// needs to be modified so users cannot modify status, author, or assignee
func (s *APIServer) handleUpdateTicket(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	ticket, err := s.db.Ticket.GetByID(id)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, ticket.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
		return fmt.Errorf("permission denied")
	}

	req := &CreateTicketRequest{}

	err = DecodeRequest(r, req)
	if err != nil {
		return err
	}

	if req.Title != "" {
		ticket.Title = req.Title
	}

	if req.Description != "" {
		ticket.Description = req.Description
	}

	if req.Status != "" {
		if ok := auth.AccountIDAuth(r, req.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
			return fmt.Errorf("permission denied")
		}

		ticket.Status = req.Status
	}

	if len(req.AssigneeIDs) > 0 {
		if ok := auth.AccountIDAuth(r, req.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
			return fmt.Errorf("permission denied")
		}

		ticket.AssigneeIDs = req.AssigneeIDs
	}

	ticket, err = s.db.Ticket.Update(ticket)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "ticket updated", Data: ticket})
}

func (s *APIServer) handleDeleteTicket(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	ticket, err := s.db.Ticket.GetByID(id)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(r, ticket.AuthorID, types.RoleAdmin, types.RoleEditor); !ok {
		return fmt.Errorf("permission denied")
	}

	err = s.db.Ticket.Delete(id)
	if err != nil {
		return err
	}

	return EncodeResponse(w, http.StatusOK, &APIResponse{Status: http.StatusOK, Message: "ticket deleted"})
}

func getAuthorID(r *http.Request) (int, error) {
	authorIDStr := r.URL.Query().Get("author_id")

	if authorIDStr == "" {
		return 0, nil
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		return authorID, fmt.Errorf("error converting author_id: %w", err)
	}

	return authorID, nil
}

func getAssigneeIDs(r *http.Request) ([]int, error) {
	assigneeIDsStr := r.URL.Query().Get("assignee_ids")

	assigneeIDs := []int{}

	if assigneeIDsStr == "" {
		return assigneeIDs, nil
	}

	for _, idStr := range strings.Split(assigneeIDsStr, "+") {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return assigneeIDs, fmt.Errorf("error converting assignee_ids: %w", err)
		}

		assigneeIDs = append(assigneeIDs, id)
	}

	return assigneeIDs, nil
}

func CreateTicket(title string, description string, authorID int, status types.Status, assigneeIDs []int) *types.Ticket {
	return &types.Ticket{
		Title:       title,
		Description: description,
		AuthorID:    authorID,
		Status:      status,
		AssigneeIDs: assigneeIDs,
	}
}

type CreateTicketRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	AuthorID    int          `json:"author"`
	Status      types.Status `json:"status"`
	AssigneeIDs []int        `json:"assignee_ids"`
}
