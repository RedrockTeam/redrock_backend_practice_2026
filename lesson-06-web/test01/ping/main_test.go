package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// doGet 用 httptest 直接请求路由表，返回响应码和响应体。
func doGet(t *testing.T, router http.Handler, target string) (int, string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))
	return recorder.Code, recorder.Body.String()
}

func TestPing(t *testing.T) {
	code, body := doGet(t, NewRouter(), "/ping")
	if code != http.StatusOK {
		t.Fatalf("GET /ping 状态码 = %d, want %d", code, http.StatusOK)
	}
	if body != "pong" {
		t.Fatalf("GET /ping 响应体 = %q, want %q", body, "pong")
	}
}

func TestUserByID(t *testing.T) {
	for _, id := range []string{"42", "hello", "redrock"} {
		code, body := doGet(t, NewRouter(), "/user/"+id)
		if code != http.StatusOK {
			t.Fatalf("GET /user/%s 状态码 = %d, want %d", id, code, http.StatusOK)
		}
		if body != id {
			t.Fatalf("GET /user/%s 响应体 = %q, want %q", id, body, id)
		}
	}
}

func TestAdd(t *testing.T) {
	for _, pair := range []struct{ num1, num2, want string }{
		{"3", "5", "8"},
		{"-2", "10", "8"},
		{"0", "0", "0"},
	} {
		target := fmt.Sprintf("/ping/add?num1=%s&num2=%s", pair.num1, pair.num2)
		code, body := doGet(t, NewRouter(), target)
		if code != http.StatusOK {
			t.Fatalf("GET %s 状态码 = %d, want %d", target, code, http.StatusOK)
		}
		if body != pair.want {
			t.Fatalf("GET %s 响应体 = %q, want %q", target, body, pair.want)
		}
	}
}
