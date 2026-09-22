package visibility_test

import (
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-01-basics/test06/visibility"
)

// 外部测试包只能通过导出的名称使用 visibility。
func TestExportedUserAPI(t *testing.T) {
	user := visibility.NewUser("小登", 18)

	if user.Name != "小登" {
		t.Fatalf("NewUser().Name = %q, want %q", user.Name, "小登")
	}
	if got := user.Age(); got != 18 {
		t.Fatalf("NewUser().Age() = %d, want 18", got)
	}

	user.Grow()
	if got := user.Age(); got != 19 {
		t.Fatalf("Age() after Grow() = %d, want 19", got)
	}
}

// Stage 是导出方法，它内部依赖的 stage 函数在另一个文件里，对外不可见。
func TestExportedStage(t *testing.T) {
	if got := visibility.NewUser("小登", 17).Stage(); got != "minor" {
		t.Fatalf("NewUser(17).Stage() = %q, want %q", got, "minor")
	}
	if got := visibility.NewUser("小登", 18).Stage(); got != "adult" {
		t.Fatalf("NewUser(18).Stage() = %q, want %q", got, "adult")
	}
}

// 这一题不能靠「把小写改成大写」通过。
// 反射可以看到未导出字段的存在，但普通代码无法读写它。
func TestAgeFieldStaysUnexported(t *testing.T) {
	userType := reflect.TypeOf(visibility.User{})

	field, ok := userType.FieldByName("age")
	if !ok {
		t.Fatal("User 仍然需要一个小写字段 age，不要把它改成大写来绕过可见性")
	}
	if field.PkgPath == "" {
		t.Fatal("字段 age 必须保持未导出")
	}
	if _, exported := userType.FieldByName("Age"); exported {
		t.Fatal("不要新增导出字段 Age；年龄只能通过 Age() 方法读取")
	}
}
