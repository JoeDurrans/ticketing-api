package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"ticketing-api/auth"
	"ticketing-api/data"
	"ticketing-api/types"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/net/websocket"
)

func (c *Client) handleCreateMessage(data json.RawMessage) error {
	req := &CreateMessageRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	id, err := gocql.RandomUUID()
	if err != nil {
		return err
	}

	currentTime := time.Now()

	accountID, err := auth.GetID(c.Conn.Request())
	if err != nil {
		return err
	}

	message, err := CreateMessage(id.String(), c.Group.TicketID, accountID, req.Content, currentTime, currentTime)
	if err != nil {
		return err
	}

	message, err = c.db.Message.Create(message)
	if err != nil {
		return err
	}

	c.Group.Broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionCreate, Message: "message created", Data: message}

	return nil
}

func (c *Client) handleDeleteMessage(data json.RawMessage) error {
	req := &DeleteMessageRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	message, err := c.db.Message.GetByID(req.ID)
	if err != nil {
		return err
	}

	log.Println(message)

	if ok := auth.AccountIDAuth(c.Conn.Request(), message.AuthorID, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	err = c.db.Message.Delete(message.ID)
	if err != nil {
		return err
	}

	c.Group.Broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionDelete, Message: "message deleted", Data: &DeleteMessageResponse{ID: message.ID}}

	return nil
}

func (c *Client) handleUpdateMessage(data []byte) error {
	req := &UpdateMessageRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	message, err := c.db.Message.GetByID(req.ID)
	if err != nil {
		return err
	}

	if ok := auth.AccountIDAuth(c.Conn.Request(), message.AuthorID, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	message.Content = req.Content

	message, err = c.db.Message.Update(message)
	if err != nil {
		return err
	}

	c.Group.Broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionUpdate, Message: "message updated", Data: message}

	return nil
}

type Client struct {
	Conn  *websocket.Conn
	Group *Group
	db    *data.DataAdapter
	Send  chan *WSMessage
	once  *sync.Once
}

func CreateClient(conn *websocket.Conn, group *Group, db *data.DataAdapter) *Client {
	return &Client{
		Conn:  conn,
		Group: group,
		db:    db,
		Send:  make(chan *WSMessage),
		once:  &sync.Once{},
	}
}

func (c *Client) Connect() {
	c.Group.Register <- c

	go c.Write()
	c.Read()
}

func (c *Client) Disconnect() {
	c.once.Do(func() {
		c.Group.Unregister <- c
		c.Conn.Close()
		close(c.Send)
	})
}

func (c *Client) Read() {
	defer c.Disconnect()

	for {
		req := &MessageRequest{}
		err := websocket.JSON.Receive(c.Conn, &req)
		if err != nil {
			c.Send <- &WSMessage{Status: StatusError, Message: err.Error()}
			return
		}

		switch req.Action {
		case ActionCreate:
			err = c.handleCreateMessage(req.Data)
			if err != nil {
				c.Send <- &WSMessage{Status: StatusError, Action: ActionCreate, Message: err.Error()}
			}
		case ActionUpdate:
			err = c.handleUpdateMessage(req.Data)
			if err != nil {
				c.Send <- &WSMessage{Status: StatusError, Action: ActionUpdate, Message: err.Error()}
			}
		case ActionDelete:
			err = c.handleDeleteMessage(req.Data)
			if err != nil {
				c.Send <- &WSMessage{Status: StatusError, Action: ActionDelete, Message: err.Error()}
			}
		}
	}
}

func (c *Client) Write() {
	defer c.Disconnect()

	for message := range c.Send {
		err := websocket.JSON.Send(c.Conn, message)
		if err != nil {
			return
		}
	}
}

func CreateMessage(id string, ticketID int, authorID int, content string, createdAt time.Time, updatedAt time.Time) (*types.Message, error) {
	return &types.Message{
		ID:        id,
		TicketID:  ticketID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type MessageRequest struct {
	Action Action          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type UpdateMessageRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type DeleteMessageRequest struct {
	ID string `json:"id"`
}

type DeleteMessageResponse struct {
	ID string `json:"id"`
}
