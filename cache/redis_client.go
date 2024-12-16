package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var RedisClient *redis.Client

// InitializeRedis khởi tạo Redis client
func InitializeRedis(redisAddr string, redisPassword string, redisDB int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Địa chỉ Redis
		Password: redisPassword,
		DB:       redisDB, // Sử dụng DB mặc định
	})

	// Kiểm tra kết nối Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối đến Redis: %v", err)
	} else {
		log.Println("Kết nối Redis thành công")
	}
}

// Set đặt giá trị vào Redis với thời gian hết hạn
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get lấy giá trị từ Redis
func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func StoreHostSession(sessionID string, hostID string, expiresAt int64) error {
	jsonString := `{"session_id":"` + sessionID + `,"host_id":"` + hostID + `","role":"Host"` + `","expires_at":` + string(expiresAt) + `"}`
	return Set("session:host:"+sessionID, jsonString, time.Hour)
}

func DeleteHostSession(sessionID string) error {
	return Delete("session:host:" + sessionID)
}

func UpdateQuizIdToHostSession(sessionID string, quizId string) error {
	sessionString, err := Get("session:host:" + sessionID)
	if err != nil {
		return err
	}

	ttl, err := RedisClient.TTL(ctx, "session:host:"+sessionID).Result()
	if err != nil {
		return err
	}

	err = Set("session:host:"+sessionID, sessionString[:len(sessionString)-1]+`,"current_quiz_id":"`+quizId+`"}`, 0)
	if err != nil {
		return err
	}

	return RedisClient.Expire(ctx, "session:host:"+sessionID, ttl).Err()
}

func UpdateStatusToHostSession(sessionID string, status string) error {
	sessionString, err := Get("session:host:" + sessionID)
	if err != nil {
		return err
	}

	ttl, err := RedisClient.TTL(ctx, "session:host:"+sessionID).Result()
	if err != nil {
		return err
	}

	err = Set("session:host:"+sessionID, sessionString[:len(sessionString)-1]+`,"status":"`+status+`"}`, 0)
	if err != nil {
		return err
	}

	return RedisClient.Expire(ctx, "session:host:"+sessionID, ttl).Err()
}

func StorePlayerSession(sessionID string, expiresAt int64) error {
	jsonString := `{"session_id":"` + sessionID + `","role":"Player"` + `","expires_at":` + string(expiresAt) + `"}`
	return Set("session:player:"+sessionID, jsonString, time.Hour)
}

func DeletePlayerSession(sessionID string) error {
	return Delete("session:player:" + sessionID)
}

func UpdateQuizIdToPlayerSession(sessionID string, quizId string) error {
	sessionString, err := Get("session:player:" + sessionID)
	if err != nil {
		return err
	}

	ttl, err := RedisClient.TTL(ctx, "session:player:"+sessionID).Result()
	if err != nil {
		return err
	}

	err = Set("session:player:"+sessionID, sessionString[:len(sessionString)-1]+`,"current_quiz_id":"`+quizId+`"}`, 0)
	if err != nil {
		return err
	}

	return RedisClient.Expire(ctx, "session:player:"+sessionID, ttl).Err()
}

func UpdateScoreToPlayerSession(sessionID string, score int) error {
	sessionString, err := Get("session:player:" + sessionID)
	if err != nil {
		return err
	}

	ttl, err := RedisClient.TTL(ctx, "session:player:"+sessionID).Result()
	if err != nil {
		return err
	}

	err = Set("session:player:"+sessionID, sessionString[:len(sessionString)-1]+`,"score":`+string(score)+`}`, 0)
	if err != nil {
		return err
	}

	return RedisClient.Expire(ctx, "session:player:"+sessionID, ttl).Err()
}
