package ledger_test

import (
	"fmt"
	"sync"
	"testing"

	"redrock/teaching-exercises/lesson-04-sync/test03/ledger"
)

func TestZeroValueLedger(t *testing.T) {
	var l ledger.Ledger

	if got := l.Balance("alice"); got != 0 {
		t.Fatalf(`零值 Ledger.Balance("alice") = %d, want 0`, got)
	}
	if got := l.Total(); got != 0 {
		t.Fatalf("零值 Ledger.Total() = %d, want 0", got)
	}

	l.Deposit("alice", 10)
	if got := l.Balance("alice"); got != 10 {
		t.Fatalf(`Deposit 后 Balance("alice") = %d, want 10`, got)
	}
	if got := l.Total(); got != 10 {
		t.Fatalf("Deposit 后 Total() = %d, want 10", got)
	}
}

func TestSequentialDeposits(t *testing.T) {
	var l ledger.Ledger
	l.Deposit("alice", 10)
	l.Deposit("bob", 5)
	l.Deposit("alice", 7)

	if got := l.Balance("alice"); got != 17 {
		t.Fatalf(`Balance("alice") = %d, want 17`, got)
	}
	if got := l.Balance("bob"); got != 5 {
		t.Fatalf(`Balance("bob") = %d, want 5`, got)
	}
	if got := l.Balance("nobody"); got != 0 {
		t.Fatalf(`Balance("nobody") = %d, want 0`, got)
	}
	if got := l.Total(); got != 22 {
		t.Fatalf("Total() = %d, want 22", got)
	}
}

// TestSnapshotIsACopy 检查 Snapshot 交出的是副本，不是内部 map 的别名。
// 这是所有权问题：调用方拿到 map 之后，账本不能因此被绕过。
func TestSnapshotIsACopy(t *testing.T) {
	var l ledger.Ledger
	l.Deposit("alice", 10)

	snapshot, total := l.Snapshot()
	if snapshot == nil {
		t.Fatal("Snapshot() 返回了 nil map；空账本也应该返回一个可写的空 map")
	}
	if total != 10 {
		t.Fatalf("Snapshot() 的 total = %d, want 10", total)
	}
	if got := snapshot["alice"]; got != 10 {
		t.Fatalf(`snapshot["alice"] = %d, want 10`, got)
	}

	snapshot["alice"] = 999
	snapshot["mallory"] = 1

	if got := l.Balance("alice"); got != 10 {
		t.Fatalf(`修改 snapshot 后 Balance("alice") = %d, want 10；Snapshot 不能把内部 map 直接交出去`, got)
	}
	if got := l.Total(); got != 10 {
		t.Fatalf("修改 snapshot 后 Total() = %d, want 10", got)
	}

	again, _ := l.Snapshot()
	if len(again) != 1 || again["alice"] != 10 {
		t.Fatalf("第二次 Snapshot() = %v，内部状态被调用方改坏了", again)
	}
}

// TestConcurrentUseKeepsInvariants 是所有并发问题的合集：
//   - 多个写者同时存款，不能丢更新；
//   - 写 map 不能把进程打死；
//   - 读者拿到的 Snapshot 必须自洽：所有账户之和 == 同一份快照的 total；
//   - 读者写的副本不能影响账本。
func TestConcurrentUseKeepsInvariants(t *testing.T) {
	const (
		rounds    = 5
		writers   = 16
		readers   = 4
		perWorker = 200
		accounts  = 4
	)

	wantEach := writers * perWorker / accounts
	wantTotal := writers * perWorker

	for round := 0; round < rounds; round++ {
		var l ledger.Ledger

		start := make(chan struct{})
		errCh := make(chan string, 1)
		report := func(format string, args ...any) {
			select {
			case errCh <- fmt.Sprintf(format, args...):
			default:
			}
		}

		var wg sync.WaitGroup

		for w := 0; w < writers; w++ {
			w := w
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for n := 0; n < perWorker; n++ {
					l.Deposit(fmt.Sprintf("acc-%d", (w+n)%accounts), 1)
				}
			}()
		}

		for r := 0; r < readers; r++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for n := 0; n < perWorker; n++ {
					snapshot, total := l.Snapshot()
					sum := 0
					for _, balance := range snapshot {
						sum += balance
					}
					if sum != total {
						report("Snapshot 不自洽：账户之和 = %d, total = %d（锁没有覆盖完整状态）", sum, total)
						return
					}
				}
			}()
		}

		close(start)
		wg.Wait()

		select {
		case message := <-errCh:
			t.Fatalf("round %d: %s", round, message)
		default:
		}

		if got := l.Total(); got != wantTotal {
			t.Fatalf("round %d: Total() = %d, want %d（有更新丢失）", round, got, wantTotal)
		}
		for i := 0; i < accounts; i++ {
			account := fmt.Sprintf("acc-%d", i)
			if got := l.Balance(account); got != wantEach {
				t.Fatalf("round %d: Balance(%q) = %d, want %d", round, account, got, wantEach)
			}
		}

		snapshot, total := l.Snapshot()
		if total != wantTotal {
			t.Fatalf("round %d: 最终 Snapshot total = %d, want %d", round, total, wantTotal)
		}
		sum := 0
		for _, balance := range snapshot {
			sum += balance
		}
		if sum != wantTotal {
			t.Fatalf("round %d: 最终 Snapshot 账户之和 = %d, want %d", round, sum, wantTotal)
		}
	}

	t.Logf("%d 轮、每轮 %d 个写者和 %d 个读者：总额和快照都保持一致。", rounds, writers, readers)
}
