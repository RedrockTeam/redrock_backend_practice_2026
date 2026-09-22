package visibility

// adultAge 是包内的小写常量。
const adultAge = 18

// stage 是包内的小写函数，定义在 stage.go，
// 但 user.go 里的 Stage 方法可以直接调用它——因为它们属于同一个包。
func stage(age int) string {
	// TODO: age 达到 adultAge 时返回 "adult"，否则返回 "minor"。
	return ""
}
