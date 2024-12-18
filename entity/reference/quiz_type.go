package reference

import (
	"errors"
	"fmt"
	"strconv"
)

type QuizType int

const quizTypeOffset = 1

const (
	PUBLIC QuizType = iota + quizTypeOffset
	PRIVATE
)

var quizTypeList = [3]string{"PUBLIC", "PRIVATE"}

func (status *QuizType) String() string {
	return quizTypeList[*status-quizTypeOffset]
}

func ParseStr2QuizType(s string) (QuizType, error) {
	for i, status := range quizTypeList {
		if status == s {
			return QuizType(i + quizTypeOffset), nil
		}
	}
	return QuizType(quizTypeOffset), errors.New(fmt.Sprintf("Failed to parse QuizType: %v", s))
}

func (status *QuizType) Scan(value interface{}) error {

	if value == nil {
		return errors.New(fmt.Sprintf("Failed to scan data: %v", value))
	}

	*status = QuizType(value.(int64))

	return nil
}

func (status *QuizType) Value() (interface{}, error) {
	if status == nil {
		return nil, nil
	}
	return status.String(), nil
}

func (status *QuizType) MarshalJSON() ([]byte, error) {
	if status == nil {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, status.String())), nil
}

func (status *QuizType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	v, err := ParseStr2QuizType(str)
	if err != nil {
		return err
	}

	*status = v

	return nil
}
