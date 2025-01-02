package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"realtime_quiz_server/cache"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/entity/reference"
	"realtime_quiz_server/service"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// represents a websocket connection
type Client struct {
	ID     string
	Role   string
	Name   string
	QuizID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	Clients         map[string]*Client
	Hosts           map[string]*Client            // Map quizId -> host
	Players         map[string]map[string]*Client // Map quizId -> playerID -> player
	Broadcast       chan []byte
	Register        chan *Client
	Unregister      chan *Client
	mu              sync.Mutex
	QuizService     *service.QuizService
	QuestionService *service.QuestionService
	AnswerService   *service.AnswerService
}

func NewHub(quizService *service.QuizService, questionService *service.QuestionService, answerService *service.AnswerService) *Hub {
	return &Hub{
		Clients:         make(map[string]*Client),
		Hosts:           make(map[string]*Client),
		Players:         make(map[string]map[string]*Client),
		Broadcast:       make(chan []byte),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		QuizService:     quizService,
		QuestionService: questionService,
		AnswerService:   answerService,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.handleRegister(client)

		case client := <-h.Unregister:
			h.handleUnregister(client)

		case message := <-h.Broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Role == "host" {
		h.Hosts[client.QuizID] = client
		fmt.Printf("Host đã kết nối tới quiz %s\n", client.QuizID)
	} else if client.Role == "player" {
		if h.Players[client.QuizID] == nil {
			h.Players[client.QuizID] = make(map[string]*Client)
		}
		h.Players[client.QuizID][client.ID] = client
		h.notifyHost(client.QuizID, "PARTICIPANT_JOINED", client.ID, client.Name)
		fmt.Printf("Player %s đã kết nối tới quiz %s\n", client.ID, client.QuizID)
	}

	h.Clients[client.ID] = client
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Role == "host" {
		delete(h.Hosts, client.QuizID)
		fmt.Printf("Host đã ngắt kết nối khỏi quiz %s\n", client.QuizID)
	} else if client.Role == "player" {
		if players, ok := h.Players[client.QuizID]; ok {
			delete(players, client.ID)
			h.notifyHost(client.QuizID, "PARTICIPANT_LEFT", client.ID, client.Name)
		}
	}

	delete(h.Clients, client.ID)
}

func (h *Hub) handleBroadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client.ID)
		}
	}
}

func (h *Hub) notifyHost(quizID, eventType, playerID string, name string) {
	if host, ok := h.Hosts[quizID]; ok {
		message := Message{
			Type:    eventType,
			Payload: map[string]string{"id": playerID, "name": name},
		}
		msgBytes, _ := json.Marshal(message)
		host.Send <- msgBytes
		fmt.Printf("Thông báo tới host quiz %s: %s\n", quizID, eventType)
	}
}

// handle messages from client
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

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Printf("Lỗi giải mã tin nhắn từ client %s: %v\n", c.ID, err)
			continue
		}

		fmt.Println("Tin nhắn từ client:", msg)

		switch msg.Type {
		case "START_QUIZ":
			if c.Role == "host" {
				duration := 5 // Ví dụ, đếm ngược 5 giây
				go hub.StartCountdown(c.QuizID, duration)
				fmt.Printf("Bắt đầu đếm ngược quiz %s\n", c.QuizID)
			}
		case "REQUEST_NEXT_QUESTION":
			if c.Role == "host" {
				// Lấy danh sách câu hỏi
				questions, err := hub.QuestionService.GetQuestions(context.Background(), c.QuizID)
				if err != nil {
					fmt.Println("Lỗi lấy câu hỏi:", err)
					return
				}
				// Get current question
				currentQuestion := int(msg.Payload.(map[string]interface{})["currentQuestion"].(float64))
				hub.SendQuiz(c.QuizID, questions, currentQuestion) // da +1 khi gui cau hoi
			}
		case "QUIZ_END":
			if c.Role == "host" {
				go hub.SendLeaderBoard(c.QuizID)
				fmt.Printf("Quiz %s đã kết thúc\n", c.QuizID)
				go hub.SendQuizEndResult(c.QuizID)
				fmt.Printf("Gửi kết quả cuối cùng của quiz %s\n", c.QuizID)
			}
		case "SUBMIT_ANSWER":
			if c.Role == "player" {
				// Lưu câu trả lời của người chơi
				questionId := strconv.Itoa(int(msg.Payload.(map[string]interface{})["questionId"].(float64)))
				answer := msg.Payload.(map[string]interface{})["answer"].(string)
				err := cache.StorePlayerAnswer(c.QuizID, c.ID, questionId, time.Now().UnixMilli(), answer)
				if err != nil {
					fmt.Println("Lỗi lưu câu trả lời:", err)
					return
				}
				fmt.Printf("Người chơi đã trả lời câu hỏi %s\n", c.QuizID)

				err = cache.UpdateCounterPlayerAnswer(c.QuizID, questionId, answer)
				if err != nil {
					fmt.Println("Lỗi cập nhật số lần trả lời:", err)
					return
				}
			}
		default:
			hub.Broadcast <- message
		}
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

// Countdown logic
func (h *Hub) StartCountdown(quizID string, duration int) {
	// check status of quiz before start
	quiz, err := h.QuizService.GetQuiz(context.Background(), quizID)
	if err != nil {
		fmt.Println("Lỗi lấy quiz:", err)
		return
	}

	if quiz.Status != reference.DRAFT {
		fmt.Println("Quiz không ở trạng thái DRAFT")
		return
	}

	// change status of quiz to started
	_, err = h.QuizService.ChangeQuizStatus(context.Background(), quizID, reference.LIVE)
	if err != nil {
		fmt.Println("Lỗi thay đổi trạng thái quiz:", err)
		return
	}

	message := Message{
		Type: "QUIZ_START_COUNTDOWN",
		Payload: map[string]string{
			"quizId":    quizID,
			"startTime": fmt.Sprintf("%d", duration),
		},
	}
	h.broadcastToQuiz(quizID, message)
	time.Sleep(time.Duration(duration) * time.Second)
	// Start quiz)
	questions, err := h.QuestionService.GetQuestions(context.Background(), quizID)
	if err != nil {
		fmt.Println("Lỗi lấy câu hỏi:", err)
		return
	}
	h.SendQuiz(quizID, questions, 0)
}

func (h *Hub) SendQuiz(quizID string, questions []*entity.Question, currentQuestion int) {
	totalQuestions := len(questions)

	// Bắt đầu vòng lặp câu hỏi
	go func() {
		// for currentQuestion < totalQuestions {
		// Lấy câu hỏi hiện tại
		question := questions[currentQuestion]

		// get answers
		answers, err := h.AnswerService.GetAnswers(context.Background(), question.ID)
		if err != nil {
			fmt.Println("Lỗi lấy câu trả lời:", err)
			return
		}
		answersText := make([]string, len(answers))
		for i, answer := range answers {
			answersText[i] = answer.AnswerText
		}

		// Store question time to redis
		err = cache.StoreQuestionTime(quizID, strconv.Itoa(int(question.ID)), time.Now().UnixMilli())
		if err != nil {
			fmt.Println("Lỗi lưu thời gian câu hỏi:", err)
			return
		}

		// Gửi câu hỏi tới client
		h.SendNextQuestion(
			quizID,
			int(question.ID),
			currentQuestion+1,
			totalQuestions,
			question.QuestionText,
			// question.Answers,
			answersText,
			int(question.TimeLimit),
		)

		// Tạo goroutine để chờ hết thời gian trả lời
		go func(quizID string, question *entity.Question) {
			time.Sleep(time.Duration(question.TimeLimit) * time.Second)

			// Gửi kết quả tự động
			// h.SendQuestionResult(quizID, question.CorrectAnswer, question.UserResponses)
			var correctAnswer string
			var userResponses = make(map[string]int)
			for _, answer := range answers {
				counter, err := cache.GetCounterPlayerAnswer(quizID, strconv.Itoa(int(question.ID)), answer.AnswerText)
				if err != nil {
					userResponses[answer.AnswerText] = 0
				} else {
					userResponses[answer.AnswerText] = int(counter)
				}
				if answer.IsCorrect {
					correctAnswer = answer.AnswerText
				}
			}
			fmt.Println("Câu trả lời:", userResponses)
			h.SendQuestionResultToHost(quizID, correctAnswer, userResponses)

			// Tính điểm cho người chơi
			questionStartTime, err := cache.GetQuestionTime(quizID, strconv.Itoa(int(question.ID)))
			if err != nil {
				fmt.Println("Lỗi lấy thời gian câu hỏi:", err)
				return
			}
			if players, ok := h.Players[quizID]; ok {
				for playerId, _ := range players {
					playerTime, playerAnswer, err := cache.GetPlayerAnswer(quizID, playerId, strconv.Itoa(int(question.ID)))
					if err != nil {
						fmt.Println("Lỗi lấy câu trả lời của người chơi1: ", quizID+" - "+playerId+" - "+strconv.Itoa(int(question.ID)))
						continue
					}
					if playerAnswer == correctAnswer {
						points := 1000 - (playerTime-questionStartTime)/int64(question.TimeLimit)
						err = cache.UpdatePlayerScore(quizID, playerId, points)
						if err != nil {
							fmt.Println("Lỗi cập nhật điểm cho người chơi:", err)
							continue
						}
					}
				}
			}

			if players, ok := h.Players[quizID]; ok {
				for playerId, client := range players {
					playerTime, playerAnswer, err := cache.GetPlayerAnswer(quizID, playerId, strconv.Itoa(int(question.ID)))
					if err != nil {
						fmt.Println("Lỗi lấy câu trả lời của người chơi2: ", quizID+" - "+playerId+" - "+strconv.Itoa(int(question.ID)))
						continue
					}
					if playerAnswer == correctAnswer {
						points := 1000 - (playerTime-questionStartTime)/int64(question.TimeLimit)
						totalPoints, err := cache.GetPlayerScore(quizID, playerId)
						if err != nil {
							fmt.Println("Lỗi lấy tổng điểm của người chơi:", err)
							continue
						}
						ranking, err := cache.GetPlayerRanking(quizID, playerId)
						if err != nil {
							fmt.Println("Lỗi lấy xếp hạng của người chơi:", err)
							continue
						}
						h.SendQuestionResultToPlayer(quizID, client, currentQuestion, true, int(points), int(totalPoints), int(ranking))
					} else {
						totalPoints, err := cache.GetPlayerScore(quizID, playerId)
						if err != nil {
							fmt.Println("Lỗi lấy tổng điểm của người chơi:", err)
							continue
						}
						ranking, err := cache.GetPlayerRanking(quizID, playerId)
						if err != nil {
							fmt.Println("Lỗi lấy xếp hạng của người chơi:", err)
							continue
						}
						h.SendQuestionResultToPlayer(quizID, client, currentQuestion, false, 0, int(totalPoints), int(ranking))
					}
				}
			}
		}(quizID, question)

		// Đợi hết thời gian trước khi chuyển câu hỏi tiếp theo
		// time.Sleep(time.Duration(questions[currentQuestion].TimeLimit) * time.Second)

		// currentQuestion++
		// }

		// Sau khi hết câu hỏi, gửi bảng xếp hạng
		// h.SendLeaderBoard(quizID)
	}()
}

func (h *Hub) SendNextQuestion(quizID string, questionId int, currentQuestion, totalQuestions int, questionContent string, answers []string, timeLimit int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Tạo payload của câu hỏi
	message := Message{
		Type: "QUESTION_REQUEST",
		Payload: map[string]interface{}{
			"quizId":          quizID,
			"questionId":      questionId,
			"question":        questionContent,
			"answers":         answers,
			"currentQuestion": currentQuestion,
			"totalQuestions":  totalQuestions,
			"timeLimit":       timeLimit,
		},
	}

	messageBytes, _ := json.Marshal(message)

	// Gửi câu hỏi tới host của quiz
	if host, ok := h.Hosts[quizID]; ok {
		host.Send <- messageBytes
	}

	// Gửi câu hỏi tới tất cả các client trong quiz
	if players, ok := h.Players[quizID]; ok {
		for _, player := range players {
			player.Send <- messageBytes
		}
	}

	fmt.Printf("Gửi câu hỏi tới quiz %s: %v\n", quizID, message)
}

func (h *Hub) SendQuestionResultToHost(quizID string, answer string, userResponses map[string]int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Đếm số người chọn từng đáp án
	// answerStats := make(map[string]int)
	// for _, userAnswer := range userResponses {
	// 	answerStats[userAnswer]++
	// }

	// convert map to array
	// answerStats := make([]int, len(userResponses))
	i := 0
	totalAnswered := 0
	for _, userAnswer := range userResponses {
		// answerStats[i] = userAnswer
		totalAnswered += userAnswer
		i++
	}

	// Tạo payload kết quả
	message := Message{
		Type: "QUESTION_RESULT",
		Payload: map[string]interface{}{
			"quizId":        quizID,
			"correctAnswer": answer,
			"totalAnswered": totalAnswered,
			"answerStats":   userResponses,
		},
	}
	messageBytes, _ := json.Marshal(message)

	// Gui ket qua cau hoi toi host
	if host, ok := h.Hosts[quizID]; ok {
		host.Send <- messageBytes
	}

	fmt.Printf("Gửi kết quả câu hỏi tới host quiz %s: %v\n", quizID, message)
}

func (h *Hub) SendQuestionResultToPlayer(quizID string, player *Client, currentQuestion int, isCorrect bool, points int, totalPoints int, ranking int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Tạo payload kết quả
	message := Message{
		Type: "PLAYER_QUESTION_RESULT",
		Payload: map[string]interface{}{
			"quizId":          quizID,
			"currentQuestion": currentQuestion,
			"isCorrect":       isCorrect,
			"points":          points,
			"totalPoints":     totalPoints,
			"ranking":         ranking + 1,
		},
	}
	messageBytes, _ := json.Marshal(message)

	// Gửi kết quả tới player
	player.Send <- messageBytes

	fmt.Printf("Gửi kết quả câu hỏi tới player quiz %s: %v\n", quizID, message)
}

func (h *Hub) SendLeaderBoard(quizID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// change status of quiz to started
	_, err := h.QuizService.ChangeQuizStatus(context.Background(), quizID, reference.ENDED)
	if err != nil {
		fmt.Println("Lỗi thay đổi trạng thái quiz:", err)
		return
	}

	// Giả sử có dữ liệu điểm số từ userResponses hoặc một hệ thống lưu trữ
	leaderboard := []map[string]interface{}{}

	// if players, ok := h.Players[quizID]; ok {
	// 	for playerID := range players {
	// 		// Ví dụ: Điểm của mỗi người chơi được tính ở đây (cần tích hợp logic scoring)
	// 		score := 100
	// 		leaderboard = append(leaderboard, map[string]interface{}{
	// 			"playerId":   playerID,
	// 			"playerName": "Player",
	// 			"score":      score,
	// 		})
	// 	}
	// }

	cacheLeaderboard, err := cache.GetLeaderboard(quizID)
	if err != nil {
		fmt.Println("Lỗi lấy bảng xếp hạng:", err)
		return
	}
	mapLeaderboard, err := cache.MutilGetPlayerScore(quizID, cacheLeaderboard)
	if err != nil {
		fmt.Println("Lỗi lấy điểm của người chơi:", err)
		return
	}
	for _, playerID := range cacheLeaderboard {
		playerInfo, err := cache.GetPlayerSession(playerID)
		if err != nil {
			fmt.Println("Lỗi lấy thông tin người chơi:", err)
			continue
		}
		// convert string to json
		var player map[string]interface{}
		err = json.Unmarshal([]byte(playerInfo), &player)
		if err != nil {
			fmt.Println("Lỗi chuyển đổi thông tin người chơi:", err, " - ", playerInfo)
			continue
		}
		playerName := player["name"].(string)
		leaderboard = append(leaderboard, map[string]interface{}{
			"playerId":   playerID,
			"playerName": playerName,
			"score":      mapLeaderboard[playerID],
		})
	}

	// Tạo payload LEADER_BOARD
	message := Message{
		Type: "LEADER_BOARD",
		Payload: map[string]interface{}{
			"quizId":      quizID,
			"leaderboard": leaderboard,
		},
	}
	messageBytes, _ := json.Marshal(message)

	// Gửi LEADER_BOARD tới host
	if host, ok := h.Hosts[quizID]; ok {
		host.Send <- messageBytes
	}
	fmt.Printf("Gửi bảng xếp hạng tới host quiz %s: %v\n", quizID, message)
}

func (h *Hub) SendQuizEndResult(quizID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// calculate score for each player
	leaderboard := []map[string]interface{}{}
	if players, ok := h.Players[quizID]; ok {
		for playerID := range players {
			// Ví dụ: Điểm của mỗi người chơi được tính ở đây (cần tích hợp logic scoring)
			score, err := cache.GetPlayerScore(quizID, playerID)
			if err != nil {
				fmt.Println("Lỗi lấy điểm của người chơi:", err)
				continue
			}
			leaderboard = append(leaderboard, map[string]interface{}{
				"playerId": playerID,
				"score":    score,
			})

			ranking, err := cache.GetPlayerRanking(quizID, playerID)
			if err != nil {
				fmt.Println("Lỗi lấy xếp hạng của người chơi:", err)
				continue
			}

			// send result to player
			message := Message{
				Type: "QUIZ_END_RESULT",
				Payload: map[string]interface{}{
					"quizId":      quizID,
					"score":       score,
					"ranking":     ranking + 1,
					"totalPlayer": len(players),
				},
			}
			msgBytes, _ := json.Marshal(message)
			players[playerID].Send <- msgBytes
		}
	}

}

func (h *Hub) broadcastToPlayer(quizID string, message Message) {
	msgBytes, _ := json.Marshal(message)
	if players, ok := h.Players[quizID]; ok {
		for _, player := range players {
			player.Send <- msgBytes
		}
	}
}

func (h *Hub) broadcastToHost(quizID string, message Message) {
	msgBytes, _ := json.Marshal(message)
	if host, ok := h.Hosts[quizID]; ok {
		host.Send <- msgBytes
	}
}

func (h *Hub) broadcastToQuiz(quizID string, message Message) {
	msgBytes, _ := json.Marshal(message)
	if host, ok := h.Hosts[quizID]; ok {
		host.Send <- msgBytes
	}
	if players, ok := h.Players[quizID]; ok {
		for _, player := range players {
			player.Send <- msgBytes
		}
	}
}
