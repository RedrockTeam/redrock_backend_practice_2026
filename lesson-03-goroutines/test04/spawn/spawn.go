package spawn

// Spawn 为 0, 1, ..., n-1 每个编号启动一个 goroutine 执行 work(id)，
// 然后**立刻返回**，不等待任何一个 work 结束。
//
// TODO: 现在它是顺序调用的：只要有一个 work 卡住，Spawn 就永远回不来。
func Spawn(n int, work func(id int)) {
	for i := 0; i < n; i++ {
		work(i)
	}
}

// Measure 调用 Spawn(n, work)，返回调用前后 runtime.NumGoroutine() 的增量。
func Measure(n int, work func(id int)) int {
	// TODO: 用 runtime.NumGoroutine() 读出调用前后的 goroutine 数量。
	return 0
}

// Chunk 是左闭右开区间 [Start, End)。
type Chunk struct {
	Start, End int
}

// Chunks 把 [0, n) 切成每段 size 个的区间，最后一段可以不满。
//
//	Chunks(10, 3) = [{0 3} {3 6} {6 9} {9 10}]
//
// n <= 0 或 size <= 0 时返回 nil。
func Chunks(n, size int) []Chunk {
	// TODO
	return nil
}

// SpawnChunks 为 Chunks(n, size) 的每一段启动一个 goroutine 执行 work(start, end)，
// 立刻返回启动了多少个 goroutine。
func SpawnChunks(n, size int, work func(start, end int)) int {
	// TODO
	return 0
}
