package toycpu

import "errors"

// NumRegs 是寄存器的个数：R0..R7。
const NumRegs = 8

// Op 是操作码。故意从 1 开始：Instr{} 的零值不是一条合法指令。
type Op int

const (
	HALT  Op = iota + 1 // 停机
	LOADI               // R[A] = Imm
	LOAD                // R[A] = Mem[Imm]
	STORE               // Mem[Imm] = R[A]
	ADD                 // R[A] = R[B] + R[C]
	ADDI                // R[A] = R[B] + Imm
	JNZ                 // 如果 R[A] != 0，PC = Imm；否则 PC + 1
)

// Instr 是一条已经译好码的指令。不同操作码只用到其中一部分字段。
type Instr struct {
	Op      Op
	A, B, C int // 寄存器编号
	Imm     int // 立即数：常数、内存地址或跳转目标
}

// CPU 是一台极简计算机：程序和数据分开放，PC 是“下一条指令的下标”。
type CPU struct {
	Regs    [NumRegs]int
	Mem     []int
	PC      int
	Program []Instr
	Halted  bool
}

// 出错时返回下面这些错误之一（测试用 errors.Is 判断）。
var (
	ErrBadOp        = errors.New("toycpu: 非法操作码")
	ErrBadRegister  = errors.New("toycpu: 寄存器编号越界")
	ErrBadAddress   = errors.New("toycpu: 内存地址越界")
	ErrPCOutOfRange = errors.New("toycpu: PC 越过了程序边界")
	ErrStepLimit    = errors.New("toycpu: 超过最大步数")
)

// New 返回一台新 CPU：装好程序，分配 memSize 个内存单元（全部为 0），PC = 0。
func New(program []Instr, memSize int) *CPU {
	// TODO
	return nil
}

// Step 执行恰好一条指令：取指 → 译码 → 执行 → 写回 → 更新 PC。
// 具体规则见 README 的“规格”一节。
func (c *CPU) Step() error {
	// TODO
	return nil
}

// Run 反复调用 Step，直到停机、出错或执行满 maxSteps 条。
// 返回成功执行的指令条数（HALT 也算一条）。
func (c *CPU) Run(maxSteps int) (int, error) {
	// TODO
	return 0, nil
}
