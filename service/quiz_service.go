package service

import (
	"context"
	"realtime_quiz_server/entity"
	"realtime_quiz_server/entity/reference"
	"realtime_quiz_server/storage"
)

type QuizService struct {
	storage storage.QuizStorage
}

func NewQuizService(storage storage.QuizStorage) *QuizService {
	return &QuizService{
		storage: storage,
	}
}

func (service *QuizService) CreateQuiz(ctx context.Context, data *entity.Quiz) (*entity.Quiz, error) {
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

func (service *QuizService) GetQuiz(ctx context.Context, quizID string) (*entity.Quiz, error) {
	var result *entity.Quiz
	cond := map[string]interface{}{"id": quizID}
	err := service.storage.Get(ctx, cond, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *QuizService) ChangeQuizStatus(ctx context.Context, quizID string, status reference.QuizStatus) (*entity.Quiz, error) {
	result, err := service.storage.ChangeQuizStatus(quizID, status)
	if err != nil {
		return nil, err
	}
	return result, nil
}
