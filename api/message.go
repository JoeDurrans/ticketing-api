package api

import (
	"net/http"
	"ticketing-api/auth"
	"ticketing-api/chat"
	"ticketing-api/types"

	"golang.org/x/net/websocket"
)

func (s *APIServer) handleChatGroup(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	ticket, err := s.db.Ticket.GetByID(id)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", r.Header.Get("Sec-WebSocket-Protocol"))

	err = auth.AccountIDAuth(r, ticket.AuthorID, types.RoleAdmin, types.RoleEditor)
	if err != nil {
		return err
	}

	group, ok := s.chatGroups.LoadOrStore(ticket.ID, chat.CreateGroup(ticket.ID, func() {
		s.chatGroups.Delete(ticket.ID)
	}))
	if !ok {
		go group.(*chat.Group).Start()
	}

	websocket.Server{
		Handler: websocket.Handler(func(conn *websocket.Conn) {
			client := chat.CreateClient(conn, group.(*chat.Group), s.db)
			client.Connect()
		}),
	}.ServeHTTP(w, r)

	return nil
}

type CreateMessageRequest struct {
	TicketID int    `json:"ticket_id"`
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}
