package entity

const TableNameQuestion = "questions"

type Question struct {
	ID           int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement:true"`
	QuizID       string `json:"quiz_id" gorm:"column:quiz_id;not null"`
	Quiz         *Quiz  `json:"quiz" gorm:"foreignKey:QuizID;references:ID"`
	QuestionText string `json:"question_text" gorm:"column:question_text;not null"`
	TimeLimit    int32  `json:"time_limit" gorm:"column:time_limit;not null;default:30"`
}

func (*Question) TableName() string {
	return TableNameQuestion
}
