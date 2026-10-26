package fanin_test

import (
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"redrock/teaching-exercises/lesson-04-sync/test02/fanin"
)

func readN(t *testing.T, out <-chan int, count int) []int {
	t.Helper()

	got := make([]int, 0, count)
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	for len(got) < count {
		select {
		case value, ok := <-out:
			if !ok {
				t.Fatalf("输出提前关闭：只收到 %v，还差 %d 个", got, count-len(got))
			}
			got = append(got, value)
		case <-timer.C:
			t.Fatalf("等了 2 秒只收到 %d/%d 个值；转发或关闭的逻辑有问题", len(got), count)
		}
	}
	return got
}

func expectClosed(t *testing.T, out <-chan int) {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	select {
	case value, ok := <-out:
		if ok {
			t.Fatalf("所有输入都结束后仍收到了值 %d；输出必须关闭", value)
		}
	case <-timer.C:
		t.Fatal("所有输入都结束后输出仍没有关闭；Merge 需要在所有转发 goroutine 退出后 close(out)")
	}
}

func TestMergeCollectsEveryValue(t *testing.T) {
	a := make(chan int, 2)
	b := make(chan int, 1)
	c := make(chan int, 3)

	a <- 1
	a <- 2
	close(a)
	b <- 3
	close(b)
	c <- 4
	c <- 5
	c <- 6
	close(c)

	out := fanin.Merge(a, b, c)
	got := readN(t, out, 6)
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("收集到的值 = %v, want %v（顺序可以不同，但一个都不能少、不能多）", got, want)
	}
	expectClosed(t, out)
}

func TestMergeClosesAfterAllInputs(t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	close(a)
	close(b)

	out := fanin.Merge(a, b)
	expectClosed(t, out)
}

func TestMergeWithNoInputs(t *testing.T) {
	// 没有输入时也要返回一个已经关闭的 channel，而不是让消费者永久等待。
	expectClosed(t, fanin.Merge())
}

func TestMergeStreamsUnbufferedInputs(t *testing.T) {
	const (
		inputs   = 3
		perInput = 50
	)

	channels := make([]chan int, inputs)
	for i := range channels {
		channels[i] = make(chan int)
	}

	var producers sync.WaitGroup
	producers.Add(inputs)
	for i := 0; i < inputs; i++ {
		i := i
		go func() {
			defer producers.Done()
			defer close(channels[i])
			for n := 0; n < perInput; n++ {
				channels[i] <- i*1000 + n
			}
		}()
	}

	out := fanin.Merge(channels[0], channels[1], channels[2])
	got := readN(t, out, inputs*perInput)
	expectClosed(t, out)
	producers.Wait()

	seen := make(map[int]bool, len(got))
	for _, value := range got {
		if seen[value] {
			t.Fatalf("值 %d 出现了不止一次", value)
		}
		seen[value] = true
	}
	for i := 0; i < inputs; i++ {
		for n := 0; n < perInput; n++ {
			value := i*1000 + n
			if !seen[value] {
				t.Fatalf("缺少值 %d；无缓冲输入下，Merge 必须边收边转发", value)
			}
		}
	}
}

func TestMergeHandlesSingleInput(t *testing.T) {
	in := make(chan int, 3)
	in <- 7
	in <- 8
	in <- 9
	close(in)

	out := fanin.Merge(in)
	got := readN(t, out, 3)
	sort.Ints(got)
	want := []int{7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("单输入时收集到 %v, want %v", got, want)
	}
	expectClosed(t, out)
}
