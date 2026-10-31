package interleave

import "errors"

// Op 是“线程”能执行的指令。每个线程只有一个私有寄存器 Reg。
type Op int

const (
	LOAD  Op = iota + 1 // Reg = Mem[Addr]
	ADDI                // Reg = Reg + Imm
	STORE               // Mem[Addr] = Reg
	INCR                // Mem[Addr] = Mem[Addr] + Imm，一步完成（“原子”加法）
)

// Instr 是一条指令。
type Instr struct {
	Op   Op
	Addr int
	Imm  int
}

// 出错时返回下面这些错误之一（测试用 errors.Is 判断）。
var (
	ErrBadThread  = errors.New("interleave: 调度表里的线程编号不存在")
	ErrThreadDone = errors.New("interleave: 调度了一个已经执行完的线程")
	ErrUnfinished = errors.New("interleave: 调度表结束时还有线程没执行完")
	ErrBadAddress = errors.New("interleave: 内存地址越界")
	ErrBadOp      = errors.New("interleave: 非法操作码")
)

// Increment 返回 x++ 在 CPU 上的样子：读 → 改 → 写，三条指令。
func Increment(addr int) []Instr {
	// TODO
	return nil
}

// AtomicIncrement 返回“一步完成”的加 1：只有一条 INCR。
func AtomicIncrement(addr int) []Instr {
	// TODO
	return nil
}

// Run 让 len(programs) 个线程共用一块内存，按 schedule 交错执行。
//
// schedule[k] = i 表示“第 k 步让线程 i 执行它的下一条指令”。
// 每个线程有自己的 PC（从 0 开始）和一个寄存器 Reg（初始为 0）。
// 返回执行完之后内存的副本；memory 本身不能被修改。
func Run(programs [][]Instr, memory []int, schedule []int) ([]int, error) {
	// TODO
	return nil, nil
}

// Schedules 列出所有合法的调度表。lengths[i] 是线程 i 的指令条数。
//
// 合法 = 线程 i 恰好出现 lengths[i] 次。结果按字典序从小到大排列，
// 例如 lengths = [2 1] 时返回 [0 0 1] [0 1 0] [1 0 0]。
func Schedules(lengths []int) [][]int {
	// TODO
	return nil
}

// Outcomes 把 programs 在所有合法调度下各跑一遍（每次都从同一份 memory 开始），
// 统计 Mem[addr] 的最终值各出现了多少次：最终值 -> 次数。
func Outcomes(programs [][]Instr, memory []int, addr int) (map[int]int, error) {
	// TODO
	return nil, nil
}
