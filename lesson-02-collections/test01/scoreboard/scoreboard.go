package scoreboard

// Totals 将每位学生的多次成绩累加到 map；输入不应被修改。
func Totals(records map[string][]int) map[string]int {
	// TODO: make a result map and sum each slice
	return nil
}

// Top 返回分数最高的前 n 名，分数相同时按姓名升序。
type Entry struct {
	Name  string
	Score int
}

func Top(totals map[string]int, n int) []Entry {
	// TODO: copy to a slice, sort, and clamp n
	return nil
}
