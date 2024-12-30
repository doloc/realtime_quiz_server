package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/storage"
)

type AnswerService struct {
	storage storage.AnswerStorage
}

func NewAnswerService(storage storage.AnswerStorage) *AnswerService {
	return &AnswerService{
		storage: storage,
	}
}

func (service *AnswerService) CreateAnswer(ctx context.Context, data *entity.Answer) (*entity.Answer, error) {
	var (
		err    error
		result entity.Answer
	)
	tx := service.storage.BeginTx()
	defer service.storage.CloseTx(tx, err)
	if err := service.storage.Create(tx, data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (service *AnswerService) GetAnswers(ctx context.Context, questionId int64) ([]*entity.Answer, error) {
	result, err := service.storage.GetAnswers(questionId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
