# 第二课练习：map 与 interface

1. 用 map 统计每位学生的总分并按规则取 Top；
2. 用同一个 `Formatter` 接口实现文本和 JSON 输出；
3. 运行 `03-map-race`，记录普通测试和 `go test -race` 的差异。

第三题故意没有同步保护。不要在本题仓促加锁，先把竞态报告中的读写位置记下来，第三课会解释 goroutine 和共享内存。

```bash
go test ./lesson-02-collections/...
go test -race ./lesson-02-collections/03-map-race/...
```
