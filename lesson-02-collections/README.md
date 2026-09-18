# 第二课练习：map 与 interface

承接第一课 `test06`：你在那里量到了线性查找的代价和两份重复的格式化循环，这一课把它们消掉。

| 题目 | 主题 |
|---|---|
| [test01](test01/) | 用 `map` 统计每位学生的总分并按规则取 Top |
| [test02](test02/) | 用同一个 `Formatter` 接口实现文本和 JSON 输出 |
| [test03](test03/) | 观察 `map` 的并发问题 |

`test03` 故意没有同步保护：普通运行可能「碰巧通过」，`go test -race` 会报告数据竞争。**不要在本题仓促加锁**，先把竞态报告里的读写位置记下来，第三课解释 goroutine 与共享内存，第四课再来修。

```bash
go test ./lesson-02-collections/...
go test -race ./lesson-02-collections/test03/...
```
