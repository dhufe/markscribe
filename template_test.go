package main

import (
	"reflect"
	"testing"
	"time"
)

func TestHumanized_Today(t *testing.T) {
	now := time.Now()
	got := humanized(now.Add(-2 * time.Hour))
	if got != "today" {
		t.Fatalf("expected 'today', got %q", got)
	}
}

func TestHumanized_OlderThanOneDay(t *testing.T) {
	now := time.Now()
	twoDaysAgo := now.AddDate(0, 0, -2)
	got := humanized(twoDaysAgo)
	// Should not say 'today' and should include 'ago'
	if got == "today" || !contains(got, "ago") {
		t.Fatalf("expected relative time ending with 'ago', got %q", got)
	}
}

// contains is a tiny helper to avoid importing strings in tests since the
// project already exposes a template helper named contains.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || indexOf(s, sub) >= 0)
}

// simple substring search to avoid extra deps in test
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestHumanized_NonTime(t *testing.T) {
	if got := humanized(42); got != "42" {
		t.Fatalf("expected '42', got %q", got)
	}
}

func TestReverse_Ints(t *testing.T) {
	in := []int{1, 2, 3, 4}
	out := reverse(in).([]int)
	want := []int{4, 3, 2, 1}
	if !reflect.DeepEqual(in, want) { // in-place
		t.Fatalf("expected in-place reversal %v, got %v", want, in)
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("expected return value %v, got %v", want, out)
	}
}

func TestReverse_Strings(t *testing.T) {
	in := []string{"a", "b", "c"}
	want := []string{"c", "b", "a"}
	_ = reverse(in)
	if !reflect.DeepEqual(in, want) {
		t.Fatalf("expected %v, got %v", want, in)
	}
}

func TestReverse_Empty(t *testing.T) {
	in := []int{}
	_ = reverse(in)
	if len(in) != 0 {
		t.Fatalf("expected empty slice to remain empty, got %v", in)
	}
}

func TestReverse_SingleElement(t *testing.T) {
	in := []string{"only"}
	_ = reverse(in)
	if len(in) != 1 || in[0] != "only" {
		t.Fatalf("expected single element slice unchanged, got %v", in)
	}
}

func TestReverse_NilSlice(t *testing.T) {
	var in []int
	_ = reverse(in)
	if in != nil || len(in) != 0 {
		t.Fatalf("expected nil slice to remain nil/empty, got %#v", in)
	}
}
