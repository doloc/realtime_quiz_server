package internal

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Client đại diện cho một kết nối WebSocket
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

// Hub quản lý tất cả các kết nối WebSocket
type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// NewHub khởi tạo một Hub mới
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run chạy vòng lặp chính của Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("Client %s đã kết nối\n", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				fmt.Printf("Client %s đã ngắt kết nối\n", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// HandleMessages xử lý tin nhắn từ client
func (c *Client) HandleMessages(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("Lỗi đọc tin nhắn từ client %s: %v\n", c.ID, err)
			break
		}

		fmt.Printf("Nhận tin nhắn từ client %s: %s\n", c.ID, message)
		hub.Broadcast <- message // Gửi tin nhắn tới tất cả client
	}
}

// WriteMessages gửi tin nhắn tới client
func (c *Client) WriteMessages() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("Lỗi gửi tin nhắn tới client %s: %v\n", c.ID, err)
			break
		}
	}
}
