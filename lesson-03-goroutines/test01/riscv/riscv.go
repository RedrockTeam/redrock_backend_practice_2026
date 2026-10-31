package riscv

import "fmt"

// Instr 是一条 32 位 RISC-V R 型指令拆出来的六个字段。
//
//	31      25 24  20 19  15 14  12 11   7 6      0
//	+---------+------+------+------+------+--------+
//	| funct7  | rs2  | rs1  |funct3|  rd  | opcode |
//	+---------+------+------+------+------+--------+
//	   7 位     5 位   5 位   3 位   5 位    7 位
type Instr struct {
	Opcode uint32 // bit 0..6
	Rd     uint32 // bit 7..11
	Funct3 uint32 // bit 12..14
	Rs1    uint32 // bit 15..19
	Rs2    uint32 // bit 20..24
	Funct7 uint32 // bit 25..31
}

// LittleEndian 把内存里按小端序存放的 4 个字节拼成一个 32 位机器码。
// b[0] 是低地址、也是最低的 8 位。
//
//	内存：33 84 72 00  ->  机器码 0x00728433
func LittleEndian(b [4]byte) uint32 {
	// TODO: 把 b[0]..b[3] 分别左移 0、8、16、24 位，再用 | 拼起来。
	return 0
}

// Decode 从机器码里取出六个字段。
func Decode(word uint32) Instr {
	// TODO: 先右移到最低位，再用掩码 & 只留下需要的那几位。
	// 下面的 Opcode 已经给出示范：7 位的掩码是 0x7f。
	return Instr{
		Opcode: word & 0x7f,
	}
}

// Encode 是 Decode 的逆运算：把六个字段拼回一个机器码。
// 字段超出自己的位宽时，只保留低位（先 & 掩码再移位）。
func Encode(in Instr) uint32 {
	// TODO: 每个字段先 & 掩码，再左移到自己的位置，最后用 | 合起来。
	return 0
}

// Disasm 把一条指令翻译成汇编文本。
//
//	opcode 0x33, funct3 0, funct7 0x00  ->  "add x8, x5, x7"
//	opcode 0x33, funct3 0, funct7 0x20  ->  "sub x8, x5, x7"
//	其它任何组合                         ->  "unknown 0x00728433"（8 位十六进制机器码）
func Disasm(in Instr) string {
	// TODO: 用 switch / if 判断是哪条指令，再用 fmt.Sprintf 拼出文本。
	return fmt.Sprintf("")
}
