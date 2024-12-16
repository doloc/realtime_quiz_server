package session

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func GenerateHostSessionID(hostID string) string {
	return uuid.NewString() // Session ID duy nhất và đảm bảo bảo mật
}

func GeneratePlayerSessionID(userID string) string {
	// Tạo một generator với nguồn mới dựa trên thời gian hiện tại
	randomGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Random salt
	salt := randomGen.Intn(1000000)

	// Dữ liệu đầu vào để hash
	data := userID + time.Now().String() + string(salt)

	// Hash dữ liệu
	hash := sha256.Sum256([]byte(data))

	// Trả về session ID dạng chuỗi hex
	return hex.EncodeToString(hash[:])
}
