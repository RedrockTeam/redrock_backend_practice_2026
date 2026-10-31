# test03：合并多条车道（补全代码）

对应课件：第四课 · 第二部分（channel、`close`、所有权）、第七部分（fan-out / fan-in）。
题型：**补全已有代码**。

## 现象

`fanin/merge.go` 已经为每个输入启动了一个转发 goroutine，但**没有关闭输出**。消费者读完所有值之后，会一直等一个永远不会到来的“结束信号”。

先跑：

```bash
go test -v ./lesson-04-sync/test03/fanin
```

你会看到测试等到超时，然后提示：

```text
所有输入都结束后输出仍没有关闭；Merge 需要在所有转发 goroutine 退出后 close(out)
```

这不是“慢”，而是**永远等不到**。想一想：消费者怎么知道“不会再有数据了”？

## 任务

只改 `fanin/merge.go`，实现：

```go
func Merge(inputs ...<-chan int) <-chan int
```

约定：

- 返回的 channel 由 `Merge` 负责关闭；
- 只有当**所有输入都关闭、所有转发 goroutine 都退出**后，才能 `close(out)`；
- `Merge` 不能关闭输入 channel：输入只由它们的所有者关闭；
- 没有输入时，返回一个**已经关闭**的 channel；
- 输出的顺序不限，但每个值恰好出现一次。

## 提示（不是答案）

- 一个 `WaitGroup` 可以记录“还有几个转发 goroutine 没退出”；
- `wg.Wait()` 不能直接写在 `Merge` 里：那样 `Merge` 会等所有输入结束才返回，就不是流式 fan-in 了。把它放进一个专门的 goroutine，等完再 `close(out)`；
- `for value := range input` 会在输入关闭时自动结束循环，这是“上游关闭 → 本级退出”的连接点；
- 谁 close 输出？**只有一个地方**。多处 close 会 panic，忘了 close 会让消费者永久等待。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestMergeCollectsEveryValue` | 多路值一个不少、不重复；结束后输出关闭 |
| `TestMergeClosesAfterAllInputs` | 输入都关闭后，输出也关闭 |
| `TestMergeWithNoInputs` | 没有输入时输出立刻关闭 |
| `TestMergeStreamsUnbufferedInputs` | 无缓冲输入下边收边转发，不要求先把所有值缓存起来 |
| `TestMergeHandlesSingleInput` | 只有一路输入时行为也正确 |

## 做完以后

这一题把“数据传递”和“结束信号”分开看：

```text
发送方   发送值，最后 close 自己的输出
Merge    转发值，等所有上游结束，close 公共输出
消费者   range 到输出关闭为止
```

下一题 `test04` 换一个方向：这里是**发送方**全部结束、由 `Merge` 关闭输出；如果是**消费者**提前不读了呢？上游会卡在发送上泄漏，需要 `select` 和 `context` 把取消传上去。
