package counter

// Counter 是一个并发安全的计数器，用 sync.Mutex（或 sync.RWMutex）保护。
//
// 要求：
//   - 零值 var c Counter 可以直接使用，不需要构造函数；
//   - 多个 goroutine 同时 Add 时不能丢更新；
//   - Value、Reset 也要和 Add 互斥。
//
// TODO: 从零实现。给结构体加上需要的字段，再补全三个方法。
type Counter struct {
}

// Add 把计数加上 delta（delta 可以是负数）。
func (c *Counter) Add(delta int) {
	// TODO
}

// Value 返回当前计数。
func (c *Counter) Value() int {
	// TODO
	return 0
}

// Reset 把计数清零。
func (c *Counter) Reset() {
	// TODO
}

// AtomicCounter 和 Counter 对外行为完全一样，但内部用 sync/atomic 实现，
// 不用锁。
//
// TODO: 从零实现。
type AtomicCounter struct {
}

// Add 把计数加上 delta。
func (c *AtomicCounter) Add(delta int) {
	// TODO
}

// Value 返回当前计数。
func (c *AtomicCounter) Value() int {
	// TODO
	return 0
}

// Reset 把计数清零。
func (c *AtomicCounter) Reset() {
	// TODO
}
