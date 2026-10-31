package sliceops_test

import (
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test02/sliceops"
)

func TestCloneIsIndependent(t *testing.T) {
	base := []int{1, 2, 3}
	clone := sliceops.Clone(base)

	if !reflect.DeepEqual(clone, []int{1, 2, 3}) {
		t.Fatalf("Clone(%v) = %v, want [1 2 3]", base, clone)
	}

	clone[0] = 100
	clone[1] = 200

	if base[0] != 1 || base[1] != 2 {
		t.Fatalf("修改 clone 后 base = %v；Clone 必须复制底层数组，不能只复制切片头", base)
	}
	if !reflect.DeepEqual(base, []int{1, 2, 3}) {
		t.Fatalf("Clone 之后 base = %v, want [1 2 3]", base)
	}

	if got := sliceops.Clone(nil); got != nil {
		t.Fatalf("Clone(nil) = %#v, want nil", got)
	}
}

func TestDoubleInPlace(t *testing.T) {
	values := []int{1, -2, 0, 5}
	sliceops.DoubleInPlace(values)

	want := []int{2, -4, 0, 10}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("DoubleInPlace 之后 values = %v, want %v；range 的 value 是副本，要用下标写回", values, want)
	}

	var nilValues []int
	sliceops.DoubleInPlace(nilValues) // 不应该 panic
}

func TestAppendOne(t *testing.T) {
	base := []int{1, 2}
	got := sliceops.AppendOne(base, 3)

	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("AppendOne(%v, 3) = %v, want [1 2 3]；append 的返回值不能丢", base, got)
	}
	if len(base) != 2 || !reflect.DeepEqual(base, []int{1, 2}) {
		t.Fatalf("AppendOne 之后 base = %v, want [1 2]；返回值是新的切片头", base)
	}

	if got := sliceops.AppendOne(nil, 7); !reflect.DeepEqual(got, []int{7}) {
		t.Fatalf("AppendOne(nil, 7) = %v, want [7]", got)
	}
}

func TestSum(t *testing.T) {
	cases := []struct {
		values []int
		want   int
	}{
		{values: []int{1, 2, 3, 4}, want: 10},
		{values: []int{-1, 1}, want: 0},
		{values: nil, want: 0},
		{values: []int{}, want: 0},
	}
	for _, c := range cases {
		if got := sliceops.Sum(c.values); got != c.want {
			t.Errorf("Sum(%v) = %d, want %d", c.values, got, c.want)
		}
	}
}

// TestBridgeSliceAssignmentIsAWindow 不检查你的实现，只把下一题要用的
// 「切片是窗口，不是数据」摆出来看。map 也是句柄，后面会再次遇到。
func TestBridgeSliceAssignmentIsAWindow(t *testing.T) {
	base := []int{1, 2, 3}
	alias := base
	alias[0] = 100

	if base[0] != 100 {
		t.Fatalf("base[0] = %d, want 100；切片赋值只复制了切片头", base[0])
	}

	t.Log("slice 赋值只复制 header（数据指针、len、cap），底层数组仍然共享。")
	t.Log("所以 Clone 必须 make + copy；map 也类似，变量只是句柄。")
}
