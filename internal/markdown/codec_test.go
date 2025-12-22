package markdown

import (
	"reflect"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	modified := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)

	orig := Note{
		ID:       "note-123",
		Title:    "Hello World",
		Body:     "Line 1\n\nLine 2",
		Tags:     []string{"foo", "bar"},
		Created:  created,
		Modified: modified,
		Links:    []string{"alpha", "beta"},
	}

	data, err := Write(orig)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	parsed, err := Read(data)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if !reflect.DeepEqual(orig, parsed) {
		t.Fatalf("round trip mismatch:\norig  = %+v\nparsed= %+v", orig, parsed)
	}
}

func TestReadRejectsMalformedFrontmatter(t *testing.T) {
	input := []byte(`---
id: bad
title: oops
`)
	if _, err := Read(input); err == nil {
		t.Fatalf("Read() error = nil, want error for malformed frontmatter")
	}
}

func TestWriteAddsFrontmatter(t *testing.T) {
	n := Note{ID: "x", Title: "t", Body: "body"}
	data, err := Write(n)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	out := string(data)
	if len(out) < 4 || out[:4] != "---\n" {
		t.Fatalf("output does not start with frontmatter, got: %q", out)
	}
	if len(out) == 0 || out[len(out)-1] != '\n' {
		t.Fatalf("output should end with newline")
	}
}

func TestHandlesCRLF(t *testing.T) {
	input := "---\r\nid: test\r\ntitle: Test\r\n---\r\n\r\ncontent"
	parsed, err := Read([]byte(input))
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	
	if parsed.ID != "test" {
		t.Fatalf("ID = %q, want %q", parsed.ID, "test")
	}
	if parsed.Title != "Test" {
		t.Fatalf("Title = %q, want %q", parsed.Title, "Test")
	}
	if parsed.Body != "content" {
		t.Fatalf("Body = %q, want %q", parsed.Body, "content")
	}
}
