package pitfalls

// Incrementer 是“能加一，能读当前值”的东西。
type Incrementer interface {
	Inc()
	Value() int
}

// Tally 是一个计数器。
type Tally struct {
	n int
}

// Inc 让计数加一。
//
// BUG: 这个方法能编译、能满足 Incrementer，但调用之后计数从来不变。
func (t Tally) Inc() {
	t.n++
}

// Value 返回当前计数。
func (t Tally) Value() int {
	return t.n
}

// NewTally 返回一个满足 Incrementer 的计数器。
// 修好 Inc 之后，这里返回的东西可能也要跟着改，否则编译不过。
func NewTally() Incrementer {
	return Tally{}
}

// IncTimes 对 i 调用 n 次 Inc。
func IncTimes(i Incrementer, n int) {
	for k := 0; k < n; k++ {
		i.Inc()
	}
}
