package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// doGet 用 httptest 直接请求路由表，返回响应码、响应体和这次请求期间打出的日志。
// 日志是中间件的职责，光看响应分辨不出来，所以把 stdout 接到管道里读回来。
func doGet(t *testing.T, router http.Handler, target string) (int, string, string) {
	t.Helper()

	previousStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建输出管道失败：%v", err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = previousStdout }()

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))

	if err := writer.Close(); err != nil {
		t.Fatalf("关闭输出管道失败：%v", err)
	}
	logs, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("读取日志失败：%v", err)
	}
	return recorder.Code, recorder.Body.String(), string(logs)
}

func TestLogEven(t *testing.T) {
	router := NewRouter() // 引擎构造在捕获之外，[GIN-debug] 输出不算日志
	for _, num := range []string{"2", "10", "42"} {
		code, body, logs := doGet(t, router, "/log/"+num)
		if code != http.StatusOK {
			t.Fatalf("GET /log/%s 状态码 = %d, want %d（偶数应放行）", num, code, http.StatusOK)
		}
		if body != `{"message":"偶数"}` {
			t.Fatalf("GET /log/%s 响应体 = %s, want {\"message\":\"偶数\"}", num, body)
		}
		if !strings.Contains(logs, num) {
			t.Fatalf("GET /log/%s 的日志里没有出现 num（日志 = %q）", num, logs)
		}
	}
}

func TestLogOdd(t *testing.T) {
	router := NewRouter()
	for _, num := range []string{"1", "3", "99"} {
		code, body, logs := doGet(t, router, "/log/"+num)
		if code != http.StatusForbidden {
			t.Fatalf("GET /log/%s 状态码 = %d, want %d（奇数要由中间件拦截）", num, code, http.StatusForbidden)
		}
		if body != `{"error":"奇数被拦截"}` {
			t.Fatalf("GET /log/%s 响应体 = %s, want {\"error\":\"奇数被拦截\"}", num, body)
		}
		if !strings.Contains(logs, num) {
			t.Fatalf("GET /log/%s 的日志里没有出现 num（日志 = %q）", num, logs)
		}
	}
}

func TestLogNotANumber(t *testing.T) {
	code, body, _ := doGet(t, NewRouter(), "/log/abc")
	if code != http.StatusBadRequest {
		t.Fatalf("GET /log/abc 状态码 = %d, want %d（不是整数应被拦截）", code, http.StatusBadRequest)
	}
	if body != `{"error":"num 不是整数"}` {
		t.Fatalf("GET /log/abc 响应体 = %s, want {\"error\":\"num 不是整数\"}", body)
	}
}
