package storage

import "realtime_quiz_server/entity"

type QuizStorage interface {
	BaseStorage
	GetQuizByID(quizID int) (*entity.Quiz, error)
}

func (storage *storage) GetQuizByID(quizID int) (*entity.Quiz, error) {
	quiz := &entity.Quiz{}
	err := storage.db.First(quiz, quizID).Error
	if err != nil {
		return nil, err
	}
	return quiz, nil
}
