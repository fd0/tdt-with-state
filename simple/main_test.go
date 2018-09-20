package main

import "testing"

func TestCapitalize(t *testing.T) {
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
		t.Fatalf("wrong result, want %q, got %q", "Österreich", res) // HL
	}
}

func TestTables(t *testing.T) {
	var tests = []struct {
		Input string
		Want  string
	}{
		{"foo", "Foo"},
		{"österreich", "Österreich"},
		{"Österreich", "Österreich"},
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
		{"österreich", "Österreich"},
		{"Österreich", "Österreich"},
	}

	for _, test := range tests {
		t.Run("", func(st *testing.T) { // HL
			result := Capitalize(test.Input)
			if result != test.Want {
				st.Fatalf("wrong result, want %q, got %q", test.Want, result) // HL
			}
		})
	}
}
