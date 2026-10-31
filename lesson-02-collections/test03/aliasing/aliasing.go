package aliasing

// 本文件里的五个函数都“看起来能跑”，但每一个都藏着一个切片共享底层数组的 bug。
// 先跑测试看现象，再逐个修。只改函数体，不要改函数签名。

// Head 返回 values 的前 n 个元素（n 超过长度时返回全部）。
//
// 调用方可能会对返回值继续 append，这不能改写 values 后面的元素。
func Head(values []int, n int) []int {
	if n > len(values) {
		n = len(values)
	}
	// BUG: 返回的窗口 cap 仍然延伸到 values 的末尾，
	// 调用方一 append，就原地写进了 values[n]。
	return values[:n]
}

// Branch 在同一个前缀 base 后面分别追加 a 和 b，返回两条互相独立的路径。
//
// 例如 base = [/ home]，Branch(base, "alice", "bob")
// 应得到 [/ home alice] 和 [/ home bob]。
func Branch(base []string, a, b string) ([]string, []string) {
	// BUG: base 的 cap 有富余时，两次 append 写的是同一个位置，
	// 第二次会把第一次的结果覆盖掉。
	left := append(base, a)
	right := append(base, b)
	return left, right
}

// RemoveIf 返回删掉所有满足 drop 的元素之后的新切片。
//
// 调用方传进来的 values 不能被修改。
func RemoveIf(values []int, drop func(int) bool) []int {
	// BUG 1: 边遍历边删除，删掉 i 之后下一个元素挪到了 i，却被 i++ 跳过。
	// BUG 2: append(values[:i], values[i+1:]...) 是原地挪动，调用方的数据被改乱了。
	for i := 0; i < len(values); i++ {
		if drop(values[i]) {
			values = append(values[:i], values[i+1:]...)
		}
	}
	return values
}

// Student 是一条名册记录。
type Student struct {
	Name  string
	Score int
}

// Pointers 返回指向 students 中每个元素的指针。
//
// 通过返回的指针修改 Score，students 里的对应元素要跟着变。
func Pointers(students []Student) []*Student {
	result := make([]*Student, 0, len(students))
	// BUG: range 给出的 student 是元素的副本，&student 指向的是副本，不是切片里的那一格。
	for _, student := range students {
		result = append(result, &student)
	}
	return result
}

// Chunk 把 values 按每段 size 个切开（最后一段可以不足 size）。size <= 0 时返回 nil。
//
// 调用方对某一段 append，不能改写下一段的内容。
func Chunk(values []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	var chunks [][]int
	for start := 0; start < len(values); start += size {
		end := start + size
		if end > len(values) {
			end = len(values)
		}
		// BUG: 每一段的 cap 都延伸到 values 末尾，和后面几段共享同一块底层数组。
		chunks = append(chunks, values[start:end])
	}
	return chunks
}
