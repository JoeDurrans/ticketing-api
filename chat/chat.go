package chat

import (
	"sync"
)

type Group struct {
	ticketID   int
	clients    *sync.Map
	register   chan *Client
	unregister chan *Client
	broadcast  chan *WSMessage
	onStop     func()
	once       *sync.Once
}

func CreateGroup(ticketID int, onStop func()) *Group {
	return &Group{
		ticketID:   ticketID,
		clients:    &sync.Map{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *WSMessage),
		onStop:     onStop,
		once:       &sync.Once{},
	}
}

func (g *Group) Start() error {
	for {
		select {
		case client := <-g.register:
			g.registerClient(client)
		case client := <-g.unregister:
			g.unregisterClient(client)
		case message := <-g.broadcast:
			g.broadcastMessage(message)
		}
	}
}

func (g *Group) Stop() {
	g.once.Do(func() {
		close(g.broadcast)
		close(g.register)
		close(g.unregister)

		if g.onStop != nil {
			g.onStop()
		}
	})
}

func (g *Group) registerClient(client *Client) {
	g.clients.Store(client, true)
}

func (g *Group) unregisterClient(client *Client) {
	g.clients.Delete(client)
	if g.isEmpty() {
		g.Stop()
	}
}

func (g *Group) isEmpty() bool {
	empty := true
	g.clients.Range(func(_, _ any) bool {
		empty = false
		return false
	})
	return empty
}

func (g *Group) broadcastMessage(message *WSMessage) {
	g.clients.Range(func(client, _ any) bool {
		if c, ok := client.(*Client); ok && c != nil {
			c.send <- message
		}

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
