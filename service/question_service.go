package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/storage"
)

type questionService struct {
	storage storage.QuestionStorage
}

func NewQuestionService(storage storage.QuestionStorage) *questionService {
	return &questionService{
		storage: storage,
	}
}

func (service *questionService) CreateQuestion(ctx context.Context, data *entity.Question) (*entity.Question, error) {
	var (
		err    error
		result entity.Question
	)
	tx := service.storage.BeginTx()
	defer service.storage.CloseTx(tx, err)
	if err := service.storage.Create(tx, data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
