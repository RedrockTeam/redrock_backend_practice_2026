package ledger

// Ledger 记录多个账户的余额，以及所有账户的总额。
//
// 下面的实现没有任何同步。多个 goroutine 同时调用时：
//   - total 会丢更新（读-改-写不是一步）；
//   - 并发写 balances 这个 map 会让整个进程直接退出；
//   - Snapshot 返回的是内部 map，调用方一改就绕过了账本。
//
// 请用本课学的同步工具修复它。注意：读也要保护，Snapshot 要返回副本。
type Ledger struct {
	balances map[string]int
	total    int
}

// Deposit 给 account 存入 amount。
func (l *Ledger) Deposit(account string, amount int) {
	if l.balances == nil {
		l.balances = make(map[string]int)
	}
	l.balances[account] += amount
	l.total += amount
}

// Balance 返回 account 的余额；账户不存在时返回 0。
func (l *Ledger) Balance(account string) int {
	return l.balances[account]
}

// Total 返回所有账户的总额。
func (l *Ledger) Total() int {
	return l.total
}

// Snapshot 返回当前所有账户余额的副本和总额。
// 调用方修改返回的 map 不能影响 Ledger 内部状态。
func (l *Ledger) Snapshot() (map[string]int, int) {
	return l.balances, l.total
}
