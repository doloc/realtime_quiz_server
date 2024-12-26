package internal

import (
	"fmt"
	"net/http"
	"realtime_quiz_server/cache"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader dùng để nâng cấp HTTP thành WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Chấp nhận tất cả kết nối
	},
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

// ServeWebSocket xử lý một kết nối WebSocket từ client
func ServeWebSocket(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	defer conn.Close()
	fmt.Println("WebSocket connected")

	// Thiết lập timeout
	conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Timeout đọc
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Gia hạn timeout sau mỗi lần nhận Pong
		return nil
	})

	// Lấy thông tin từ query parameters
	query := c.Request.URL.Query()
	role := query.Get("role")
	quizID := query.Get("quizId")
	sessionID := query.Get("sessionId")

	if role == "host" {
		if _, err := cache.GetHostSession(sessionID); err != nil {
			fmt.Println("Không tìm thấy session host:", sessionID)
			return
		}
	} else {
		if _, err := cache.GetPlayerSession(sessionID); err != nil {
			fmt.Println("Không tìm thấy session player:", sessionID)
			return
		}
	}

	client := &Client{
		ID:     sessionID,
		Role:   role,
		QuizID: quizID,
		Conn:   conn,
		Send:   make(chan []byte),
	}

	// Đăng ký client vào Hub
	hub.Register <- client

	// Xử lý tin nhắn từ client
	go client.HandleMessages(hub)

	// Gửi tin nhắn đến client
	go client.WriteMessages()

	// Gửi Ping định kỳ tới client
	// go client.SendPing()
}

// SendPing gửi một thông điệp ping đến client
func (c *Client) SendPing() {
	// check if client is still connected
	if c.Conn == nil {
		return
	}
	ticker := time.NewTicker(30 * time.Second) // Gửi ping mỗi 30 giây
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := c.Conn.WriteMessage(websocket.PingMessage, nil) // Gửi Ping
			if err != nil {
				fmt.Printf("Lỗi gửi ping tới client %s: %v\n", c.ID, err)
				c.Conn.Close() // Đóng kết nối nếu có lỗi
				return
			}
		}
	}
}
