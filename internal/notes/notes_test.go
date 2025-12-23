package notes

import (
	"os"
	"testing"
	"time"

	"github.com/DeDude/weave2/internal/markdown"
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

func TestResolvePath(t *testing.T) {
	tests := []struct {
		vaultPath string
		id        string
		want      string
	}{
		{"/vault", "my-note-20250122223045", "/vault/2025/01/my-note-20250122223045.md"},
		{"/vault", "test-20251231235959", "/vault/2025/12/test-20251231235959.md"},
		{"/vault", "note-20240101000000", "/vault/2024/01/note-20240101000000.md"},
		{"/home/notes", "hello-20250615120000", "/home/notes/2025/06/hello-20250615120000.md"},
		{"/vault", "20250122223045", "/vault/2025/01/20250122223045.md"},
	}

	for _, tt := range tests {
		got := ResolvePath(tt.vaultPath, tt.id)

		if got != tt.want {
			t.Errorf("ResolvePath(%q, %q) = %q, want %q", tt.vaultPath, tt.id, got, tt.want)
		}
	}
}

func TestCreate(t *testing.T) {
	vaultPath := t.TempDir()

	note := markdown.Note{
		Title: "My Test Note",
		Body:  "This is the body",
		Tags:  []string{"test"},
	}

	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)

	id, err := Create(vaultPath, note, ts)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	expectedID := "my-test-note-20250122223045"

	if id != expectedID {
		t.Errorf("Create() id = %q, want %q", id, expectedID)
	}

	expectedPath := vaultPath + "/2025/01/my-test-note-20250122223045.md"

	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("File not created at %q", expectedPath)
	}

	data, err := os.ReadFile(expectedPath)

	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	parsed, err := markdown.Read(data)

	if err != nil {
		t.Fatalf("markdown.Read() error = %v", err)
	}

	if parsed.ID != expectedID {
		t.Errorf("parsed.ID = %q, want %q", parsed.ID, expectedID)
	}

	if parsed.Title != "My Test Note" {
		t.Errorf("parsed.Title = %q, want %q", parsed.Title, "My Test Note")
	}

	if parsed.Body != "This is the body" {
		t.Errorf("parsed.Body = %q, want %q", parsed.Body, "This is the body")
	}
}

func TestRead(t *testing.T) {
	vaultPath := t.TempDir()

	note := markdown.Note{
		Title: "Test Note",
		Body:  "Test body",
		Tags:  []string{"tag1", "tag2"},
	}

	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)
	id, err := Create(vaultPath, note, ts)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	loaded, err := Read(vaultPath, id)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if loaded.ID != id {
		t.Errorf("loaded.ID = %q, want %q", loaded.ID, id)
	}

	if loaded.Title != "Test Note" {
		t.Errorf("loaded.Title = %q, want %q", loaded.Title, "Test Note")
	}

	if loaded.Body != "Test body" {
		t.Errorf("loaded.Body = %q, want %q", loaded.Body, "Test body")
	}

	if len(loaded.Tags) != 2 {
		t.Errorf("len(loaded.Tags) = %d, want 2", len(loaded.Tags))
	}
}

func TestReadNonExistent(t *testing.T) {
	vaultPath := t.TempDir()

	_, err := Read(vaultPath, "nonexistent-20250122223045")

	if err == nil {
		t.Fatal("Read() error = nil, want error for non-existent note")
	}
}

func TestUpdate(t *testing.T) {
	vaultPath := t.TempDir()

	note := markdown.Note{
		Title: "Original Title",
		Body:  "original body",
		Tags:  []string{"tag1"},
	}

	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)
	id, err := Create(vaultPath, note, ts)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	note.Title = "Updated Title"
	note.Body = "Updated body"
	note.Tags = []string{"tag1", "tag2"}

	updateTime := time.Date(2025, 1, 22, 23, 0, 0, 0, time.UTC)
	err = Update(vaultPath, id, note, updateTime)

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	loaded, err := Read(vaultPath, id)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if loaded.Title != "Updated Title" {
		t.Errorf("loaded.Title = %q, want %q", loaded.Title, "Updated Title")
	}

	if loaded.Body != "Updated body" {
		t.Errorf("loaded.Body = %q, want %q", loaded.Body, "Updated body")
	}

	if loaded.Modified != updateTime {
		t.Errorf("loaded.Modified = %v, want %v", loaded.Modified, updateTime)
	}

	if loaded.Created != ts {
		t.Errorf("loaded.Created = %v, want %v (should not change)", loaded.Created, ts)
	}
}

func TestDelete(t *testing.T) {
	vaultPath := t.TempDir()

	note := markdown.Note{
		Title: "Test Note",
		Body:  "Test body",
	}

	ts := time.Date(2025, 1, 22, 22, 30, 45, 0, time.UTC)
	id, err := Create(vaultPath, note, ts)

	if err != nil {
		t.Fatalf("Create() erorr = %v", err)
	}

	filePath := ResolvePath(vaultPath, id)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("File should exist before delete")
	}

	err = Delete(vaultPath, id)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("File should not exist after delete")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	vaultPath := t.TempDir()

	err := Delete(vaultPath, "nonexistent-20250122223045")
	if err == nil {
		t.Fatal("Delete() error = nil, want error for non-existent note")
	}
}

func TestList(t *testing.T) {
	vaultPath := t.TempDir()

	notes := []markdown.Note{
		{Title: "Note one", Body: "Body 1"},
		{Title: "Note two", Body: "Body 2"},
		{Title: "Note three", Body: "Body 3"},
	}

	ts1 := time.Date(2025, 1, 22, 10, 0, 0, 0, time.UTC)
	ts2 := time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC)
	ts3 := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	_, err := Create(vaultPath, notes[0], ts1)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = Create(vaultPath, notes[1], ts2)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = Create(vaultPath, notes[2], ts3)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	loaded, err := List(vaultPath)
	
	if err != nil {
		t.Fatalf("List() error = %v", err)	
	}

	if len(loaded) != 3 {
		t.Fatalf("List() returned %d notes, want 3", len(loaded))
	}

	titles := make(map[string]bool)
	for _, note := range loaded {
		titles[note.Title] = true
	}

	if !titles["Note one"] || !titles["Note two"] || !titles["Note three"] {
		t.Error("Not all notes were loaded")
	}
}

func TestListEmpty(t *testing.T){
	vaultPath := t.TempDir()

	loaded, err := List(vaultPath)
	
	if err != nil {
		t.Fatalf("List() error = %v", err)	
	}

	if len(loaded) != 0 {
		t.Errorf("List() returned %d notes, want 0", len(loaded))
	}
}
