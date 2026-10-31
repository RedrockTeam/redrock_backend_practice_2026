package cache

// Cache 是一张「字符串 -> 整数」的本地缓存。
//
// 现在它没有锁：这一课只在一个执行流里使用它。
// 下一课会启动多个 goroutine；第四课再决定怎么保护这份共享内存。
type Cache struct {
	values map[string]int
}

// New 返回一份空缓存。
func New() *Cache {
	return &Cache{values: make(map[string]int)}
}

// Put 写入一个键值对；零值 Cache 也应该能用。
func (c *Cache) Put(key string, value int) {
	if c.values == nil {
		c.values = make(map[string]int)
	}
	c.values[key] = value
}

// Get 读取一个键；不存在时返回 0 和 false。
func (c *Cache) Get(key string) (int, bool) {
	value, ok := c.values[key]
	return value, ok
}

// Delete 删除一个键；键不存在时什么也不做。
func (c *Cache) Delete(key string) {
	delete(c.values, key)
}

// Len 返回缓存里有多少个键。
func (c *Cache) Len() int {
	return len(c.values)
}

// Snapshot 返回当前内容的副本。
// 调用方修改返回的 map 不应该影响 Cache。
func (c *Cache) Snapshot() map[string]int {
	// TODO: 创建新 map，把 c.values 的键值逐个复制进去。
	return c.values
}
