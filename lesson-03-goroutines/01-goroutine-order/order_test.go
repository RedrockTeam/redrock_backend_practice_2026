package order

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRunKeepsOrder(t *testing.T) {
	tasks := []func() string{func() string { return fmt.Sprint(1) }, func() string { return fmt.Sprint(2) }}
	if got := Run(tasks); !reflect.DeepEqual(got, []string{"1", "2"}) {
		t.Fatalf("Run() = %#v", got)
	}
}
