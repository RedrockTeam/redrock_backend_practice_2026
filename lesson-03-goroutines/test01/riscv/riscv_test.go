package riscv_test

import (
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test01/riscv"
)

// 这三组“内存字节 -> 机器码”来自课件第三课「cpu 工作简介」里的那张图，只用来练小端序。
// （图里两条 lw 的机器码是示意值，真正的 lw x5, 0(x6) 编码是 0x00032283，见 TestDecodeEveryField。）
var courseware = []struct {
	asm   string
	bytes [4]byte
	word  uint32
}{
	{asm: "图中第 1 行", bytes: [4]byte{0x83, 0x22, 0x31, 0x00}, word: 0x00312283},
	{asm: "图中第 2 行", bytes: [4]byte{0x83, 0x23, 0x43, 0x00}, word: 0x00432383},
	{asm: "图中第 3 行 add x8, x5, x7", bytes: [4]byte{0x33, 0x84, 0x72, 0x00}, word: 0x00728433},
}

func TestLittleEndian(t *testing.T) {
	for _, c := range courseware {
		if got := riscv.LittleEndian(c.bytes); got != c.word {
			t.Errorf("%s: LittleEndian(% x) = 0x%08x, want 0x%08x；b[0] 是最低 8 位", c.asm, c.bytes, got, c.word)
		}
	}
	if got := riscv.LittleEndian([4]byte{0x78, 0x56, 0x34, 0x12}); got != 0x12345678 {
		t.Errorf("LittleEndian(78 56 34 12) = 0x%08x, want 0x12345678", got)
	}
	if got := riscv.LittleEndian([4]byte{0xff, 0xff, 0xff, 0xff}); got != 0xffffffff {
		t.Errorf("LittleEndian(ff ff ff ff) = 0x%08x, want 0xffffffff；注意先转成 uint32 再移位", got)
	}
}

func TestDecodeCoursewareAdd(t *testing.T) {
	got := riscv.Decode(0x00728433)
	want := riscv.Instr{Opcode: 0x33, Rd: 8, Funct3: 0, Rs1: 5, Rs2: 7, Funct7: 0}
	if got != want {
		t.Fatalf("Decode(0x00728433) = %+v, want %+v（课件里 add x8, x5, x7 的拆解）", got, want)
	}
}

func TestDecodeEveryField(t *testing.T) {
	// 每个字段取各自位宽的最大值，相邻字段不能互相“串位”。
	cases := []struct {
		word uint32
		want riscv.Instr
	}{
		{0x403100b3, riscv.Instr{Opcode: 0x33, Rd: 1, Funct3: 0, Rs1: 2, Rs2: 3, Funct7: 0x20}},
		{0xffffffff, riscv.Instr{Opcode: 0x7f, Rd: 31, Funct3: 7, Rs1: 31, Rs2: 31, Funct7: 0x7f}},
		{0x00000f80, riscv.Instr{Rd: 31}},
		{0x00007000, riscv.Instr{Funct3: 7}},
		{0x000f8000, riscv.Instr{Rs1: 31}},
		{0x01f00000, riscv.Instr{Rs2: 31}},
		{0xfe000000, riscv.Instr{Funct7: 0x7f}},
		{0x00032283, riscv.Instr{Opcode: 0x03, Rd: 5, Funct3: 2, Rs1: 6}}, // 真正的 lw x5, 0(x6)
	}
	for _, c := range cases {
		if got := riscv.Decode(c.word); got != c.want {
			t.Errorf("Decode(0x%08x) = %+v, want %+v", c.word, got, c.want)
		}
	}
}

func TestEncodeIsTheInverseOfDecode(t *testing.T) {
	words := []uint32{0x00728433, 0x403100b3, 0x00032283, 0xffffffff, 0x00000000, 0x12345678, 0x9abcdef0}
	for _, word := range words {
		if got := riscv.Encode(riscv.Decode(word)); got != word {
			t.Errorf("Encode(Decode(0x%08x)) = 0x%08x, want 原样拼回", word, got)
		}
	}

	add := riscv.Instr{Opcode: 0x33, Rd: 8, Funct3: 0, Rs1: 5, Rs2: 7, Funct7: 0}
	if got := riscv.Encode(add); got != 0x00728433 {
		t.Errorf("Encode(add x8, x5, x7) = 0x%08x, want 0x00728433", got)
	}

	// 超出位宽的字段只保留低位，不能污染相邻字段。
	overflow := riscv.Instr{Opcode: 0x33, Rd: 32 + 8, Rs1: 5, Rs2: 7}
	if got := riscv.Encode(overflow); got != 0x00728433 {
		t.Errorf("Encode(Rd = 40) = 0x%08x, want 0x00728433；Rd 只有 5 位，先 & 0x1f 再移位", got)
	}
}

func TestDisasm(t *testing.T) {
	cases := []struct {
		word uint32
		want string
	}{
		{0x00728433, "add x8, x5, x7"},
		{0x403100b3, "sub x1, x2, x3"},
		{0x01f00fb3, "add x31, x0, x31"},
		{0x00032283, "unknown 0x00032283"}, // lw：opcode 不是 0x33
		{0x02728433, "unknown 0x02728433"}, // funct7 = 1：既不是 add 也不是 sub
		{0x00729433, "unknown 0x00729433"}, // funct3 = 1
	}
	for _, c := range cases {
		if got := riscv.Disasm(riscv.Decode(c.word)); got != c.want {
			t.Errorf("Disasm(Decode(0x%08x)) = %q, want %q", c.word, got, c.want)
		}
	}
}

// TestBridgeFetchDecodeFromMemory 把三步串起来：从一段“内存”里每 4 字节取一条指令，
// 拼成机器码，译码，翻译成汇编。这就是 CPU 的“取指 → 译码”，下一题 test02 再加上“执行”。
func TestBridgeFetchDecodeFromMemory(t *testing.T) {
	memory := []byte{
		0x33, 0x84, 0x72, 0x00,
		0xb3, 0x00, 0x31, 0x40,
	}
	want := []string{"add x8, x5, x7", "sub x1, x2, x3"}

	for pc := 0; pc+4 <= len(memory); pc += 4 {
		var b [4]byte
		copy(b[:], memory[pc:pc+4])
		word := riscv.LittleEndian(b)
		got := riscv.Disasm(riscv.Decode(word))
		t.Logf("PC = 0x%04x  bytes = % x  word = 0x%08x  %s", pc, b, word, got)
		if got != want[pc/4] {
			t.Fatalf("PC = 0x%04x 处译出 %q, want %q", pc, got, want[pc/4])
		}
	}
	t.Log("PC 每次加 4：这就是课件里的“更新 PC = PC + 4”。")
}
