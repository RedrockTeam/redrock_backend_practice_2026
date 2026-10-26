# test05：从 channel 到字节流（第五课预告）

对应课件：第四课（channel、goroutine、关闭与取消）+ **第五课**（网络、TCP、HTTP）。
题型：**桥接题**：先用第四课的工具处理一个“像网络连接”的东西，然后亲手撞上第五课要解决的问题。

这题不改课件，也不要求你提前学 TCP。你只需要知道：`net.Pipe()` 返回一对**内存里的连接**，两端都能 `Read` / `Write` / `Close`。它和真实网线的共同点是——**它们都只是一根字节流，没有“消息”这个概念**。

## 先理解“字节流”

在 channel 里，发一个 `int` 就是一个 `int`，边界天然存在。连接不是这样：

```text
发送方 Write("alpha\nbeta\ngamma\n")
                    │
                    ▼
接收方 Read(buf)  ← 可能一次拿到 "alpha\nbeta\ngamma\n"
                    也可能先拿到 "alp"，下一次才拿到 "ha\nbeta\n..."
```

TCP 不保存消息边界。这就是“粘包 / 半包”的根源。HTTP 用请求行、头、`Content-Length` 或 chunked 编码来划分消息；UDP 则保留数据报边界。第五课会展开讲。

## 任务

只改 `wire/wire.go`，实现：

```go
type Conn interface {
    io.Reader
    io.Writer
    io.Closer
}

func Serve(conn Conn) error
func Request(conn Conn, text string) (string, error)
```

### 协议

请求和响应都是**一行文本**，以 `\n` 结束：

```text
请求：hello\n
响应：ECHO HELLO\n
```

- 服务端：每读到完整的一行，就回一行 `"ECHO " + 大写文本`；
- 客户端：发送一行，读取**恰好一行**响应（返回值不含 `\n`）；
- 连接正常关闭（`io.EOF`）时 `Serve` 返回 `nil`；
- 要能处理：一次 `Read` 拿到多行、一行被拆成多次 `Read`、多个客户端同时连接。

starter 把“一次 `Read`”当成“一条消息”，所以：

- 单条消息能跑通；
- 一次写入多行会错位；
- 一行拆成两次写入会提前回包。

```bash
go test -v ./lesson-04-sync/test05/wire
```

## 需要的最小新工具

这一题要用标准库 `bufio`。它不改变协议，只是帮你把字节流按行切分：

```go
reader := bufio.NewReader(conn)

line, err := reader.ReadString('\n') // line 含结尾的 '\n'
// err 可能是 io.EOF；line 里可能仍有数据，要先处理数据再处理 err
```

或者：

```go
scanner := bufio.NewScanner(conn)
for scanner.Scan() {
    line := scanner.Text() // 不含换行符
}
if err := scanner.Err(); err != nil { ... }
```

写响应可以用：

```go
fmt.Fprintf(conn, "ECHO %s\n", strings.ToUpper(line))
```

注意：`bufio.Reader` 会预读。严格的一问一答里，每次 `Request` 新建一个 reader 通常也能工作；但如果服务端可能连续发来多条响应，就应该让一个连接共用一个 reader，否则预读到的下一条响应可能随旧 reader 一起丢掉。

## 提示（不是答案）

- 服务端需要一个循环：读一行 → 处理 → 写一行 → 继续；读到 `io.EOF` 才退出；
- 半行要留在缓冲区里等下一次读，不能“收到多少就回多少”；
- 一次读到多行时，要逐行处理；
- 客户端要读到 `\n` 才停，不能假设“一次 Read 就是一条响应”；
- 多个客户端之间不要共享可变状态；每个连接有自己的 reader；
- 如果测试卡住，先检查：写的时候对方在读吗？读的时候对方在写吗？`net.Pipe` 的读写是成对阻塞的。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestRequestResponse` | 连续多次请求/响应，协议正确 |
| `TestServeHandlesSeveralLinesInOneRead` | 一次写入多行，服务端逐行回包 |
| `TestServeBuffersHalfLine` | 一行拆成两次写入，服务端能拼起来 |
| `TestConcurrentClients` | 8 个客户端并发，各自请求互不干扰 |
| `TestServeReturnsOnClose` | 连接关闭后 `Serve` 返回，不泄漏 |
| `TestBridgeOneReadCarriesManyMessages` | 观察原始连接：一次 `Read` 可以拿到多行 |

最后一个是“桥”测试：它不检查你的实现，而是把下一课要解释的现象直接摆在你面前。

## 做完以后：第五课

你现在应该能回答：

1. 为什么 channel 里的“一条消息”，到了连接上就没了？
2. 为什么服务端不能假设“一次 Read = 一条请求”？
3. 如果两台机器要通信，除了“有一根能读写的管道”，还需要知道什么？

第五课的回答：

- 需要地址和端口：`net.Listen("tcp", ":8080")`、`net.Dial("tcp", "host:port")`；
- TCP 建立连接要三次握手，断开要四次挥手；
- TCP 是字节流，所以应用层要自己定边界；HTTP 就是一套已经定好的边界和语义；
- `net.Conn` 正好满足这道题的 `wire.Conn`，你写的 `Serve` 可以几乎原样搬到真实连接上。

```text
第四课  channel 在一个进程内传递“值”
   │
test05  连接是一根字节流，没有消息边界
   │
第五课  IP / 端口 / TCP / HTTP：给字节流加上地址、连接和协议
```
