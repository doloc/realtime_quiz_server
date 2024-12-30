package storage

import "realtime_quiz_server/entity"

type QuestionStorage interface {
	BaseStorage
	GetQuestions(quizID string) ([]*entity.Question, error)
}

func (storage *storage) GetQuestions(quizID string) ([]*entity.Question, error) {
	var questions []*entity.Question
	if err := storage.db.Preload("Quiz").Where(map[string]interface{}{"quiz_id": quizID}).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
