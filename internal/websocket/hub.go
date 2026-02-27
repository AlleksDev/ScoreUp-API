package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
)

const (
	// MaxClients es el número máximo de conexiones WS simultáneas.
	MaxClients = 2048
)

// Hub mantiene el conjunto de clientes activos y difunde mensajes.
type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	stop       chan struct{} // señal de shutdown graceful
	count      atomic.Int64  // contador atómico para acceso rápido
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		stop:       make(chan struct{}),
	}
}

// Run es el event-loop principal. Debe ejecutarse en su propia goroutine.
func (h *Hub) Run() {
	for {
		select {
		case <-h.stop:
			h.mu.Lock()
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			h.mu.Unlock()
			log.Println("[WS] Hub detenido correctamente")
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.count.Add(1)
			h.mu.Unlock()
			log.Printf("[WS] Cliente conectado (channel=%s, role=%s) — total: %d",
				client.channel, client.role, h.count.Load())

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.count.Add(-1)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// El cliente no puede seguir el ritmo; se desconecta.
					close(client.send)
					delete(h.clients, client)
					h.count.Add(-1)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Shutdown detiene el Hub de forma ordenada.
func (h *Hub) Shutdown() {
	close(h.stop)
}

// ClientCount devuelve el número de clientes conectados (acceso lock-free).
func (h *Hub) ClientCount() int64 {
	return h.count.Load()
}

// IsFull indica si se alcanzó el límite de conexiones.
func (h *Hub) IsFull() bool {
	return h.count.Load() >= MaxClients
}

// BroadcastToRole envía un mensaje a todos los clientes con cierto rol.
func (h *Hub) BroadcastToRole(role string, message []byte) {
	var slow []*Client
	h.mu.RLock()
	for client := range h.clients {
		if client.role == role {
			select {
			case client.send <- message:
			default:
				slow = append(slow, client)
			}
		}
	}
	h.mu.RUnlock()
	// Desregistrar clientes lentos fuera del RLock para evitar race condition.
	for _, c := range slow {
		h.unregister <- c
	}
}

// BroadcastToAll envía un mensaje a todos los clientes vía el canal broadcast.
func (h *Hub) BroadcastToAll(message []byte) {
	h.broadcast <- message
}

// BroadcastToChannel envía data a todos los clientes suscritos a un canal específico.
func (h *Hub) BroadcastToChannel(channel string, message []byte) {
	var slow []*Client
	h.mu.RLock()
	for client := range h.clients {
		if client.channel == channel {
			select {
			case client.send <- message:
			default:
				slow = append(slow, client)
			}
		}
	}
	h.mu.RUnlock()
	// Desregistrar clientes lentos fuera del RLock.
	for _, c := range slow {
		h.unregister <- c
	}
}

// BroadcastJSON serializa data a JSON y lo envía al canal indicado.
func (h *Hub) BroadcastJSON(channel string, data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("[WS] Error serializando broadcast para '%s': %v", channel, err)
		return
	}
	h.BroadcastToChannel(channel, msg)
}
