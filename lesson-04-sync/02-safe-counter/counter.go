package counter

// Counter 的方法会被多个 goroutine 调用；请用合适的锁保护 n。
type Counter struct{ n int }

func (c *Counter) Add(delta int) { c.n += delta }
func (c *Counter) Value() int    { return c.n }
