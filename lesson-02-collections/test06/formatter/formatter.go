package formatter

import "encoding/json"

// Student 是格式化器的输入。
type Student struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// Formatter 只描述“能把一个 Student 格式化成字符串”。
// 任何类型，只要实现了 Format，就自动满足它。
type Formatter interface {
	Format(Student) (string, error)
}

// TextFormatter 输出人读的文本：Lin: 95。
type TextFormatter struct{}

func (TextFormatter) Format(s Student) (string, error) {
	// TODO: 用 fmt.Sprintf("%s: %d", s.Name, s.Score)。
	return "", nil
}

// JSONFormatter 输出机器读的 JSON。
type JSONFormatter struct{}

func (JSONFormatter) Format(s Student) (string, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}

// FormatAll 按顺序调用每个 formatter，把结果收集起来。
// 遇到第一个错误就返回已经收集到的结果和这个错误。
func FormatAll(formatters []Formatter, s Student) ([]string, error) {
	// TODO: 遍历 formatters，调用 Format，append 结果；err != nil 时立刻返回。
	return nil, nil
}
