package main

import "testing"

func TestNormal(t *testing.T) {
	res := Capitalize("foo")
	if res != "Foo" {
		t.Fatalf("wrong result, want %q, got %q", "Foo", res)
	}

	res = Capitalize("bar")
	if res != "Bar" {
		t.Fatalf("wrong result, want %q, got %q", "Bar", res)
	}
}

func TestFail(t *testing.T) {
	res := Capitalize("foo")
	if res != "Foo" {
		t.Fatalf("wrong result, want %q, got %q", "Foo", res)
	}

	res = Capitalize("bar")
	if res != "Bar" {
		t.Fatalf("wrong result, want %q, got %q", "Bar", res)
	}

	res = Capitalize("österreich")
	if res != "Österreich" {
		t.Fatalf("wrong result, want %q, got %q", "Österreich", res)
	}
}

func TestTables(t *testing.T) {
	var tests = []struct {
		Input string
		Want  string
	}{
		{"foo", "Foo"},
		{"bar", "Bar"},
		{"österreich", "Österreich"},
	}

	for _, test := range tests {
		result := Capitalize(test.Input)
		if result != test.Want {
			t.Errorf("wrong result, want %q, got %q", test.Want, result)
		}
	}
}

func TestSubtests(t *testing.T) {
	var tests = []struct {
		Input string
		Want  string
	}{
		{"foo", "Foo"},
		{"bar", "Bar"},
		{"österreich", "Österreich"},
	}

	for _, test := range tests {
		t.Run("", func(st *testing.T) {
			result := Capitalize(test.Input)
			if result != test.Want {
				st.Fatalf("wrong result, want %q, got %q", test.Want, result)
			}
		})
	}
}
