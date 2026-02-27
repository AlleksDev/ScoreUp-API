package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// En producción, reemplazar por una validación de origen concreta.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSHandler maneja la conexión WS entrante.
type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// HandleConnection upgradea una petición HTTP a WebSocket.
func (h *WSHandler) HandleConnection(c *gin.Context) {
	role := c.Query("role")
	userID := c.Query("user_id")
	channel := c.Query("channel")

	if role == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role y user_id son requeridos"})
		return
	}

	// Rechazar conexión si el hub está lleno.
	if h.hub.IsFull() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Límite de conexiones WebSocket alcanzado",
		})
		log.Printf("[WS] Conexión rechazada: límite de %d clientes alcanzado", MaxClients)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Error al actualizar a websocket: %v", err)
		return
	}

	client := NewClient(h.hub, conn, role, userID, channel)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
