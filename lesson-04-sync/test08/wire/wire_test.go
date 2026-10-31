package wire_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test08/wire"
)

// newPipe 造一对“像网线一样”的内存连接。
// net.Pipe 的两端都实现了 wire.Conn，但数据只是一段字节流：
// 它不保留消息边界，也不帮你分隔每一条消息。
func newPipe(t *testing.T) (client, server net.Conn) {
	t.Helper()

	client, server = net.Pipe()
	deadline := time.Now().Add(3 * time.Second)
	_ = client.SetDeadline(deadline)
	_ = server.SetDeadline(deadline)

	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	return client, server
}

func startServer(t *testing.T, conn net.Conn) <-chan error {
	t.Helper()

	done := make(chan error, 1)
	go func() {
		done <- wire.Serve(conn)
	}()
	return done
}

func waitServer(t *testing.T, done <-chan error) {
	t.Helper()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) && !errors.Is(err, net.ErrClosed) {
			t.Fatalf("连接关闭后 Serve 返回 %v；正常关闭应当返回 nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("连接关闭后 Serve 没有返回；读循环是否只在收到数据时才检查错误？")
	}
}

func TestRequestResponse(t *testing.T) {
	client, server := newPipe(t)
	done := startServer(t, server)

	for _, testCase := range []struct {
		input string
		want  string
	}{
		{input: "hello", want: "ECHO HELLO"},
		{input: "world", want: "ECHO WORLD"},
		{input: "mixed Case 123", want: "ECHO MIXED CASE 123"},
	} {
		got, err := wire.Request(client, testCase.input)
		if err != nil {
			t.Fatalf("Request(%q) 返回错误: %v", testCase.input, err)
		}
		if got != testCase.want {
			t.Fatalf("Request(%q) = %q, want %q", testCase.input, got, testCase.want)
		}
	}

	_ = client.Close()
	waitServer(t, done)
}

// TestServeHandlesSeveralLinesInOneRead 是 TCP 粘包的预演：
// 对面一次 Write 可能包含三条消息。服务端不能把“一次 Read”当成“一条消息”。
func TestServeHandlesSeveralLinesInOneRead(t *testing.T) {
	client, server := newPipe(t)
	done := startServer(t, server)

	writeErr := make(chan error, 1)
	go func() {
		_, err := io.WriteString(client, "alpha\nbeta\ngamma\n")
		writeErr <- err
	}()

	reader := bufio.NewReader(client)
	for _, want := range []string{"ECHO ALPHA", "ECHO BETA", "ECHO GAMMA"} {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("读取服务端响应失败: %v", err)
		}
		got := strings.TrimSuffix(line, "\n")
		if got != want {
			t.Fatalf("读到 %q, want %q；一次 Read 拿到多行时，必须逐行处理", got, want)
		}
	}

	if err := <-writeErr; err != nil {
		t.Fatalf("写入请求失败: %v", err)
	}
	_ = client.Close()
	waitServer(t, done)
}

// TestServeBuffersHalfLine 是 TCP 半包的预演：
// "hello\n" 可能分两次到达。服务端要等在内存里攒够一整行再处理。
func TestServeBuffersHalfLine(t *testing.T) {
	client, server := newPipe(t)
	done := startServer(t, server)

	go func() {
		_, _ = io.WriteString(client, "hel")
		time.Sleep(20 * time.Millisecond)
		_, _ = io.WriteString(client, "lo\n")
	}()

	reader := bufio.NewReader(client)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("读取服务端响应失败: %v", err)
	}
	got := strings.TrimSuffix(line, "\n")
	if got != "ECHO HELLO" {
		t.Fatalf("读到 %q, want %q；半行要缓存到下一次 Read，而不是收到一半就先回一条", got, "ECHO HELLO")
	}

	_ = client.Close()
	waitServer(t, done)
}

func TestConcurrentClients(t *testing.T) {
	const (
		clientCount  = 8
		requestsEach = 5
	)

	var wg sync.WaitGroup
	errCh := make(chan error, clientCount)

	for c := 0; c < clientCount; c++ {
		c := c
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, server := net.Pipe()
			defer client.Close()
			defer server.Close()

			deadline := time.Now().Add(3 * time.Second)
			_ = client.SetDeadline(deadline)
			_ = server.SetDeadline(deadline)

			done := make(chan error, 1)
			go func() {
				done <- wire.Serve(server)
			}()

			for i := 0; i < requestsEach; i++ {
				text := fmt.Sprintf("client-%d-request-%d", c, i)
				got, err := wire.Request(client, text)
				if err != nil {
					errCh <- fmt.Errorf("client %d: Request(%q): %w", c, text, err)
					return
				}
				want := "ECHO " + strings.ToUpper(text)
				if got != want {
					errCh <- fmt.Errorf("client %d: Request(%q) = %q, want %q", c, text, got, want)
					return
				}
			}

			_ = client.Close()
			select {
			case err := <-done:
				if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) && !errors.Is(err, net.ErrClosed) {
					errCh <- fmt.Errorf("client %d: Serve 返回 %v", c, err)
				}
			case <-time.After(2 * time.Second):
				errCh <- fmt.Errorf("client %d: 关闭连接后 Serve 没有返回", c)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}
}

func TestServeReturnsOnClose(t *testing.T) {
	client, server := newPipe(t)
	done := startServer(t, server)

	_ = client.Close()
	waitServer(t, done)
}

// TestBridgeOneReadCarriesManyMessages 不改学生代码，只观察 net.Pipe 的原始行为：
// 字节流没有消息边界，一次 Read 可能拿到半行，也可能拿到多行。
// 这正是下一课 TCP 要处理的问题，也是 HTTP 为什么要定义请求行、头和 Content-Length。
func TestBridgeOneReadCarriesManyMessages(t *testing.T) {
	client, server := newPipe(t)

	go func() {
		_, _ = io.WriteString(client, "one\ntwo\nthree\n")
	}()

	buffer := make([]byte, 64)
	n, err := server.Read(buffer)
	if err != nil {
		t.Fatalf("读取原始连接失败: %v", err)
	}

	got := string(buffer[:n])
	if !strings.Contains(got, "one\n") || !strings.Contains(got, "two\n") {
		t.Fatalf("一次 Read 收到 %q；这里期望能同时看到多行，以便观察字节流没有消息边界", got)
	}

	t.Logf("一次 Read 收到了 %q", got)
	t.Log("这只是一根字节管道：它不知道哪里是一条消息的结尾。")
	t.Log("下一课的 TCP 也只是一根字节流；HTTP 用请求行、头、Content-Length 来划分消息。")
}
