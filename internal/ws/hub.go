package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.springy.io/api/document"
	"go.springy.io/internal/events"
	"log"
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
	hub  *Hub
	once sync.Once
)

func init() {
	log.Println("🌱 [Initializing Hub] 🌱")
	once.Do(func() {
		hub = &Hub{
			broadcast:  make(chan []byte),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[*Client]bool),
		}
	})
}

// Performs the ws upgrade
func Upgrade(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("💩", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), requests: make(map[string]document.DocumentRequest)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.write()
	go client.read()
}

func Run() {

	// Subscribe to websocket events
	subscriber := make(chan events.Event)
	events.Subscribe(events.Websocket, subscriber)

	for {
		select {
		case e := <-subscriber:
			if client, ok := e.Sender.(*Client); ok {
				if snapshot, ok := e.Data.(document.DocumentSnapshot); ok {
					go client.writeResponse(snapshot.Value)
				}
			}
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
