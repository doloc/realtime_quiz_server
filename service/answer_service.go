package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/storage"
)

type answerService struct {
	storage storage.AnswerStorage
}

func NewAnswerService(storage storage.AnswerStorage) *answerService {
	return &answerService{
		storage: storage,
	}
}

func (service *answerService) CreateAnswer(ctx context.Context, data *entity.Answer) (*entity.Answer, error) {
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
