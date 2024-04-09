package chat

import (
	"log"
	"sync"
)

type Group struct {
	TicketID   int
	Clients    *sync.Map
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *WSMessage
	OnStop     func()
}

func CreateGroup(ticketID int, onStop func()) *Group {
	return &Group{
		TicketID:   ticketID,
		Clients:    &sync.Map{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *WSMessage),
		OnStop:     onStop,
	}
}

func (g *Group) Start() error {
	log.Println("Starting chat group for ticket", g.TicketID)
	for {
		select {
		case client := <-g.Register:
			g.registerClient(client)
		case client := <-g.Unregister:
			g.unregisterClient(client)
		case message := <-g.Broadcast:
			g.broadcastMessage(message)
		}
	}
}

func (g *Group) Stop() {
	close(g.Broadcast)
	close(g.Register)
	close(g.Unregister)
	if g.OnStop != nil {
		g.OnStop()
	}
}

func (g *Group) registerClient(client *Client) {
	g.Clients.Store(client, true)
}

func (g *Group) unregisterClient(client *Client) {
	g.Clients.Delete(client)
	if g.isEmpty() {
		g.Stop()
		log.Println("Chat group for ticket", g.TicketID, "is empty, stopping")
	}
}

func (g *Group) isEmpty() bool {
	empty := true
	g.Clients.Range(func(_, _ any) bool {
		empty = false
		return false
	})
	return empty
}

func (g *Group) broadcastMessage(message *WSMessage) {
	g.Clients.Range(func(client, _ any) bool {
		client.(*Client).Send <- message
		return true
	})
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type WSMessage struct {
	Status  Status `json:"status"`
	Action  Action `json:"action"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
