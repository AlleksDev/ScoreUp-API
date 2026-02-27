package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleConnection(c *gin.Context) {
	role := c.Query("role")
	userID := c.Query("user_id")

	if role == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role y user_id son requeridos"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("error al actualizar a websocket: %v", err)
		return
	}

	client := NewClient(h.hub, conn, role, userID)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}