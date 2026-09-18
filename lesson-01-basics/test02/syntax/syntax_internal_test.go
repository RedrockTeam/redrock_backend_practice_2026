package syntax

import "testing"

// 这个测试文件声明的是 package syntax，和被测代码同属一个包，
// 因此可以直接调用小写的 add。外部测试文件做不到这一点。
func TestPackagePrivateAdd(t *testing.T) {
	if got := add(5, 6); got != 11 {
		t.Fatalf("add(5, 6) = %d, want 11", got)
	}
}
