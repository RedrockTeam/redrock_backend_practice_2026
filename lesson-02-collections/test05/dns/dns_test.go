package dns_test

import (
	"errors"
	"reflect"
	"testing"

	"redrock/teaching-exercises/lesson-02-collections/test05/dns"
)

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"www.baidu.com":     "www.baidu.com",
		"WWW.Baidu.COM":     "www.baidu.com",
		" www.baidu.com.  ": "www.baidu.com",
		"localhost.":        "localhost",
		"   ":               "",
	}
	for in, want := range cases {
		if got := dns.Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAddAndLookupAreCaseInsensitive(t *testing.T) {
	r := dns.New()
	if err := r.Add("WWW.Baidu.com.", "110.242.68.66"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	for _, query := range []string{"www.baidu.com", "WWW.BAIDU.COM", "www.baidu.com."} {
		ip, ok := r.Lookup(query)
		if !ok || ip != "110.242.68.66" {
			t.Errorf("Lookup(%q) = %q, %v, want 110.242.68.66, true；查询前也要规范化", query, ip, ok)
		}
	}
	if ip, ok := r.Lookup("missing.com"); ok || ip != "" {
		t.Errorf(`Lookup("missing.com") = %q, %v, want "", false`, ip, ok)
	}
	if r.Len() != 1 {
		t.Errorf("Len() = %d, want 1", r.Len())
	}
}

func TestAddRejectsEmptyInput(t *testing.T) {
	r := dns.New()
	_ = r.Add("ok.com", "1.1.1.1")
	for _, c := range [][2]string{{"", "1.1.1.1"}, {"  .", "1.1.1.1"}, {"a.com", ""}, {"a.com", "  "}} {
		if err := r.Add(c[0], c[1]); !errors.Is(err, dns.ErrInvalid) {
			t.Errorf("Add(%q, %q) error = %v, want ErrInvalid", c[0], c[1], err)
		}
	}
	if r.Len() != 1 {
		t.Fatalf("非法输入之后 Len() = %d, want 1；出错时不能修改 Resolver", r.Len())
	}
}

func TestOverwriteUpdatesReverseIndex(t *testing.T) {
	r := dns.New()
	_ = r.Add("a.com", "1.1.1.1")
	_ = r.Add("A.com", "2.2.2.2")

	if ip, _ := r.Lookup("a.com"); ip != "2.2.2.2" {
		t.Fatalf(`覆盖之后 Lookup("a.com") = %q, want 2.2.2.2`, ip)
	}
	if got := r.ReverseLookup("1.1.1.1"); len(got) != 0 {
		t.Fatalf(`覆盖之后 ReverseLookup("1.1.1.1") = %v, want 空；旧的反查记录没删掉`, got)
	}
	if got := r.ReverseLookup("2.2.2.2"); !reflect.DeepEqual(got, []string{"a.com"}) {
		t.Fatalf(`ReverseLookup("2.2.2.2") = %v, want [a.com]`, got)
	}
	if r.Len() != 1 {
		t.Fatalf("覆盖之后 Len() = %d, want 1", r.Len())
	}
}

func TestRemoveKeepsBothDirectionsInSync(t *testing.T) {
	r := dns.New()
	_ = r.Add("a.com", "1.1.1.1")
	_ = r.Add("b.com", "1.1.1.1")

	r.Remove("A.COM.")
	if _, ok := r.Lookup("a.com"); ok {
		t.Fatal(`Remove("A.COM.") 之后仍能查到 a.com；删除前也要规范化`)
	}
	if got := r.ReverseLookup("1.1.1.1"); !reflect.DeepEqual(got, []string{"b.com"}) {
		t.Fatalf(`Remove 之后 ReverseLookup("1.1.1.1") = %v, want [b.com]`, got)
	}
	r.Remove("never-existed.com") // 不应该 panic
	r.Remove("b.com")
	if got := r.ReverseLookup("1.1.1.1"); len(got) != 0 {
		t.Fatalf("全部删除之后 ReverseLookup = %v, want 空", got)
	}
	if r.Len() != 0 {
		t.Fatalf("全部删除之后 Len() = %d, want 0", r.Len())
	}
}

func TestReverseLookupIsSortedAndACopy(t *testing.T) {
	r := dns.New()
	for _, d := range []string{"cdn.example.com", "api.example.com", "www.example.com"} {
		_ = r.Add(d, "10.0.0.1")
	}
	_ = r.Add("other.com", "10.0.0.2")

	got := r.ReverseLookup(" 10.0.0.1 ")
	want := []string{"api.example.com", "cdn.example.com", "www.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReverseLookup(\"10.0.0.1\") = %v, want %v；结果要按字典序，IP 也要去空白", got, want)
	}

	got[0] = "hacked.com"
	if again := r.ReverseLookup("10.0.0.1"); len(again) == 0 || again[0] != "api.example.com" {
		t.Fatalf("修改返回的切片之后再查 = %v；ReverseLookup 不能交出内部切片", again)
	}
}

func TestDomainsSorted(t *testing.T) {
	r := dns.New()
	for _, d := range []string{"localhost", "www.baidu.com", "Api.Example.com"} {
		_ = r.Add(d, "127.0.0.1")
	}
	want := []string{"api.example.com", "localhost", "www.baidu.com"}
	if got := r.Domains(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Domains() = %v, want %v；map 的遍历顺序是随机的，要自己排序", got, want)
	}
	if got := dns.New().Domains(); len(got) != 0 {
		t.Fatalf("空 Resolver 的 Domains() = %v, want 空", got)
	}
}

func TestZeroValueResolver(t *testing.T) {
	var r dns.Resolver
	if _, ok := r.Lookup("a.com"); ok {
		t.Fatal("零值 Resolver 不应该查到任何东西")
	}
	r.Remove("a.com") // 不应该 panic
	if err := r.Add("a.com", "1.1.1.1"); err != nil {
		t.Fatalf("零值 Resolver.Add() error = %v", err)
	}
	if ip, ok := r.Lookup("a.com"); !ok || ip != "1.1.1.1" {
		t.Fatalf(`零值 Resolver 上 Lookup("a.com") = %q, %v, want 1.1.1.1, true`, ip, ok)
	}
	if got := r.ReverseLookup("1.1.1.1"); !reflect.DeepEqual(got, []string{"a.com"}) {
		t.Fatalf("零值 Resolver 上 ReverseLookup = %v, want [a.com]", got)
	}
}
