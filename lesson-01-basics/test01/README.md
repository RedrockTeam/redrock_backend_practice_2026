# test01：把 C 的 Hello World 搬到 Go

对应课件：第一课 · 第二部分 1〜2.3（Hello World、分号与大括号、打印与大小写）。

## 任务

编辑 `hellogo/main.go`，让程序运行后输出**两行**：

```text
hello, go
hello, 红岩网校
```

其中：

- 第一行必须用 `fmt.Println` 打印；
- 第二行必须用 `fmt.Printf` 打印。

末尾有换行，前后没有多余空格、没有多余空行。

## 你需要用到的知识

课件里的 C 版本是这样的：

```c
#include <stdio.h>

int main(void) {
    printf("hello world!\n");
    return 0;
}
```

按课件第二部分第 3 节给出的迁移清单改写：

1. 第一行写 `package main`，声明这是一个可执行程序；
2. 需要打印就 `import "fmt"`；
3. `int main(void)` 改成 `func main()`，去掉 `return 0`；
4. `printf` 改成 `fmt.Println` 或 `fmt.Printf`；
5. 行尾不写分号，`{` 必须和 `func main()` 在同一行。

## 为什么要求两个都用一遍

因为它们有一个很容易踩的区别：

| | 换行 | 写法 |
|---|---|---|
| `fmt.Println("hello, go")` | **自动换行** | 直接给内容 |
| `fmt.Printf("hello, go\n")` | **不换行** | `\n` 要自己写 |

光看输出看不出你用了哪个，所以测试会直接读你的源码，确认两个都出现过。忘了写 `\n`，两行就会粘成一行。

另外注意大小写：包名 `fmt` 是小写，函数 `Println`、`Printf` 首字母大写。Go 就是用首字母大小写决定一个名字能不能被别的包使用的——第五部分会把这条规则讲透。

## 自己验证

```bash
go run ./lesson-01-basics/test01/hellogo
go test ./lesson-01-basics/test01/...
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestHelloWorld` | `main()` 的标准输出逐字等于那两行 |
| `TestUsesBothPrintlnAndPrintf` | 源码里 `fmt.Println` 和 `fmt.Printf` 都用到了 |

测试的做法是临时替换 `os.Stdout`、调用 `main()`、再把捕获到的内容和期望值比较。你现在不需要看懂这段测试代码，第 test07 题会带你自己写测试。
