package roster

import (
	"fmt"
	"strings"
)

// Student 是名册里的一条学生记录。
// 三个字段都是大写的，因为别的包需要构造和读取它们。
type Student struct {
	ID    string
	Name  string
	Score int
}

// Roster 是一份按加入顺序保存的学生名册。
//
// 两个字段都是小写的：别的包不能直接改名册内容，
// 只能通过下面这些导出方法操作——这就是 test04 的可见性规则在实际结构里的用法。
type Roster struct {
	students    []Student
	comparisons int
}

// New 返回一份空名册。
func New() *Roster {
	// TODO: 返回一个 *Roster。
	// 切片的零值是 nil，可以直接 append，所以这里不需要预先分配。
	return nil
}

// Add 把一条记录追加到名册末尾。
// 名字两端的空白要先去掉：Add(Student{Name: " Lin "}) 存进去的是 "Lin"。
func (r *Roster) Add(student Student) {
	// TODO: 用 strings.TrimSpace 清理 student.Name，再 append 到 r.students。
}

// Len 返回名册里的记录条数。
func (r *Roster) Len() int {
	// TODO: 返回切片长度。
	return 0
}

// Comparisons 返回「从名册里取出一条记录来比较」累计发生了多少次。
//
// 它不是业务功能，而是这道题的观测口：
// 我们要用它量出线性查找的代价到底是多少。
func (r *Roster) Comparisons() int {
	// TODO: 返回 r.comparisons。
	return 0
}

// FindByID 按学号查找学生，找不到时第二个返回值为 false。
//
// 要求：每检查一条记录，就把 r.comparisons 加 1。
// 也就是说命中第 k 条（从 1 开始数）就增加 k 次，完全找不到就增加 Len() 次。
func (r *Roster) FindByID(id string) (Student, bool) {
	// TODO: 用 for ... range 遍历 r.students。
	// 每轮循环先给 r.comparisons 加 1，再判断 ID 是否相等。
	// 相等就返回这条记录和 true；循环结束仍未找到就返回零值和 false。
	return Student{}, false
}

// TotalScoreByName 返回所有叫这个名字的学生的分数之和。
// 名册里允许出现同名的人（学号不同），所以这里必须把整份名册看完。
//
// 同样要求：每检查一条记录，就把 r.comparisons 加 1。
func (r *Roster) TotalScoreByName(name string) int {
	// TODO: 遍历整份名册累加分数，记得同步累加 r.comparisons。
	return 0
}

// FormatText 把名册渲染成人看的纯文本，每行一条记录：
//
//	S001 Lin 90
//	S002 Wang 80
//
// 每行以 "\n" 结尾；空名册返回空字符串 ""。
func FormatText(r *Roster) string {
	var builder strings.Builder
	// TODO: 遍历 r.students，用 fmt.Sprintf("%s %s %d\n", ...) 拼出每一行，
	// 再用 builder.WriteString 写进去，最后 return builder.String()。
	builder.WriteString(fmt.Sprintf("%d 条记录\n", r.Len()))
	return builder.String()
}

// FormatCSV 把同一份名册渲染成机器读的 CSV：
//
//	id,name,score
//	S001,Lin,90
//	S002,Wang,80
//
// 第一行固定是表头 "id,name,score"；每行以 "\n" 结尾；
// 空名册只返回表头那一行。
func FormatCSV(r *Roster) string {
	var builder strings.Builder
	// TODO: 先写表头，再遍历 r.students 写出每一行。
	return builder.String()
}
