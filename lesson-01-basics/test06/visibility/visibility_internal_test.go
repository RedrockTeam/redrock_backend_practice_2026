package visibility

import "testing"

// 同包测试可以直接访问小写字段、小写方法和小写函数。
func TestPackagePrivateGrow(t *testing.T) {
	user := &User{Name: "小登", age: 18}
	user.grow()

	if user.age != 19 {
		t.Fatalf("age after grow() = %d, want 19", user.age)
	}
}

// stage 定义在 stage.go，这个测试文件同属 visibility 包，所以可以直接调用它。
func TestPackagePrivateStage(t *testing.T) {
	if got := stage(adultAge - 1); got != "minor" {
		t.Fatalf("stage(%d) = %q, want %q", adultAge-1, got, "minor")
	}
	if got := stage(adultAge); got != "adult" {
		t.Fatalf("stage(%d) = %q, want %q", adultAge, got, "adult")
	}
	if got := stage(adultAge + 40); got != "adult" {
		t.Fatalf("stage(%d) = %q, want %q", adultAge+40, got, "adult")
	}
}
