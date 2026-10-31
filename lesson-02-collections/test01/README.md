# test01：指针与数组——值传递还是传地址

对应课件：第二课 · 指针和数组。
题型：**补全代码 + 预测输出**（入门）。

## 背景

C 里数组名会退化成指针，传进函数改的就是原数组；Go 不是。Go 里**数组是值**，赋值、传参都会把整个数组复制一份。想让函数改到调用方的变量，就要显式传地址 `&x`，函数里用 `*p` 读写。

这一题把这件事拆成六个小函数和三道预测题。

## 任务

### 1. 补全 `pointer/pointer.go`

| 函数 | 要求 |
|---|---|
| `Inc(p *int)` | 把 `p` 指向的整数加 1；`p == nil` 时什么也不做 |
| `Swap(a, b *int)` | 交换两个指针**指向的值** |
| `Birthday(s *Student)` | 通过结构体指针让 `Age` 加 1 |
| `Renamed(s Student, name string) Student` | 返回改名后的副本，调用方的 `s` 不变 |
| `Scaled(a [3]int, k int) [3]int` | 返回每个元素乘 `k` 的新数组，原数组不变 |
| `ScaleInPlace(a *[3]int, k int)` | 原地把调用方数组的每个元素乘 `k` |

### 2. 预测 `pointer/demo.go` 的输出

`demo.go` 是课件里的三段示例，**不要改它**。只读代码，把你认为的返回值写进 `predict.go` 的三个常量，格式和 `fmt.Sprint` 一致（数组 `"[1 2 3]"`，数字 `"42"`）。

先写预测，再跑测试。猜错了测试会告诉你实际值；改之前先想清楚自己错在哪。

## 提示（不是答案）

- `*p` 是“p 指向的那个变量”，`&x` 是“x 的地址”；
- 先判断 `p == nil`，再解引用，否则会 panic；
- 结构体指针访问字段可以直接写 `s.Age`，Go 自动解引用，等价于 `(*s).Age`；
- 参数 `a [3]int` 本身就已经是副本，直接改它再返回就行；
- `for i := range a` 对数组指针 `a *[3]int` 同样可用，`a[i]` 也会自动解引用。

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestInc` | 通过指针加一；`nil` 不 panic |
| `TestSwap` | 交换指向的值 |
| `TestBirthdayThroughStructPointer` | 结构体指针修改字段 |
| `TestRenamedDoesNotTouchOriginal` | 结构体值传递，原件不变 |
| `TestArrayIsCopiedByValue` | 数组值传递，原件不变 |
| `TestScaleInPlaceThroughArrayPointer` | 传 `*[3]int` 改到原数组 |
| `TestPredictions` | 三个预测常量与实际输出一致 |
| `TestBridgeArrayCopiesButSliceShares` | 观察：数组赋值复制、切片赋值共享（不检查你的实现） |

## 做完以后

`TestBridgeArrayCopiesButSliceShares` 里同样是 `b := a`：数组复制了全部元素，切片却和原来共享数据。原因是切片头里藏着一个**指向底层数组的指针**。

下一题 `test02` 从这里开始：切片是“窗口”，不是数据本身。
