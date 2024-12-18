package reference

import (
	"errors"
	"fmt"
	"strconv"
)

type QuizStatus int

const quizStatusOffset = 1

const (
	DRAFT QuizStatus = iota + quizStatusOffset
	LIVE
	ENDED
)

var quizStatusList = [3]string{"DRAFT", "LIVE", "ENDED"}

func (status *QuizStatus) String() string {
	return quizStatusList[*status-quizStatusOffset]
}

func ParseStr2QuizStatus(s string) (QuizStatus, error) {
	for i, status := range quizStatusList {
		if status == s {
			return QuizStatus(i + quizStatusOffset), nil
		}
	}
	return QuizStatus(quizStatusOffset), errors.New(fmt.Sprintf("Failed to parse QuizStatus: %v", s))
}

func (status *QuizStatus) Scan(value interface{}) error {

	if value == nil {
		return errors.New(fmt.Sprintf("Failed to scan data: %v", value))
	}

	*status = QuizStatus(value.(int64))

	return nil
}

func (status *QuizStatus) Value() (interface{}, error) {
	if status == nil {
		return nil, nil
	}
	return status.String(), nil
}

func (status *QuizStatus) MarshalJSON() ([]byte, error) {
	if status == nil {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, status.String())), nil
}

func (status *QuizStatus) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	v, err := ParseStr2QuizStatus(str)
	if err != nil {
		return err
	}

	*status = v

	return nil
}
