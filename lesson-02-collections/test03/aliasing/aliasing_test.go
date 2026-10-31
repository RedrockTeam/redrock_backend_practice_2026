package aliasing_test

import (
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test03/aliasing"
)

func TestHeadDoesNotLeakCapacity(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	head := aliasing.Head(values, 2)
	if !reflect.DeepEqual(head, []int{1, 2}) {
		t.Fatalf("Head(values, 2) = %v, want [1 2]", head)
	}

	head = append(head, 99)
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("对 Head 的结果 append 之后 values = %v, want [1 2 3 4 5]；返回的窗口 cap 太大了", values)
	}
	if !reflect.DeepEqual(head, []int{1, 2, 99}) {
		t.Fatalf("append 之后 head = %v, want [1 2 99]", head)
	}

	if got := aliasing.Head(values, 10); !reflect.DeepEqual(got, values) {
		t.Fatalf("Head(values, 10) = %v, want 全部 %v", got, values)
	}
}

func TestBranchPathsAreIndependent(t *testing.T) {
	base := make([]string, 2, 8) // 故意留出富余的 cap
	base[0], base[1] = "/", "home"

	left, right := aliasing.Branch(base, "alice", "bob")
	if !reflect.DeepEqual(left, []string{"/", "home", "alice"}) {
		t.Fatalf("left = %v, want [/ home alice]；两次 append 是不是写进了同一格？", left)
	}
	if !reflect.DeepEqual(right, []string{"/", "home", "bob"}) {
		t.Fatalf("right = %v, want [/ home bob]", right)
	}

	left[0] = "changed"
	if base[0] != "/" || right[0] != "/" {
		t.Fatalf("修改 left[0] 之后 base[0] = %q, right[0] = %q, want 都是 \"/\"；三条路径要互相独立", base[0], right[0])
	}
}

func TestRemoveIfHandlesAdjacentMatches(t *testing.T) {
	values := []int{2, 4, 5, 6, 8, 9}
	isEven := func(v int) bool { return v%2 == 0 }

	got := aliasing.RemoveIf(values, isEven)
	if !reflect.DeepEqual(got, []int{5, 9}) {
		t.Fatalf("RemoveIf(%v, isEven) = %v, want [5 9]；相邻的两个偶数，第二个是不是被跳过了？", []int{2, 4, 5, 6, 8, 9}, got)
	}

	if got := aliasing.RemoveIf(nil, isEven); len(got) != 0 {
		t.Fatalf("RemoveIf(nil, ...) = %v, want 空", got)
	}
	if got := aliasing.RemoveIf([]int{2, 4}, isEven); len(got) != 0 {
		t.Fatalf("RemoveIf([2 4], isEven) = %v, want 空", got)
	}
}

func TestRemoveIfDoesNotMutateInput(t *testing.T) {
	values := []int{1, 2, 3, 4}
	aliasing.RemoveIf(values, func(v int) bool { return v == 2 })
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4}) {
		t.Fatalf("RemoveIf 之后调用方的 values = %v, want [1 2 3 4]；不能原地挪动别人的数据", values)
	}
}

func TestPointersPointIntoTheSlice(t *testing.T) {
	students := []aliasing.Student{{"Lin", 80}, {"Wang", 90}, {"Zhao", 70}}
	ptrs := aliasing.Pointers(students)
	if len(ptrs) != 3 {
		t.Fatalf("len(Pointers(...)) = %d, want 3", len(ptrs))
	}
	for i, p := range ptrs {
		if p != &students[i] {
			t.Fatalf("Pointers(...)[%d] 没有指向 students[%d]；range 的 value 是副本", i, i)
		}
	}
	ptrs[1].Score = 100
	if students[1].Score != 100 {
		t.Fatalf("通过指针改 Score 之后 students[1].Score = %d, want 100", students[1].Score)
	}
}

func TestChunkPiecesDoNotOverlap(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7}
	chunks := aliasing.Chunk(values, 3)
	want := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	if !reflect.DeepEqual(chunks, want) {
		t.Fatalf("Chunk(values, 3) = %v, want %v", chunks, want)
	}

	chunks[0] = append(chunks[0], 100)
	if !reflect.DeepEqual(chunks[1], []int{4, 5, 6}) {
		t.Fatalf("对第 0 段 append 之后第 1 段 = %v, want [4 5 6]；每一段的 cap 要收紧", chunks[1])
	}
	if !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5, 6, 7}) {
		t.Fatalf("对第 0 段 append 之后 values = %v, want 不变", values)
	}

	if got := aliasing.Chunk(values, 0); got != nil {
		t.Fatalf("Chunk(values, 0) = %v, want nil", got)
	}
}

// TestBridgeAppendWritesInPlaceWhenCapAllows 不检查你的实现，复现课件里的例子：
// cap 足够时 append 原地写入，两个从同一个切片 append 出来的结果互相覆盖。
func TestBridgeAppendWritesInPlaceWhenCapAllows(t *testing.T) {
	s := make([]int, 2, 5)
	s[0], s[1] = 1, 2
	first := append(s, 3)
	second := append(s, 9)

	if first[2] != 9 || second[2] != 9 {
		t.Fatalf("first = %v, second = %v；课件里说它们共享底层数组", first, second)
	}
	t.Logf("first = %v, second = %v：两次 append 写的是同一格。", first, second)
	t.Log("完整切片表达式 s[low:high:max] 能把 cap 收到 max-low，下一次 append 就会被迫分配新数组。")
}
