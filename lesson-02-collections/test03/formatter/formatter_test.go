package formatter_test

import (
	"errors"
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test03/formatter"
)

func TestTextAndJSONFormatters(t *testing.T) {
	student := formatter.Student{Name: "Lin", Score: 95}

	text, err := (formatter.TextFormatter{}).Format(student)
	if err != nil {
		t.Fatalf("TextFormatter.Format() error = %v", err)
	}
	if text != "Lin: 95" {
		t.Fatalf("TextFormatter.Format() = %q, want %q", text, "Lin: 95")
	}

	jsonText, err := (formatter.JSONFormatter{}).Format(student)
	if err != nil {
		t.Fatalf("JSONFormatter.Format() error = %v", err)
	}
	if jsonText != `{"name":"Lin","score":95}` {
		t.Fatalf("JSONFormatter.Format() = %q, want %q", jsonText, `{"name":"Lin","score":95}`)
	}
}

func TestFormatAllUsesTheInterface(t *testing.T) {
	student := formatter.Student{Name: "Lin", Score: 95}
	formatters := []formatter.Formatter{
		formatter.TextFormatter{},
		formatter.JSONFormatter{},
	}

	got, err := formatter.FormatAll(formatters, student)
	if err != nil {
		t.Fatalf("FormatAll() error = %v", err)
	}
	want := []string{"Lin: 95", `{"name":"Lin","score":95}`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FormatAll() = %#v, want %#v", got, want)
	}
}

type brokenFormatter struct{}

func (brokenFormatter) Format(formatter.Student) (string, error) {
	return "", errors.New("boom")
}

func TestFormatAllStopsOnError(t *testing.T) {
	student := formatter.Student{Name: "Lin", Score: 95}
	formatters := []formatter.Formatter{
		formatter.TextFormatter{},
		brokenFormatter{},
		formatter.JSONFormatter{},
	}

	got, err := formatter.FormatAll(formatters, student)
	if err == nil {
		t.Fatal("FormatAll() error = nil, want 非空错误；第二个 formatter 已经失败")
	}
	want := []string{"Lin: 95"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FormatAll() 在出错前的结果 = %#v, want %#v", got, want)
	}
}

func TestFormatAllHandlesEmpty(t *testing.T) {
	got, err := formatter.FormatAll(nil, formatter.Student{Name: "Lin", Score: 95})
	if err != nil {
		t.Fatalf("FormatAll(nil) error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("FormatAll(nil) = %#v, want 空结果", got)
	}
}
