package pointer

import "fmt"

// 本文件是课件里的三段示例，已经写好，不要修改。
// 先读代码、在 predict.go 里写下你的预测，再运行测试核对。

func setFirst(a [3]int) {
	a[0] = 100
}

// DemoArrayByValue 把数组传给函数并在函数里修改它，返回之后打印原数组。
func DemoArrayByValue() string {
	x := [3]int{1, 2, 3}
	setFirst(x)
	return fmt.Sprint(x)
}

func inc(p *int) {
	*p = *p + 1
}

// DemoPointerArgument 把变量的地址传给函数，返回之后打印变量。
func DemoPointerArgument() string {
	x := 10
	inc(&x)
	return fmt.Sprint(x)
}

// DemoStructPointer 通过结构体指针修改字段，然后打印原结构体的字段。
func DemoStructPointer() string {
	s := Student{Name: "Tom", Age: 18}
	p := &s
	p.Age = 20
	copied := s
	copied.Age = 30
	return fmt.Sprint(s.Age)
}
