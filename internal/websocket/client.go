package websocket

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	role    string
	userID  string
	channel string
}

func NewClient(hub *Hub, conn *websocket.Conn, role, userID, channel string) *Client {
	return &Client{
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		role:    role,
		userID:  userID,
		channel: channel,
	}
}

// ReadPump lee mensajes del cliente WS.
// Se ejecuta en su propia goroutine por cada conexión.
func (c *Client) ReadPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC][WS ReadPump] user=%s channel=%s: %v\n%s", c.userID, c.channel, r, debug.Stack())
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure) {
				// Solo loguear cierres inesperados, no reconexiones normales.
			}
			break
		}
		c.hub.broadcast <- message
	}
}

// WritePump envía mensajes al cliente WS con write-batching.
// Se ejecuta en su propia goroutine por cada conexión.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC][WS WritePump] user=%s channel=%s: %v\n%s", c.userID, c.channel, r, debug.Stack())
		}
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub cerró el canal.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Enviar el primer mensaje como frame individual.
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

			// Drenar mensajes pendientes — cada uno en su propio frame WS
			// para que el cliente pueda parsear cada JSON por separado.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if err := c.conn.WriteMessage(websocket.TextMessage, <-c.send); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
