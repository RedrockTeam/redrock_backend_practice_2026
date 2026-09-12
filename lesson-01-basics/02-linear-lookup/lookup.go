package lookup

import "strings"

type Student struct {
	Name  string
	Score int
}

// FindByName 在线性列表中查找名字，找不到时返回 false。
func FindByName(students []Student, name string) (Student, bool) {
	// TODO: implement
	return Student{}, false
}

// FormatStudents 用一行一个学生的形式输出，名字两端空白应被清理。
func FormatStudents(students []Student) string {
	var b strings.Builder
	for _, student := range students {
		// TODO: write one line
		_ = student
	}
	return b.String()
}
