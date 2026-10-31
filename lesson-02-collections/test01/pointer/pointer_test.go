package pointer_test

import (
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test01/pointer"
)

func TestInc(t *testing.T) {
	x := 41
	pointer.Inc(&x)
	if x != 42 {
		t.Fatalf("Inc(&x) 之后 x = %d, want 42；要通过 *p 修改 p 指向的变量", x)
	}
	pointer.Inc(&x)
	if x != 43 {
		t.Fatalf("第二次 Inc(&x) 之后 x = %d, want 43", x)
	}
	pointer.Inc(nil) // 不应该 panic
}

func TestSwap(t *testing.T) {
	a, b := 1, 2
	pointer.Swap(&a, &b)
	if a != 2 || b != 1 {
		t.Fatalf("Swap(&a, &b) 之后 a, b = %d, %d, want 2, 1；交换的是指针指向的值，不是指针本身", a, b)
	}
}

func TestBirthdayThroughStructPointer(t *testing.T) {
	s := pointer.Student{Name: "Tom", Age: 18}
	pointer.Birthday(&s)
	if s.Age != 19 {
		t.Fatalf("Birthday(&s) 之后 s.Age = %d, want 19", s.Age)
	}
}

func TestRenamedDoesNotTouchOriginal(t *testing.T) {
	s := pointer.Student{Name: "Tom", Age: 18}
	got := pointer.Renamed(s, "Jerry")
	if got.Name != "Jerry" || got.Age != 18 {
		t.Fatalf("Renamed(s, \"Jerry\") = %+v, want {Name:Jerry Age:18}", got)
	}
	if s.Name != "Tom" {
		t.Fatalf("Renamed 之后原来的 s.Name = %q, want \"Tom\"；结构体是值传递", s.Name)
	}
}

func TestArrayIsCopiedByValue(t *testing.T) {
	a := [3]int{1, 2, 3}
	got := pointer.Scaled(a, 10)
	if got != [3]int{10, 20, 30} {
		t.Fatalf("Scaled(%v, 10) = %v, want [10 20 30]", a, got)
	}
	if a != [3]int{1, 2, 3} {
		t.Fatalf("Scaled 之后原数组 = %v, want [1 2 3]", a)
	}
}

func TestScaleInPlaceThroughArrayPointer(t *testing.T) {
	a := [3]int{1, -2, 3}
	pointer.ScaleInPlace(&a, 3)
	if a != [3]int{3, -6, 9} {
		t.Fatalf("ScaleInPlace(&a, 3) 之后 a = %v, want [3 -6 9]；传 *[3]int 才能改到调用方的数组", a)
	}
}

func TestPredictions(t *testing.T) {
	cases := []struct {
		name    string
		predict string
		actual  string
	}{
		{"PredictArrayByValue", pointer.PredictArrayByValue, pointer.DemoArrayByValue()},
		{"PredictPointerArgument", pointer.PredictPointerArgument, pointer.DemoPointerArgument()},
		{"PredictStructPointer", pointer.PredictStructPointer, pointer.DemoStructPointer()},
	}
	for _, c := range cases {
		if c.predict == "TODO" {
			t.Errorf("%s 还是 \"TODO\"：先读 demo.go 写下你的预测", c.name)
			continue
		}
		if c.predict != c.actual {
			t.Errorf("%s = %q，实际输出是 %q；想一想这里传的是值还是地址", c.name, c.predict, c.actual)
		}
	}
}

// TestBridgeArrayCopiesButSliceShares 不检查你的实现。
// 它把数组和切片放在一起：同样写 b := a，结果完全不同。下一题 test02 从这里开始。
func TestBridgeArrayCopiesButSliceShares(t *testing.T) {
	array := [3]int{1, 2, 3}
	arrayCopy := array
	arrayCopy[0] = 100

	slice := []int{1, 2, 3}
	sliceCopy := slice
	sliceCopy[0] = 100

	if array[0] != 1 {
		t.Fatalf("array[0] = %d, want 1；数组赋值复制全部元素", array[0])
	}
	if slice[0] != 100 {
		t.Fatalf("slice[0] = %d, want 100；切片赋值只复制切片头", slice[0])
	}
	t.Log("数组 b := a 复制了 3 个元素；切片 b := a 只复制了 (指针, len, cap)。")
	t.Log("切片里藏着一个指向底层数组的指针，所以它表现得像“窗口”。test02 会把这件事做成题。")
}
