package pitfalls

import "fmt"

// ValidationError 描述一个不合法的字段。
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

// User 是要校验的输入。
type User struct {
	Name string
	Age  int
}

// Validate 校验 u：Name 不能为空，Age 必须在 [0, 150] 内。
// 校验通过时返回 nil；失败时返回 *ValidationError，Field 是 "name" 或 "age"。
// 两个字段都不合法时，报告 name。
//
// BUG: 调用方写 if err := Validate(u); err != nil，合法的 User 也会走进错误分支。
// 先写个小程序打印 err == nil 看看，再修。
func Validate(u User) error {
	var problem *ValidationError
	if u.Name == "" {
		problem = &ValidationError{Field: "name", Reason: "must not be empty"}
	} else if u.Age < 0 || u.Age > 150 {
		problem = &ValidationError{Field: "age", Reason: "out of range"}
	}
	return problem
}
