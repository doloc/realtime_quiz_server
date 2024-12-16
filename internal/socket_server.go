package internal

import (
	"fmt"
	"net/http"

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

	client := &Client{
		ID:   r.RemoteAddr, // Sử dụng địa chỉ IP làm ID tạm thời
		Conn: conn,
		Send: make(chan []byte),
	}

	// Đăng ký client vào Hub
	hub.Register <- client

	// Xử lý tin nhắn từ client
	go client.HandleMessages(hub)

	// Gửi tin nhắn đến client
	go client.WriteMessages()

	fmt.Printf("Client %s đã kết nối\n", client.ID)
}
