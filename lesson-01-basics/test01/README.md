# test01：把 C 的 Hello World 搬到 Go

对应课件：第一课 · 第二部分 1〜3（Hello World、分号与大括号、打印与大小写）。

## 任务

编辑 `hellogo/main.go`，让程序运行后**只输出一行**：

```text
hello, go
```

注意末尾有换行，前后没有多余空格、没有多余空行。

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

`fmt.Println` 会自动换行，`fmt.Printf` 不会，需要自己写 `\n`。两种写法都能通过测试。

## 自己验证

```bash
go run ./lesson-01-basics/test01/hellogo
go test ./lesson-01-basics/test01/...
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestHelloWorld` | `main()` 的标准输出等于 `"hello, go\n"` |

测试的做法是临时替换 `os.Stdout`、调用 `main()`、再把捕获到的内容和期望值比较。你现在不需要看懂这段测试代码，第 test05 题会带你自己写测试。
