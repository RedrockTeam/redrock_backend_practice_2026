package interleave_test

import (
	"errors"
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test03/interleave"
)

func twoIncrements() [][]interleave.Instr {
	return [][]interleave.Instr{interleave.Increment(0), interleave.Increment(0)}
}

func TestIncrementPrograms(t *testing.T) {
	want := []interleave.Instr{
		{Op: interleave.LOAD, Addr: 3},
		{Op: interleave.ADDI, Imm: 1},
		{Op: interleave.STORE, Addr: 3},
	}
	if got := interleave.Increment(3); !reflect.DeepEqual(got, want) {
		t.Fatalf("Increment(3) = %+v, want %+v", got, want)
	}

	wantAtomic := []interleave.Instr{{Op: interleave.INCR, Addr: 3, Imm: 1}}
	if got := interleave.AtomicIncrement(3); !reflect.DeepEqual(got, wantAtomic) {
		t.Fatalf("AtomicIncrement(3) = %+v, want %+v", got, wantAtomic)
	}
}

func runFinal(t *testing.T, programs [][]interleave.Instr, memory []int, schedule []int) int {
	t.Helper()
	got, err := interleave.Run(programs, memory, schedule)
	if err != nil {
		t.Fatalf("Run(schedule = %v) error = %v", schedule, err)
	}
	if len(got) != len(memory) {
		t.Fatalf("Run(schedule = %v) 返回的内存长度 = %d, want %d", schedule, len(got), len(memory))
	}
	return got[0]
}

func TestRunSerialSchedules(t *testing.T) {
	for _, schedule := range [][]int{{0, 0, 0, 1, 1, 1}, {1, 1, 1, 0, 0, 0}} {
		if got := runFinal(t, twoIncrements(), []int{0}, schedule); got != 2 {
			t.Errorf("串行调度 %v：Mem[0] = %d, want 2", schedule, got)
		}
	}
	if got := runFinal(t, twoIncrements(), []int{40}, []int{0, 0, 0, 1, 1, 1}); got != 42 {
		t.Errorf("从 40 开始串行加两次：Mem[0] = %d, want 42", got)
	}
}

func TestRunInterleavedScheduleLosesAnUpdate(t *testing.T) {
	for _, schedule := range [][]int{{0, 1, 0, 1, 0, 1}, {0, 0, 1, 1, 1, 0}, {1, 0, 0, 0, 1, 1}} {
		if got := runFinal(t, twoIncrements(), []int{0}, schedule); got != 1 {
			t.Errorf("交错调度 %v：Mem[0] = %d, want 1；两个线程都读到了旧值 0，后写的覆盖了先写的", schedule, got)
		}
	}
}

func TestRunKeepsThreadsPrivateAndInputIntact(t *testing.T) {
	// 线程 0 把 Mem[0] 搬到 Mem[1]，线程 1 给 Mem[2] 加 5。寄存器互不干扰。
	programs := [][]interleave.Instr{
		{{Op: interleave.LOAD, Addr: 0}, {Op: interleave.STORE, Addr: 1}},
		{{Op: interleave.LOAD, Addr: 2}, {Op: interleave.ADDI, Imm: 5}, {Op: interleave.STORE, Addr: 2}},
	}
	memory := []int{7, 0, 10}
	got, err := interleave.Run(programs, memory, []int{0, 1, 1, 0, 1})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !reflect.DeepEqual(got, []int{7, 7, 15}) {
		t.Fatalf("Run() = %v, want [7 7 15]；每个线程要有自己的寄存器", got)
	}
	if !reflect.DeepEqual(memory, []int{7, 0, 10}) {
		t.Fatalf("Run 修改了传进来的 memory：%v；要先复制一份再执行", memory)
	}
}

func TestRunErrors(t *testing.T) {
	cases := []struct {
		name     string
		programs [][]interleave.Instr
		schedule []int
		want     error
	}{
		{"线程编号越界", twoIncrements(), []int{0, 0, 0, 2}, interleave.ErrBadThread},
		{"线程编号为负", twoIncrements(), []int{-1}, interleave.ErrBadThread},
		{"调度已结束的线程", twoIncrements(), []int{0, 0, 0, 0, 1, 1, 1}, interleave.ErrThreadDone},
		{"调度表太短", twoIncrements(), []int{0, 0, 0, 1, 1}, interleave.ErrUnfinished},
		{"地址越界", [][]interleave.Instr{interleave.Increment(5)}, []int{0, 0, 0}, interleave.ErrBadAddress},
		{"非法操作码", [][]interleave.Instr{{{}}}, []int{0}, interleave.ErrBadOp},
	}
	for _, tc := range cases {
		if _, err := interleave.Run(tc.programs, []int{0}, tc.schedule); !errors.Is(err, tc.want) {
			t.Errorf("%s: Run() error = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestSchedules(t *testing.T) {
	small := interleave.Schedules([]int{2, 1})
	want := [][]int{{0, 0, 1}, {0, 1, 0}, {1, 0, 0}}
	if !reflect.DeepEqual(small, want) {
		t.Fatalf("Schedules([2 1]) = %v, want %v（字典序）", small, want)
	}

	if got := interleave.Schedules([]int{2}); !reflect.DeepEqual(got, [][]int{{0, 0}}) {
		t.Fatalf("Schedules([2]) = %v, want [[0 0]]", got)
	}

	// 没有指令要执行时，只有一种调度：什么都不做。
	if got := interleave.Schedules([]int{0, 0}); len(got) != 1 || len(got[0]) != 0 {
		t.Fatalf("Schedules([0 0]) = %v, want 恰好一个空调度", got)
	}

	all := interleave.Schedules([]int{3, 3})
	if len(all) != 20 {
		t.Fatalf("Schedules([3 3]) 共 %d 种, want 20（6 个位置里挑 3 个给线程 0）", len(all))
	}
	if !reflect.DeepEqual(all[0], []int{0, 0, 0, 1, 1, 1}) || !reflect.DeepEqual(all[19], []int{1, 1, 1, 0, 0, 0}) {
		t.Fatalf("Schedules([3 3]) 首尾 = %v, %v；要按字典序", all[0], all[19])
	}
	seen := map[string]bool{}
	for _, schedule := range all {
		key := ""
		count := [2]int{}
		for _, thread := range schedule {
			key += string(rune('0' + thread))
			count[thread]++
		}
		if count != [2]int{3, 3} {
			t.Fatalf("调度 %v 不合法：每个线程要恰好出现 3 次", schedule)
		}
		if seen[key] {
			t.Fatalf("调度 %v 重复出现", schedule)
		}
		seen[key] = true
	}

	if got := len(interleave.Schedules([]int{3, 3, 3})); got != 1680 {
		t.Fatalf("Schedules([3 3 3]) 共 %d 种, want 1680", got)
	}
}

func TestOutcomesOfTwoIncrements(t *testing.T) {
	got, err := interleave.Outcomes(twoIncrements(), []int{0}, 0)
	if err != nil {
		t.Fatalf("Outcomes() error = %v", err)
	}
	want := map[int]int{1: 18, 2: 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("两个线程各 x++ 一次，所有调度的结果分布 = %v, want %v", got, want)
	}
	t.Logf("20 种调度里，只有 %d 种得到 2，%d 种丢了一次更新。", got[2], got[1])
}

func TestOutcomesOfThreeIncrements(t *testing.T) {
	programs := [][]interleave.Instr{interleave.Increment(0), interleave.Increment(0), interleave.Increment(0)}
	got, err := interleave.Outcomes(programs, []int{0}, 0)
	if err != nil {
		t.Fatalf("Outcomes() error = %v", err)
	}
	total := 0
	for _, n := range got {
		total += n
	}
	if total != 1680 {
		t.Fatalf("三个线程的调度共统计了 %d 次, want 1680", total)
	}
	want := map[int]int{1: 1404, 2: 270, 3: 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("三个线程各 x++ 一次的结果分布 = %v, want %v", got, want)
	}
	t.Logf("1680 种调度里只有 %d 种是对的。", got[3])
}

// TestBridgeAtomicIncrementNeverLoses：读-改-写合成一步之后，再怎么交错也不会丢。
// 第四课 test02 的 Mutex / atomic 做的就是这件事：让“读-改-写”在别人眼里像一步。
func TestBridgeAtomicIncrementNeverLoses(t *testing.T) {
	programs := [][]interleave.Instr{
		interleave.AtomicIncrement(0),
		interleave.AtomicIncrement(0),
		interleave.AtomicIncrement(0),
	}
	got, err := interleave.Outcomes(programs, []int{0}, 0)
	if err != nil {
		t.Fatalf("Outcomes() error = %v", err)
	}
	if !reflect.DeepEqual(got, map[int]int{3: 6}) {
		t.Fatalf("三个原子加 1 的结果分布 = %v, want map[3:6]", got)
	}
	t.Log("原子加法：6 种调度全部得到 3。")
	t.Log("真实 CPU 上 x++ 是三条指令，调度器可以在任意两条之间切换线程——这就是第四课要解决的数据竞争。")
}
