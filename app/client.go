package app

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	openBracket   = []byte{'['}
	closeBracket   = []byte{']'}
	comma   = []byte{','}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {

	// The hub
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Deferred requests to process onDisconnect
	requests map[string]Request
}

// read sends messages from the websocket connection to the hub.
//
// The application runs read in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()

		for k, v := range c.requests {
			log.Printf("🐳 %s = %v\n", k, v)
			go processRequest(c, &v)
			delete(c.requests, k)
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		// Parse the request and send it to Mongo
		request := Request{}
		err := c.conn.ReadJSON(&request)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		log.Printf("⭐️: %v\n", request)

		if request.OnDisconnect {
			// Defer the request to process on disconnect
			c.requests[request.Uid] = request
		} else {
			// Immediately process the requests
			processRequest(c, &request)
		}
	}
}

// write sends messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()

		for k, v := range c.requests {
			log.Printf("🐳: %s = %v\n", k, v)
			delete(c.requests, k)
		}

	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			//w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Print("c.conn.NextWriter", err)
				return
			}
			n := len(c.send)

			if n > 0 {
				w.Write(openBracket)
			}

			w.Write(message)

			// Add queued messages to the array
			for i := 0; i < n; i++ {
				w.Write(comma)
				w.Write(<-c.send)
			}

			if n > 0 {
				w.Write(closeBracket)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeResponse(data map[string]interface{}) {
	message, _ := json.Marshal(data)
	buffer := new(bytes.Buffer)
	error := json.Compact(buffer, message)
	if error != nil {
		log.Print("💩 Error compacting JSON: ", error)
		return
	}
	c.hub.broadcast <- buffer.Bytes()
}
