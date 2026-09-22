package even

import "testing"

// TestIsEven 是需要由你亲手写完的第一个 Go 测试。
//
// 要求：
//  1. 至少分别检查 IsEven(2)、IsEven(3)、IsEven(0) 和 IsEven(-2) 四种输入；
//  2. 结果和预期不一致时，用 t.Error / t.Errorf / t.Fatal / t.Fatalf 报告失败；
//  3. 写完测试后先运行一次，看它失败；再去修 even.go，让它变绿。
//
// 期望的行为：2 是偶数，3 不是，0 是偶数，-2 也是偶数。
func TestIsEven(t *testing.T) {
	// TODO: 删掉下面这行，写出你自己的检查。
	//
	// 一个检查长这样：
	//
	//	if !IsEven(2) {
	//		t.Errorf("IsEven(2) = false, want true")
	//	}
	//
	// t.Error 报告失败但继续往下跑，t.Fatal 报告失败并立刻结束这个测试。
	t.Fatal("TODO: 完成 TestIsEven")
}
