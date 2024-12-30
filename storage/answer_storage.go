package storage

import "realtime_quiz_server/entity"

type AnswerStorage interface {
	BaseStorage
	GetAnswers(questionId int64) ([]*entity.Answer, error)
}

func (storage *storage) GetAnswers(questionId int64) ([]*entity.Answer, error) {
	var answers []*entity.Answer
	if err := storage.db.Preload("Question").Where(map[string]interface{}{"question_id": questionId}).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}
