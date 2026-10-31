package lazyfetch

import "context"

// Config 是客户端第一次使用时才去加载的配置。
type Config struct {
	Endpoint string
}

// Client 是一个访问「慢服务」的客户端：
//
//   - 配置很贵，只在第一次用到时加载一次（懒加载）；
//   - 每次 Fetch 都在后台 goroutine 里调用慢服务，调用方可以用 ctx 设超时。
//
// 这段代码在单个 goroutine、服务又很快的时候看起来完全正常。
// 它至少藏着三个问题，测试会逐个把它们逼出来。
// 只改这个文件；New、Config、Fetch 的签名不要变。
type Client struct {
	load  func() Config
	fetch func(endpoint, key string) string

	cfg *Config
}

// New 创建客户端。load 负责加载配置；fetch 负责真正去请求慢服务。
func New(load func() Config, fetch func(endpoint, key string) string) *Client {
	return &Client{load: load, fetch: fetch}
}

// Config 返回配置；第一次调用时才执行 load，之后都复用同一份。
// 多个 goroutine 同时第一次调用时，load 也只能被执行一次。
func (c *Client) Config() Config {
	if c.cfg == nil {
		cfg := c.load()
		c.cfg = &cfg
	}
	return *c.cfg
}

// Fetch 用当前配置的 Endpoint 请求 key。
//
//   - 服务在 ctx 结束前返回：返回结果和 nil；
//   - ctx 先结束：立刻返回 "" 和一个错误。这个错误的信息里要带上 key，
//     并且 errors.Is(err, ctx.Err()) 为真；
//   - 不管哪种情况，后台 goroutine 最终都要能退出。
func (c *Client) Fetch(ctx context.Context, key string) (string, error) {
	cfg := c.Config()

	result := make(chan string)
	go func() {
		result <- c.fetch(cfg.Endpoint, key)
	}()

	select {
	case value := <-result:
		return value, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
