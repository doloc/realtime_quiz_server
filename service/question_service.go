package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/storage"
)

type QuestionService struct {
	storage storage.QuestionStorage
}

func NewQuestionService(storage storage.QuestionStorage) *QuestionService {
	return &QuestionService{
		storage: storage,
	}
}

func (service *QuestionService) CreateQuestion(ctx context.Context, data *entity.Question) (*entity.Question, error) {
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

func (service *QuestionService) GetQuestions(ctx context.Context, quizId string) ([]*entity.Question, error) {
	result, err := service.storage.GetQuestions(quizId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
