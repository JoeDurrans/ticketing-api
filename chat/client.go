package chat

import (
	"encoding/json"
	"fmt"
	"log"
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

	accountID, err := auth.GetID(c.conn.Request())
	if err != nil {
		return err
	}

	message, err := types.CreateMessage(id.String(), c.group.ticketID, accountID, req.Content, currentTime, currentTime)
	if err != nil {
		return err
	}

	message, err = c.db.Message.Create(message)
	if err != nil {
		return err
	}

	c.group.broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionCreate, Message: "message created", Data: message}

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

	if ok := auth.AccountIDAuth(c.conn.Request(), message.AuthorID, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	err = c.db.Message.Delete(message.ID)
	if err != nil {
		return err
	}

	c.group.broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionDelete, Message: "message deleted", Data: &DeleteMessageResponse{ID: message.ID}}

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

	if ok := auth.AccountIDAuth(c.conn.Request(), message.AuthorID, types.RoleAdmin); !ok {
		return fmt.Errorf("permission denied")
	}

	message.Content = req.Content

	message, err = c.db.Message.Update(message)
	if err != nil {
		return err
	}

	c.group.broadcast <- &WSMessage{Status: StatusSuccess, Action: ActionUpdate, Message: "message updated", Data: message}

	return nil
}

type Client struct {
	conn  *websocket.Conn
	group *Group
	db    *data.DataAdapter
	send  chan *WSMessage
}

func CreateClient(conn *websocket.Conn, group *Group, db *data.DataAdapter) *Client {
	return &Client{
		conn:  conn,
		group: group,
		db:    db,
		send:  make(chan *WSMessage),
	}
}

func (c *Client) Connect() {
	c.group.register <- c

	go c.Write()
	c.Read()
}

func (c *Client) Disconnect() {
	c.group.unregister <- c
	c.conn.Close()
	close(c.send)
}

func (c *Client) Read() {
	defer c.Disconnect()

	for {
		req := &MessageRequest{}
		err := websocket.JSON.Receive(c.conn, &req)
		if err != nil {
			c.send <- &WSMessage{Status: StatusError, Message: err.Error()}
			return
		}

		switch req.Action {
		case ActionCreate:
			err = c.handleCreateMessage(req.Data)
			if err != nil {
				c.send <- &WSMessage{Status: StatusError, Action: ActionCreate, Message: err.Error()}
			}
		case ActionUpdate:
			err = c.handleUpdateMessage(req.Data)
			if err != nil {
				c.send <- &WSMessage{Status: StatusError, Action: ActionUpdate, Message: err.Error()}
			}
		case ActionDelete:
			err = c.handleDeleteMessage(req.Data)
			if err != nil {
				c.send <- &WSMessage{Status: StatusError, Action: ActionDelete, Message: err.Error()}
			}
		}
	}
}

func (c *Client) Write() {
	for message := range c.send {
		err := websocket.JSON.Send(c.conn, message)
		if err != nil {
			log.Println("error sending message:", err)
			return
		}
	}
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
