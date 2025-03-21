package websocket

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn   *websocket.Conn
	UserID int
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

var manager *Manager

func InitManager() {
	manager = NewManager()
	go manager.Run()
}

func GetManager() *Manager {
	return manager
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.Conn.Close()
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.RLock()
			for client := range m.clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Conn.Close()
					delete(m.clients, client)
				}
			}
			m.mutex.RUnlock()
		}
	}
}

func (m *Manager) BroadcastToAll(message []byte) {
	m.broadcast <- message
}

func WebsocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Get user ID from context (set by JWT middleware)
		userID := c.Locals("user_id").(int)

		client := &Client{
			Conn:   c,
			UserID: userID,
		}

		manager.register <- client
		defer func() {
			manager.unregister <- client
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			// Echo the message back (optional)
			manager.broadcast <- msg
		}
	})
}
