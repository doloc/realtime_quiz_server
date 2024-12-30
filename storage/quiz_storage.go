package storage

import (
	"realtime_quiz_server/entity"
	"realtime_quiz_server/entity/reference"
)

type QuizStorage interface {
	BaseStorage
	ChangeQuizStatus(quizID string, status reference.QuizStatus) (*entity.Quiz, error)
}

func (storage *storage) ChangeQuizStatus(quizID string, status reference.QuizStatus) (*entity.Quiz, error) {
	quiz := &entity.Quiz{}
	if err := storage.db.Model(quiz).Where("id = ?", quizID).Update("status", status).Error; err != nil {
		return nil, err
	}
	return quiz, nil
}
