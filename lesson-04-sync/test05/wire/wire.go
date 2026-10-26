package wire

import (
	"fmt"
	"io"
	"strings"
)

// Conn 是这道题使用的连接抽象：能读、能写、能关闭。
// 下一课的 net.Conn 也满足这个接口。
type Conn interface {
	io.Reader
	io.Writer
	io.Closer
}

// Serve 在一个连接上提供回声服务。
//
// 协议：请求和响应都是一行文本，以 '\n' 结束。
// 每收到一行，就返回一行 "ECHO " + strings.ToUpper(那一行)。
// 例如 "hello\n" -> "ECHO HELLO\n"。
// 连接正常关闭（io.EOF）时返回 nil。
//
// TODO: 现在它把“一次 Read”当成“一条消息”。但字节流没有消息边界：
// 一次 Read 可能拿到半行，也可能拿到好几行。请按 '\n' 切分并处理残留。
func Serve(conn Conn) error {
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			line := strings.TrimSuffix(string(buffer[:n]), "\n")
			if _, writeErr := fmt.Fprintf(conn, "ECHO %s\n", strings.ToUpper(line)); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// Request 发送一行文本，并读取一行响应（返回值不含结尾的 '\n'）。
//
// TODO: 现在它只调用一次 Read。请保证它恰好消费“一行”响应，
// 不要把下一条响应也读进来。
func Request(conn Conn, text string) (string, error) {
	if _, err := fmt.Fprintf(conn, "%s\n", text); err != nil {
		return "", err
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(buffer[:n]), "\n"), nil
}
