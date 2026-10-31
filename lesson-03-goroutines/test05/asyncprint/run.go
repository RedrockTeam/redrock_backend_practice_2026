package asyncprint

import (
	"fmt"
	"io"
)

// Run 启动 workers 个 goroutine，每个 goroutine 向 out 写一行：
//
//	worker 0
//	worker 1
//	...
//
// 第三课先用 time.Sleep 等它们跑完；第四课会换成 WaitGroup。
// TODO: 现在它是顺序打印的，一个 goroutine 都没有启动。
func Run(out io.Writer, workers int) {
	for i := 0; i < workers; i++ {
		fmt.Fprintf(out, "worker %d\n", i)
	}
}
