package toycpu_test

import (
	"errors"
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test02/toycpu"
)

// newCPU 包一层：starter 的 New 返回 nil 时让单个测试失败，而不是 panic 拖垮整个包。
func newCPU(t *testing.T, program []toycpu.Instr, memSize int) *toycpu.CPU {
	t.Helper()
	c := toycpu.New(program, memSize)
	if c == nil {
		t.Fatal("New(...) 返回 nil；先把 New 写好")
	}
	return c
}

func TestNew(t *testing.T) {
	program := []toycpu.Instr{{Op: toycpu.HALT}}
	c := newCPU(t, program, 4)

	if c.PC != 0 || c.Halted {
		t.Fatalf("New 之后 PC = %d, Halted = %v, want 0, false", c.PC, c.Halted)
	}
	if !reflect.DeepEqual(c.Mem, []int{0, 0, 0, 0}) {
		t.Fatalf("New(..., 4).Mem = %v, want 4 个 0", c.Mem)
	}
	if c.Regs != [toycpu.NumRegs]int{} {
		t.Fatalf("New 之后 Regs = %v, want 全 0", c.Regs)
	}
	if len(c.Program) != 1 || c.Program[0].Op != toycpu.HALT {
		t.Fatalf("New 之后 Program = %v, want 装入传进来的程序", c.Program)
	}
}

func TestStepEachInstruction(t *testing.T) {
	cases := []struct {
		name    string
		in      toycpu.Instr
		setup   func(c *toycpu.CPU)
		check   func(c *toycpu.CPU) bool
		wantPC  int
		explain string
	}{
		{
			name:    "LOADI",
			in:      toycpu.Instr{Op: toycpu.LOADI, A: 3, Imm: 42},
			check:   func(c *toycpu.CPU) bool { return c.Regs[3] == 42 },
			wantPC:  1,
			explain: "R3 应该等于 42",
		},
		{
			name:    "LOAD",
			in:      toycpu.Instr{Op: toycpu.LOAD, A: 1, Imm: 2},
			setup:   func(c *toycpu.CPU) { c.Mem[2] = 7 },
			check:   func(c *toycpu.CPU) bool { return c.Regs[1] == 7 },
			wantPC:  1,
			explain: "R1 应该等于 Mem[2] = 7",
		},
		{
			name:    "STORE",
			in:      toycpu.Instr{Op: toycpu.STORE, A: 4, Imm: 3},
			setup:   func(c *toycpu.CPU) { c.Regs[4] = 9 },
			check:   func(c *toycpu.CPU) bool { return c.Mem[3] == 9 },
			wantPC:  1,
			explain: "Mem[3] 应该等于 R4 = 9",
		},
		{
			name:    "ADD",
			in:      toycpu.Instr{Op: toycpu.ADD, A: 0, B: 1, C: 2},
			setup:   func(c *toycpu.CPU) { c.Regs[1], c.Regs[2] = 5, -2 },
			check:   func(c *toycpu.CPU) bool { return c.Regs[0] == 3 },
			wantPC:  1,
			explain: "R0 应该等于 R1 + R2 = 3",
		},
		{
			name:    "ADDI",
			in:      toycpu.Instr{Op: toycpu.ADDI, A: 5, B: 5, Imm: -1},
			setup:   func(c *toycpu.CPU) { c.Regs[5] = 10 },
			check:   func(c *toycpu.CPU) bool { return c.Regs[5] == 9 },
			wantPC:  1,
			explain: "R5 应该等于 R5 + (-1) = 9",
		},
		{
			name:    "JNZ 跳转",
			in:      toycpu.Instr{Op: toycpu.JNZ, A: 6, Imm: 0},
			setup:   func(c *toycpu.CPU) { c.Regs[6] = 1 },
			check:   func(c *toycpu.CPU) bool { return true },
			wantPC:  0,
			explain: "R6 != 0 时 PC 应该跳到 Imm",
		},
		{
			name:    "JNZ 不跳",
			in:      toycpu.Instr{Op: toycpu.JNZ, A: 6, Imm: 0},
			check:   func(c *toycpu.CPU) bool { return true },
			wantPC:  1,
			explain: "R6 == 0 时 PC 应该 + 1",
		},
	}

	for _, tc := range cases {
		program := []toycpu.Instr{tc.in, {Op: toycpu.HALT}}
		c := newCPU(t, program, 4)
		if tc.setup != nil {
			tc.setup(c)
		}
		if err := c.Step(); err != nil {
			t.Errorf("%s: Step() error = %v, want nil", tc.name, err)
			continue
		}
		if !tc.check(c) {
			t.Errorf("%s: 执行后 Regs = %v, Mem = %v；%s", tc.name, c.Regs, c.Mem, tc.explain)
		}
		if c.PC != tc.wantPC {
			t.Errorf("%s: 执行后 PC = %d, want %d", tc.name, c.PC, tc.wantPC)
		}
	}
}

func TestHalt(t *testing.T) {
	c := newCPU(t, []toycpu.Instr{{Op: toycpu.HALT}, {Op: toycpu.LOADI, A: 1, Imm: 1}}, 1)

	if err := c.Step(); err != nil {
		t.Fatalf("HALT: Step() error = %v", err)
	}
	if !c.Halted || c.PC != 0 {
		t.Fatalf("HALT 之后 Halted = %v, PC = %d, want true, 0（停机时 PC 不再前进）", c.Halted, c.PC)
	}
	if err := c.Step(); err != nil || c.Regs[1] != 0 || c.PC != 0 {
		t.Fatalf("停机后再 Step：error = %v, R1 = %d, PC = %d；停机的 CPU 什么都不做", err, c.Regs[1], c.PC)
	}
}

// sumProgram 计算 1 + 2 + ... + Mem[0]，结果写到 Mem[1]。
var sumProgram = []toycpu.Instr{
	{Op: toycpu.LOAD, A: 1, Imm: 0},        // 0: R1 = n
	{Op: toycpu.LOADI, A: 2, Imm: 0},       // 1: R2 = 0
	{Op: toycpu.ADD, A: 2, B: 2, C: 1},     // 2: R2 += R1
	{Op: toycpu.ADDI, A: 1, B: 1, Imm: -1}, // 3: R1 -= 1
	{Op: toycpu.JNZ, A: 1, Imm: 2},         // 4: if R1 != 0 goto 2
	{Op: toycpu.STORE, A: 2, Imm: 1},       // 5: Mem[1] = R2
	{Op: toycpu.HALT},                      // 6
}

func TestRunSumProgram(t *testing.T) {
	for _, n := range []int{1, 5, 100} {
		c := newCPU(t, sumProgram, 2)
		c.Mem[0] = n

		steps, err := c.Run(10000)
		if err != nil {
			t.Fatalf("n = %d: Run() error = %v", n, err)
		}
		if want := n * (n + 1) / 2; c.Mem[1] != want {
			t.Fatalf("n = %d: Mem[1] = %d, want %d", n, c.Mem[1], want)
		}
		if want := 2 + 3*n + 2; steps != want {
			t.Fatalf("n = %d: Run() 执行了 %d 条, want %d（HALT 也算一条）", n, steps, want)
		}
		if !c.Halted {
			t.Fatalf("n = %d: Run 结束后 Halted = false", n)
		}
	}
}

func TestRunStopsAtStepLimit(t *testing.T) {
	forever := []toycpu.Instr{
		{Op: toycpu.LOADI, A: 1, Imm: 1},
		{Op: toycpu.JNZ, A: 1, Imm: 1}, // 永远跳回自己
	}
	c := newCPU(t, forever, 0)

	steps, err := c.Run(100)
	if !errors.Is(err, toycpu.ErrStepLimit) {
		t.Fatalf("死循环程序 Run(100) error = %v, want ErrStepLimit", err)
	}
	if steps != 100 {
		t.Fatalf("Run(100) 执行了 %d 条, want 正好 100", steps)
	}

	// 恰好在第 maxSteps 条停机，不算超限。
	c = newCPU(t, []toycpu.Instr{{Op: toycpu.LOADI, A: 1, Imm: 1}, {Op: toycpu.HALT}}, 0)
	if steps, err := c.Run(2); err != nil || steps != 2 {
		t.Fatalf("第 2 条正好是 HALT：Run(2) = %d, %v, want 2, nil", steps, err)
	}
	// 已经停机的 CPU 再 Run 什么都不做。
	if steps, err := c.Run(10); err != nil || steps != 0 {
		t.Fatalf("已停机的 CPU Run(10) = %d, %v, want 0, nil", steps, err)
	}
}

func TestRunReportsErrors(t *testing.T) {
	cases := []struct {
		name    string
		program []toycpu.Instr
		want    error
		steps   int
	}{
		{"零值指令", []toycpu.Instr{{}}, toycpu.ErrBadOp, 0},
		{"未知操作码", []toycpu.Instr{{Op: 99}}, toycpu.ErrBadOp, 0},
		{"寄存器越界 A", []toycpu.Instr{{Op: toycpu.LOADI, A: toycpu.NumRegs, Imm: 1}}, toycpu.ErrBadRegister, 0},
		{"寄存器越界 C", []toycpu.Instr{{Op: toycpu.ADD, A: 0, B: 1, C: -1}}, toycpu.ErrBadRegister, 0},
		{"读地址越界", []toycpu.Instr{{Op: toycpu.LOAD, A: 0, Imm: 4}}, toycpu.ErrBadAddress, 0},
		{"写地址越界", []toycpu.Instr{{Op: toycpu.LOADI, A: 0, Imm: 1}, {Op: toycpu.STORE, A: 0, Imm: -1}}, toycpu.ErrBadAddress, 1},
		{"没有 HALT，跑出程序末尾", []toycpu.Instr{{Op: toycpu.LOADI, A: 0, Imm: 1}}, toycpu.ErrPCOutOfRange, 1},
		{"跳到负数地址", []toycpu.Instr{{Op: toycpu.LOADI, A: 0, Imm: 1}, {Op: toycpu.JNZ, A: 0, Imm: -3}}, toycpu.ErrPCOutOfRange, 2},
	}

	for _, tc := range cases {
		c := newCPU(t, tc.program, 4)
		steps, err := c.Run(100)
		if !errors.Is(err, tc.want) {
			t.Errorf("%s: Run() error = %v, want %v", tc.name, err, tc.want)
			continue
		}
		if steps != tc.steps {
			t.Errorf("%s: Run() 执行了 %d 条, want %d（出错的那一条不算）", tc.name, steps, tc.steps)
		}
	}
}

func TestErrorLeavesStateUnchanged(t *testing.T) {
	c := newCPU(t, []toycpu.Instr{{Op: toycpu.STORE, A: 1, Imm: 10}}, 2)
	c.Regs[1] = 5
	before := append([]int(nil), c.Mem...)

	if err := c.Step(); !errors.Is(err, toycpu.ErrBadAddress) {
		t.Fatalf("Step() error = %v, want ErrBadAddress", err)
	}
	if c.PC != 0 || !reflect.DeepEqual(c.Mem, before) || c.Halted {
		t.Fatalf("出错后 PC = %d, Mem = %v, Halted = %v；出错的指令不能改动任何状态", c.PC, c.Mem, c.Halted)
	}
}

// TestBridgeIncrementIsThreeInstructions 把 Go 里的一句 x++ 放到这台 CPU 上：
// 它要三条指令（读 → 改 → 写），中间任何一步之后都可能被打断。下一题 test03 就从这里开始。
func TestBridgeIncrementIsThreeInstructions(t *testing.T) {
	increment := []toycpu.Instr{
		{Op: toycpu.LOAD, A: 1, Imm: 0},       // 读
		{Op: toycpu.ADDI, A: 1, B: 1, Imm: 1}, // 改
		{Op: toycpu.STORE, A: 1, Imm: 0},      // 写
		{Op: toycpu.HALT},
	}
	c := newCPU(t, increment, 1)
	c.Mem[0] = 41

	for i := 0; i < 3; i++ {
		if err := c.Step(); err != nil {
			t.Fatalf("第 %d 步 error = %v", i+1, err)
		}
		t.Logf("第 %d 步之后：R1 = %d，Mem[0] = %d", i+1, c.Regs[1], c.Mem[0])
	}
	if c.Mem[0] != 42 {
		t.Fatalf("x++ 之后 Mem[0] = %d, want 42", c.Mem[0])
	}
	t.Log("第 1、2 步之后，新值只在寄存器里，内存里还是旧值。如果此时另一个执行流也来读 Mem[0]……")
}
