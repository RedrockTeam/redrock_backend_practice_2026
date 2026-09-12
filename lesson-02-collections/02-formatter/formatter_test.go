package formatter

import "testing"

func TestFormatters(t *testing.T) {
	student := Student{Name: "Lin", Score: 95}
	text, err := (TextFormatter{}).Format(student)
	if err != nil || text != "Lin: 95" {
		t.Fatalf("text = %q, %v", text, err)
	}
	jsonText, err := (JSONFormatter{}).Format(student)
	if err != nil || jsonText != `{"name":"Lin","score":95}` {
		t.Fatalf("json = %q, %v", jsonText, err)
	}
}
