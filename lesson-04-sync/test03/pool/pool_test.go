package pool

import (
	"context"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	jobs := []Job{func(context.Context) (int, error) { return 1, nil }, func(context.Context) (int, error) { return 0, errors.New("bad") }}
	got := Run(context.Background(), 2, jobs)
	if len(got) != 2 || got[0].Value != 1 || got[1].Err == nil {
		t.Fatalf("Run() = %#v", got)
	}
}
