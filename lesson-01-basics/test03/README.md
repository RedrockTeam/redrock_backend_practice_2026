# test03：把一整个 C 程序改写成 Go

对应课件：第一课 · 第二部分 3（「把 C 代码改成 Go 代码」的十步清单）。

前两题是一句一句地学语法，这一题是把它们连起来用一次：给你一个完整的 C 程序，要求逐条按清单改写成 Go，输出一模一样。

## 任务

题面是 `reference.c`（它不参与编译，只是给你看的）：

```c
#include <stdio.h>

int square(int n) {
    return n * n;
}

int main(void) {
    int total = 0;

    for (int i = 1; i <= 5; i++) {
        if (i % 2 == 0) {
            printf("%d even %d\n", i, square(i));
        } else {
            printf("%d odd %d\n", i, square(i));
        }
        total = total + square(i);
    }

    printf("total %d\n", total);
    return 0;
}
```

把它改写进 `translate/translate.go`，让程序输出：

```text
1 odd 1
2 even 4
3 odd 9
4 even 16
5 odd 25
total 55
```

**结构要保持一致**：`square` 仍然是一个独立函数并且被 `main` 调用，循环仍然是循环，`if/else` 仍然是 `if/else`。把六行结果直接写成六条 `Println` 不算过。

## 课件里的十步清单

1. 新建 `.go` 文件，第一行写 `package main`；
2. 需要打印就 `import "fmt"`；
3. `int main(void)` 改成 `func main()`；
4. 删掉 `return 0`；
5. `printf(...)` 改成 `fmt.Printf(...)` 或 `fmt.Println(...)`；
6. 删掉所有行尾分号；
7. 确保 `{` 跟在 `func`、`if`、`for` 同一行；
8. 变量声明改成 `var 名字 类型`，函数内可用 `名字 := 值`；
9. 函数定义改成 `func 函数名(参数) 返回类型`；
10. 用 `go run` 运行。

## 四个容易卡住的地方

**分号。** C 的每行末尾都有 `;`，Go 全部删掉——但 `for i := 1; i <= 5; i++` 里那两个分号是语法要求的，必须留着。测试会逐个 token 检查：留在行尾的手写分号会被点名，`for` 头部里的不算。

```go
total := 0        // 对
total := 0;       // 会被判不通过
```

**大括号。** `{` 换行写会让编译器在上一行自动插入分号，整段就散了：

```go
for i := 1; i <= 5; i++
{                          // 编译错误
```

**循环变量的作用域。** C 写 `for (int i = 1; ...)`，Go 写 `for i := 1; ...`——没有 `int`，用 `:=` 推断类型。

**格式串不变。** `printf("%d even %d\n", ...)` 里的 `%d` 和 `\n` 在 `fmt.Printf` 里含义完全一样，照抄即可。这是 C 程序员迁移到 Go 时少数不用改的地方。

## 自己验证

```bash
go run ./lesson-01-basics/test03/translate
go test ./lesson-01-basics/test03/...
```

## 评分点

| 测试 | 检查内容 |
|---|---|
| `TestTranslationOutput` | 输出和 `reference.c` 逐字一致 |
| `TestTranslationKeepsSquare` | `square` 是一个单参数单返回值的函数，并且 `main` 调用了它 |
| `TestTranslationKeepsControlFlow` | `main` 里有循环，也有带 `else` 的 `if` |
| `TestTranslationDropsSemicolons` | 除 `for` 头部外没有手写分号 |

后三个测试是直接读你的源码判的（Go 标准库自带 `go/parser`，能把 Go 代码解析成语法树）。它们的作用是保证你是在**翻译**，而不是在凑输出。
