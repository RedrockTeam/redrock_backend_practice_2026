package missing_test

import (
	"crypto/sha256"
	"encoding/hex"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"redrock/teaching-exercises/lesson-03-goroutines/test07/missing"
)

// TestBridgeObserveCollect 不要求 Collect 正确——它就是这一题的观察对象。
// 它用你写的 Verify 数出每一轮少了多少个结果，用 go test -v 看输出。
func TestBridgeObserveCollect(t *testing.T) {
	if raceEnabled {
		t.Skip("Collect 本身就有数据竞争（这正是观察对象），-race 下跳过；去掉 -race 再看输出")
	}
	const size = 200

	values := make([]int, size)
	for i := range values {
		values[i] = i + 1
	}

	for round := 0; round < 5; round++ {
		got := missing.Collect(values)
		if len(got) != size {
			t.Fatalf("round %d: len(Collect(...)) = %d, want %d", round, len(got), size)
		}
		t.Logf("round %d: %d 个结果里有 %d 个还没写好", round, size, missing.Verify(values, got))
	}
	t.Log("每一轮的数字可能都不一样：结果取决于调度，而不是你的代码。")
	t.Log("第四课 test01 会用 WaitGroup 让这个数字稳定地变成 0。")
}

func TestVerifyCountsWrongPositions(t *testing.T) {
	cases := []struct {
		name   string
		values []int
		got    []int
		want   int
	}{
		{"全部正确", []int{1, 2, 3}, []int{1, 4, 9}, 0},
		{"全是零值", []int{1, 2, 3}, []int{0, 0, 0}, 3},
		{"中间一个错", []int{1, 2, 3}, []int{1, 5, 9}, 1},
		{"负数的平方", []int{-2, 0}, []int{4, 0}, 0},
		{"0 的平方本来就是 0", []int{0, 0}, []int{0, 0}, 0},
		{"got 太短", []int{1, 2, 3, 4}, []int{1, 4}, 2},
		{"got 是 nil", []int{1, 2}, nil, 2},
		{"got 更长，多的不管", []int{1}, []int{1, 99}, 0},
		{"都为空", nil, nil, 0},
	}
	for _, c := range cases {
		if got := missing.Verify(c.values, c.got); got != c.want {
			t.Errorf("%s: Verify(%v, %v) = %d, want %d", c.name, c.values, c.got, got, c.want)
		}
	}
}

// 答案以摘要形式保存，失败信息只告诉你哪一题不对，不告诉你正确选项。
var answerDigests = map[string]string{
	"Q1": "f532ebc92b4c200499ad5eac3cd7c847d2129f070c1506d96a57288447483bfd",
	"Q2": "5f3bdb3a64f4d1ec191c3ef3276fa106d38af3ca716ee7d93437ce1c2cafb407",
	"Q3": "3811d47c5a0bcd5082e0a161581ce3f22e2594abaa895894cfbce1170fa4cb83",
	"Q4": "3a67210d69dc5da5814aeda9243c0070c6569fc96daf8146d9f89756930d7f3c",
	"Q5": "f9eda4cd02197420ebe78d42a5b757ea6f088c3770d4d95dde2e382016cd5c5b",
	"Q6": "67b483ad147878f586d86b1ffc941f14e6992bf75f41b84d0131d1fa7c951695",
}

func checkAnswers(t *testing.T, answers map[string]string) {
	t.Helper()
	for _, question := range []string{"Q1", "Q2", "Q3", "Q4", "Q5", "Q6"} {
		answer, ok := answers[question]
		if !ok {
			continue
		}
		switch answer {
		case "A", "B", "C":
		case "":
			t.Errorf("%s 还没作答", question)
			continue
		default:
			t.Errorf("%s = %q；答案只能是 \"A\"、\"B\" 或 \"C\"", question, answer)
			continue
		}
		sum := sha256.Sum256([]byte("missing:" + question + ":" + answer))
		if hex.EncodeToString(sum[:]) != answerDigests[question] {
			t.Errorf("%s 的答案不对；回到 README 对应的课件段落再想一想", question)
		}
	}
}

func TestObservationAnswers(t *testing.T) {
	checkAnswers(t, map[string]string{"Q1": missing.Q1, "Q2": missing.Q2, "Q3": missing.Q3})
}

func TestConceptAnswers(t *testing.T) {
	checkAnswers(t, map[string]string{"Q4": missing.Q4, "Q5": missing.Q5, "Q6": missing.Q6})
}

func TestCollectHandlesEmpty(t *testing.T) {
	if got := missing.Collect(nil); len(got) != 0 {
		t.Fatalf("Collect(nil) = %#v, want 空切片", got)
	}
}

// TestCollectDoesNotUseSleep 挡住「睡一会儿大概就好了」的修法。
// collect.go 是观察对象，这一题不需要改它；第四课 test01 会用 WaitGroup 建立“写 → 返回”的顺序。
func TestCollectDoesNotUseSleep(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("列目录失败: %v", err)
	}

	fset := token.NewFileSet()
	foundAt := ""
	checked := 0

	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		checked++

		file, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatalf("解析 %s 失败: %v", name, err)
		}

		timeNames := map[string]bool{"time": true}
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if path != "time" {
				continue
			}
			if imp.Name != nil {
				timeNames[imp.Name.Name] = true
			}
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			if timeNames[ident.Name] && selector.Sel.Name == "Sleep" {
				foundAt = fset.Position(call.Pos()).String()
			}
			return true
		})
	}

	if checked == 0 {
		t.Fatal("没有找到可解析的实现文件；请保留 collect.go")
	}
	if foundAt != "" {
		t.Fatalf("%s 使用了 time.Sleep。Sleep 只是让当前函数等一会儿，解决不了“函数已经返回、调用方看不到结果”的问题", foundAt)
	}
}
