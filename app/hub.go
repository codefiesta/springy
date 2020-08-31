package app

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var (
	shared *Hub
	once   sync.Once
)

// Our singleton instance of the hub
func NewHub() *Hub {

	once.Do(func() {
		shared = &Hub{
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
		}
	})
	return shared
}

// Performs the ws upgrade
func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("💩", err)
		return
	}

	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256), requests: make(map[string]Request)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.write()
	go client.read()
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
