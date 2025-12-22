package notes

import (
	"testing"
	"time"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My Note Title", "my-note-title"},
		{"Hello World!", "hello-world"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"Special@#$Chars", "specialchars"},
		{"   Leading and Trailing   ", "leading-and-trailing"},
		{"CamelCaseTitle", "camelcasetitle"},
		{"123 Numbers", "123-numbers"},
		{"", ""},
		{"---", ""},
		{"a", "a"},
	}

	for _, tt := range tests {
		got := slugify(tt.input)

		if got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatTimestamp(t *testing.T) {
	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)
	want := "20250122223045"
	got := formatTimestamp(ts)

	if got != want {
		t.Errorf("formatTimestamp() = %q, want %q", got, want)
	}

	loc, _ := time.LoadLocation("America/New_York")
	tsLocal := time.Date(2025, 1, 22, 17, 30, 45, 0, loc)
	want = "20250122223045"
	got = formatTimestamp(tsLocal)

	if got != want {
		t.Errorf("formatTimestamp() with local time = %q, want %q", got, want)
	}
}

func TestGenerateID(t *testing.T) {
	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)

	tests := []struct {
		title string
		want  string
	}{
		{"My Note Title", "my-note-title-20250122223045"},
		{"Hello World", "hello-world-20250122223045"},
		{"Test", "test-20250122223045"},
		{"", "20250122223045"},
	}

	for _, tt := range tests {
		got := GenerateID(tt.title, ts)

		if got != tt.want {
			t.Errorf("GenerateID(%q, ts) = %q, want %q", tt.title, got, tt.want)
		}
	}
}
