package entity

import "realtime_quiz_server/entity/reference"

const TableNameQuiz = "quizzes"

type Quiz struct {
	ID          string               `json:"id" gorm:"column:id;primaryKey"`
	Title       string               `json:"title" gorm:"column:title;not null"`
	Description string               `json:"description" gorm:"column:description;not null"`
	Type        reference.QuizType   `json:"type" gorm:"column:type;not null;default:1"`
	Status      reference.QuizStatus `json:"status" gorm:"column:status;not null;default:1"`
	BaseEntity
}

func (*Quiz) TableName() string {
	return TableNameQuiz
}
