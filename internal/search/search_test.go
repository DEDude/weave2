package search

import (
	"testing"
	"time"

	"github.com/DeDude/weave2/internal/markdown"
	"github.com/DeDude/weave2/internal/notes"
	"github.com/DeDude/weave2/internal/links"
)

func TestSearchEmpty(t *testing.T) {
	vaultPath := t.TempDir()
	
	results, err := Search(vaultPath, Query{Term: "test"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 0 {
		t.Errorf("Search() returned %d results, want 0", len(results))
	}
}

func TestSearchTitle(t *testing.T) {
	vaultPath := t.TempDir()
	
	note := markdown.Note{
		Title: "Test Note",
		Body:  "Some content",
	}
	
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := notes.Create(vaultPath, note, ts)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	results, err := Search(vaultPath, Query{Term: "test"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(results))
	}
	
	if results[0].Note.Title != "Test Note" {
		t.Errorf("Result title = %q, want %q", results[0].Note.Title, "Test Note")
	}
}

func TestSearchBody(t *testing.T) {
	vaultPath := t.TempDir()
	
	note := markdown.Note{
		Title: "My Note",
		Body:  "This contains the keyword test",
	}
	
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := notes.Create(vaultPath, note, ts)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	results, err := Search(vaultPath, Query{Term: "keyword"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(results))
	}
}

func TestSearchTags(t *testing.T) {
	vaultPath := t.TempDir()
	
	note := markdown.Note{
		Title: "My Note",
		Body:  "Content",
		Tags:  []string{"golang", "testing"},
	}
	
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := notes.Create(vaultPath, note, ts)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	results, err := Search(vaultPath, Query{Term: "golang"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	vaultPath := t.TempDir()
	
	note := markdown.Note{
		Title: "Test Note",
		Body:  "Content",
	}
	
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := notes.Create(vaultPath, note, ts)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	results, err := Search(vaultPath, Query{Term: "TEST"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(results))
	}
}

func TestSearchScoring(t *testing.T) {
	vaultPath := t.TempDir()
	
	note1 := markdown.Note{
		Title: "Test",
		Body:  "Content",
	}
	
	note2 := markdown.Note{
		Title: "Test test test",
		Body:  "More test content",
	}
	
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := notes.Create(vaultPath, note1, ts)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	ts2 := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	_, err = notes.Create(vaultPath, note2, ts2)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	
	results, err := Search(vaultPath, Query{Term: "test"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	
	if len(results) != 2 {
		t.Fatalf("Search() returned %d results, want 2", len(results))
	}
	
	if results[0].Score <= results[1].Score {
		t.Errorf("Results not sorted by score: [0].Score=%d, [1].Score=%d", results[0].Score, results[1].Score)
	}
}

func TestSearchLinks(t *testing.T) {
	vaultPath := t.TempDir()

	note := markdown.Note{
		Title: "My Note",
		Body: "Content",
		Links: []links.Link{
			{ID: "target-note=123", Type: "linksTo", Label: ""},
		},
	}

	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := notes.Create(vaultPath, note, ts)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	results, err := Search(vaultPath, Query{Term: "target"})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(results))
	}
}
