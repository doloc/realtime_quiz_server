package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/storage"
)

type quizService struct {
	storage storage.QuizStorage
}

func NewQuizService(storage storage.QuizStorage) *quizService {
	return &quizService{
		storage: storage,
	}
}

func (service *quizService) CreateQuiz(ctx context.Context, data *entity.Quiz) (*entity.Quiz, error) {
	var (
		err    error
		result entity.Quiz
	)
	tx := service.storage.BeginTx()
	defer service.storage.CloseTx(tx, err)
	if err := service.storage.Create(tx, data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
