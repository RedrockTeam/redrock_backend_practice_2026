package pitfalls_test

import (
	"errors"
	"fmt"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test07/pitfalls"
)

func TestValidateReturnsRealNil(t *testing.T) {
	err := pitfalls.Validate(pitfalls.User{Name: "Lin", Age: 18})
	if err != nil {
		t.Fatalf("合法 User 的 Validate() = %v（%T）, err != nil 为 true；接口里装了一个 nil 指针，接口本身不是 nil", err, err)
	}
}

func TestValidateReportsField(t *testing.T) {
	cases := []struct {
		user  pitfalls.User
		field string
	}{
		{pitfalls.User{Name: "", Age: 18}, "name"},
		{pitfalls.User{Name: "Lin", Age: -1}, "age"},
		{pitfalls.User{Name: "Lin", Age: 151}, "age"},
		{pitfalls.User{Name: "", Age: 999}, "name"},
	}
	for _, c := range cases {
		err := pitfalls.Validate(c.user)
		var ve *pitfalls.ValidationError
		if !errors.As(err, &ve) || ve == nil {
			t.Errorf("Validate(%+v) = %v, want *ValidationError", c.user, err)
			continue
		}
		if ve.Field != c.field {
			t.Errorf("Validate(%+v).Field = %q, want %q", c.user, ve.Field, c.field)
		}
	}
	if err := pitfalls.Validate(pitfalls.User{Name: "Lin", Age: 0}); err != nil {
		t.Errorf("Age = 0 是合法的，Validate() = %v", err)
	}
}

func TestTallyCountsThroughInterface(t *testing.T) {
	counter := pitfalls.NewTally()
	pitfalls.IncTimes(counter, 3)
	if got := counter.Value(); got != 3 {
		t.Fatalf("IncTimes(counter, 3) 之后 Value() = %d, want 3；值接收者的方法改的是副本", got)
	}
}

func TestPointerSatisfiesIncrementer(t *testing.T) {
	var counter pitfalls.Incrementer = &pitfalls.Tally{}
	counter.Inc()
	counter.Inc()
	if got := counter.Value(); got != 2 {
		t.Fatalf("(*Tally).Inc 两次之后 Value() = %d, want 2", got)
	}
}

func TestPointImplementsStringer(t *testing.T) {
	var p any = pitfalls.Point{X: 1, Y: 2}
	if _, ok := p.(fmt.Stringer); !ok {
		t.Fatal("Point 没有实现 fmt.Stringer；需要一个 String() string 方法")
	}
	if got := fmt.Sprint(pitfalls.Point{X: 3, Y: -4}); got != "(3, -4)" {
		t.Fatalf("fmt.Sprint(Point{3, -4}) = %q, want %q", got, "(3, -4)")
	}
}

func TestDescribeTypeSwitch(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, "nil"},
		{errors.New("boom"), "error: boom"},
		{pitfalls.Point{X: 1, Y: 2}, "stringer: (1, 2)"},
		{42, "int: 42"},
		{"hi", `string: "hi"`},
		{[]int{1, 2, 3}, "[]int of len 3"},
		{[]int(nil), "[]int of len 0"},
		{3.5, "other: float64"},
		{map[string]int{}, "other: map[string]int"},
	}
	for _, c := range cases {
		if got := pitfalls.Describe(c.in); got != c.want {
			t.Errorf("Describe(%#v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func safeEqual(t *testing.T, a, b any) (equal, comparable bool) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SafeEqual(%v, %v) panic: %v；不可比较的类型不能直接用 ==", a, b, r)
		}
	}()
	return pitfalls.SafeEqual(a, b)
}

func TestSafeEqualComparable(t *testing.T) {
	cases := []struct {
		a, b any
		want bool
	}{
		{1, 1, true},
		{1, 2, false},
		{"a", "a", true},
		{1, "1", false},
		{pitfalls.Point{X: 1, Y: 2}, pitfalls.Point{X: 1, Y: 2}, true},
		{nil, nil, true},
	}
	for _, c := range cases {
		equal, comparable := safeEqual(t, c.a, c.b)
		if !comparable || equal != c.want {
			t.Errorf("SafeEqual(%v, %v) = %v, %v, want %v, true", c.a, c.b, equal, comparable, c.want)
		}
	}
}

func TestSafeEqualUncomparableDoesNotPanic(t *testing.T) {
	cases := [][2]any{
		{[]int{1, 2}, []int{1, 2}},
		{map[string]int{"a": 1}, map[string]int{"a": 1}},
	}
	for _, c := range cases {
		equal, comparable := safeEqual(t, c[0], c[1])
		if equal || comparable {
			t.Errorf("SafeEqual(%v, %v) = %v, %v, want false, false", c[0], c[1], equal, comparable)
		}
	}
}

// TestBridgeInterfaceHoldsTypeAndValue 不检查你的实现：接口值 = (动态类型, 动态值)。
// 只有两者都为空，接口才 == nil。这正是 Validate 那个 bug 的根源。
func TestBridgeInterfaceHoldsTypeAndValue(t *testing.T) {
	var p *pitfalls.ValidationError
	var err error = p
	if err == nil {
		t.Fatal("装了 (*ValidationError)(nil) 的接口居然 == nil？")
	}
	t.Logf("err 的动态类型 = %T，动态值 = %v，err == nil 为 %v", err, p, err == nil)
	t.Log("所以返回 error 的函数，成功时要直接 return nil，而不是返回一个类型化的 nil 指针。")
}
