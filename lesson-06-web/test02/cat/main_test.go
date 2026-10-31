package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	image, err := os.ReadFile("cat.jpg")
	if err != nil {
		t.Fatalf("找不到 cat.jpg（应和 main.go 同目录）：%v", err)
	}

	recorder := httptest.NewRecorder()
	NewRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/cat", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /cat 状态码 = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "image/jpeg") {
		t.Fatalf("GET /cat 的 Content-Type = %q, want image/jpeg", got)
	}
	if body := recorder.Body.String(); body != string(image) {
		t.Fatalf("GET /cat 响应体与 cat.jpg 不一致（响应 %d 字节，文件 %d 字节）", len(body), len(image))
	}
}
