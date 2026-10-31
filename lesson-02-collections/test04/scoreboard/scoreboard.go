package scoreboard

// Totals 把每位学生的多次成绩累加到一张新 map。
// records 本身不能被修改；某个学生没有成绩时，总分为 0。
func Totals(records map[string][]int) map[string]int {
	// TODO: 创建结果 map，遍历 records，把每条成绩加起来。
	// 不要直接改 records 里的切片。
	return nil
}

// Entry 是排行榜上的一行。
type Entry struct {
	Name  string
	Score int
}

// Top 返回分数最高的前 n 名。
// 分数相同时按姓名升序；n <= 0 时返回空结果；n 大于人数时返回全部。
func Top(totals map[string]int, n int) []Entry {
	// TODO: 把 map 转成切片，按「分数降序、姓名升序」排序，再截取前 n 个。
	return nil
}
