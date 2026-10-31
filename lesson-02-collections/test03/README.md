# test03：切片别名——五个“看起来能跑”的 bug

对应课件：第二课 · 切片（`append` 容量足够时原地写入、`copy`、`range`）。
题型：**代码纠错**。

## 背景

`test02` 里你已经知道：切片头 = (指针, len, cap)，赋值只复制切片头。这一题把这件事放进五个真实场景里：每个函数都能编译、单独看也像是对的，但只要调用方**继续使用**返回值或原切片，数据就会被悄悄改掉。

## 现象

```bash
go test -v ./lesson-02-collections/test03/aliasing
```

`aliasing/aliasing.go` 里每个函数上方都写了 `BUG:` 注释，说明症状。你要做的是理解为什么，然后修好。

| 函数 | 症状 |
|---|---|
| `Head` | 对返回值 `append`，原切片后面的元素被覆盖 |
| `Branch` | 从同一个前缀 `append` 两次，两条路径变成一样 |
| `RemoveIf` | 相邻两个要删的元素只删掉一个；调用方的数据被挪乱 |
| `Pointers` | 通过返回的指针修改，原切片纹丝不动 |
| `Chunk` | 对第一段 `append`，第二段的内容变了 |

## 任务

只改 `aliasing/aliasing.go` 的函数体，不改签名，不改测试。

## 提示（不是答案）

- 回到课件“容量足够时：原地写入”那一节：`t := append(s, 3)` 和 `u := append(s, 9)` 写的是同一格；
- **完整切片表达式** `s[low:high:max]`：结果的 `len = high-low`、`cap = max-low`。cap 收紧之后，下一次 `append` 只能分配新数组；
- 想要一份完全独立的数据：`make` + `copy`，或 `append([]T(nil), s...)`；
- 边遍历边删除很难写对。换个思路：**新建**一个结果切片，只把要留下的放进去；
- `for _, v := range s` 的 `v` 是副本，`&v` 是副本的地址；要拿元素本身的地址，用下标。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestHeadDoesNotLeakCapacity` | 对返回值 append 不影响原切片 |
| `TestBranchPathsAreIndependent` | 两条路径与前缀互相独立 |
| `TestRemoveIfHandlesAdjacentMatches` | 相邻匹配都被删掉；`nil`、全删 |
| `TestRemoveIfDoesNotMutateInput` | 调用方的切片不被修改 |
| `TestPointersPointIntoTheSlice` | 指针指向切片里的元素本身 |
| `TestChunkPiecesDoNotOverlap` | 每一段 cap 独立；`size <= 0` 返回 `nil` |
| `TestBridgeAppendWritesInPlaceWhenCapAllows` | 观察：课件里 `t`、`u` 互相覆盖的例子（不检查你的实现） |

## 做完以后

五个 bug 其实是同一件事：**两个切片头指向同一块底层数组，而你以为它们是两份数据**。

记住两条规则：交出切片给别人继续 `append` 时收紧 cap 或复制；要元素的地址用下标。下一题 `test04` 换成 map：map 也是句柄，但它不会让你踩 cap 的坑，而是让“按 key 查”变成一步。
