package lookup

import "testing"

func TestFindByName(t *testing.T) {
	students := []Student{{Name: "Lin", Score: 90}, {Name: "Wang", Score: 80}}
	got, ok := FindByName(students, "Wang")
	if !ok || got.Score != 80 {
		t.Fatalf("FindByName() = %+v, %v", got, ok)
	}
	if _, ok := FindByName(students, "None"); ok {
		t.Fatal("missing student should return false")
	}
}

func TestFormatStudents(t *testing.T) {
	got := FormatStudents([]Student{{Name: " Lin ", Score: 90}})
	if got != "Lin: 90\n" {
		t.Fatalf("FormatStudents() = %q", got)
	}
}
