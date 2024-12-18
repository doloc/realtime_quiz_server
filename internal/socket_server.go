package internal

import (
	"fmt"
	"net/http"
	"realtime_quiz_server/cache"

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
func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Lỗi nâng cấp WebSocket:", err)
		return
	}

	// Lấy thông tin từ query parameters
	query := r.URL.Query()
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
}
