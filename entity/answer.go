package entity

const TableNameAnswer = "answers"

type Answer struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement:true"`
	QuestionID int64     `json:"question_id" gorm:"column:question_id;not null"`
	Question   *Question `json:"question" gorm:"foreignKey:QuestionID;references:ID"`
	AnswerText string    `json:"answer_text" gorm:"column:answer_text;not null"`
	IsCorrect  bool      `json:"is_correct" gorm:"column:is_correct;not null"`
}

func (*Answer) TableName() string {
	return TableNameAnswer
}
