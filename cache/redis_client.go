package cache

import (
	"context"
	"log"
	"strconv"
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
	jsonString := `{"session_id":"` + sessionID + `,"host_id":"` + hostID + `","role":"Host"` + `,"expires_at":` + strconv.Itoa(int(expiresAt)) + `}`
	return Set("session:host:"+sessionID, jsonString, time.Hour)
}

func GetHostSession(sessionID string) (string, error) {
	return Get("session:host:" + sessionID)
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

func StorePlayerSession(sessionID string, name string, expiresAt int64) error {
	jsonString := `{"session_id":"` + sessionID + `","role":"Player"` + `,"name":"` + name + `","expires_at":` + strconv.Itoa(int(expiresAt)) + `}`
	return Set("session:player:"+sessionID, jsonString, time.Hour)
}

func GetPlayerSession(sessionID string) (string, error) {
	return Get("session:player:" + sessionID)
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

func UpdatePlayerScore(quizId string, sessionID string, score int64) error {
	key := "leaderboard:" + quizId
	_, err := RedisClient.ZIncrBy(ctx, key, float64(score), sessionID).Result()
	RedisClient.Expire(ctx, key, time.Hour)
	return err
}

func GetLeaderboard(quizId string) ([]string, error) {
	key := "leaderboard:" + quizId
	return RedisClient.ZRevRange(ctx, key, 0, 9).Result()
}

func GetPlayerScore(quizId string, sessionID string) (int64, error) {
	key := "leaderboard:" + quizId
	score, err := RedisClient.ZScore(ctx, key, sessionID).Result()
	return int64(score), err
}

func MutilGetPlayerScore(quizId string, sessionIDs []string) (map[string]int64, error) {
	key := "leaderboard:" + quizId
	result, err := RedisClient.ZMScore(ctx, key, sessionIDs...).Result()
	if err != nil {
		return nil, err
	}
	scores := make(map[string]int64)
	for i, score := range result {
		scores[sessionIDs[i]] = int64(score)
	}
	return scores, nil
}

func GetPlayerRanking(quizId string, sessionID string) (int64, error) {
	key := "leaderboard:" + quizId
	return RedisClient.ZRevRank(ctx, key, sessionID).Result()
}

func StoreQuestionTime(quizId string, questionId string, startTime int64) error {
	key := "question:" + quizId + ":start_time"
	return Set(key, startTime, time.Hour)
}

func GetQuestionTime(quizId string, questionId string) (int64, error) {
	key := "question:" + quizId + ":start_time"
	return RedisClient.Get(ctx, key).Int64()
}

func StorePlayerAnswer(quizId string, sessionID string, questionId string, answerTime int64, answer string) error {
	key := "answer:" + quizId + ":" + sessionID
	_, err := RedisClient.HSet(ctx, key, questionId+":time", answerTime, questionId+":answer", answer).Result()
	RedisClient.Expire(ctx, key, time.Hour)
	return err
}

func GetPlayerAnswer(quizId string, sessionID string, questionId string) (int64, string, error) {
	key := "answer:" + quizId + ":" + sessionID
	time, err := RedisClient.HGet(ctx, key, questionId+":time").Int64()
	if err != nil {
		return 0, "", err
	}
	answer, err := RedisClient.HGet(ctx, key, questionId+":answer").Result()
	return time, answer, err
}

func UpdateCounterPlayerAnswer(quizId string, questionId string, answer string) error {
	key := "question:" + quizId + ":" + questionId + ":answer_count"
	_, err := RedisClient.HIncrBy(ctx, key, answer, 1).Result()
	RedisClient.Expire(ctx, key, time.Hour)
	return err
}

func GetCounterPlayerAnswer(quizId string, questionId string, answer string) (int64, error) {
	key := "question:" + quizId + ":" + questionId + ":answer_count"
	return RedisClient.HGet(ctx, key, answer).Int64()
}
