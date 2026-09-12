package formatter

import "encoding/json"

type Student struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}
type Formatter interface{ Format(Student) (string, error) }

type TextFormatter struct{}

func (TextFormatter) Format(s Student) (string, error) {
	// TODO: return "Name: Score"
	return "", nil
}

type JSONFormatter struct{}

func (JSONFormatter) Format(s Student) (string, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}
