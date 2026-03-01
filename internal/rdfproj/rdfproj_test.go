package rdfproj

import (
	"strings"
	"testing"
	"time"

	"github.com/DeDude/weave2/internal/links"
	"github.com/DeDude/weave2/internal/markdown"
)

func TestNoteToTriples_Basic(t *testing.T) {
	note := markdown.Note{
		ID:    "test-note-20250101000000",
		Title: "Test Note",
		Body:  "Test body content",
		Type:  "Note",
	}

	triples := NoteToTriples(note, "http://example.org")

	if len(triples) < 5 {
		t.Fatalf("Expected at least 5 triples, got %d", len(triples))
	}

	hasType := false
	hasID := false
	hasTitle := false
	hasBody := false

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "rdf-syntax-ns#type") {
			hasType = true
		}
		if strings.Contains(s, "identifier") {
			hasID = true
		}
		if strings.Contains(s, "dc/terms/title") {
			hasTitle = true
		}
		if strings.Contains(s, "schema.org/text") {
			hasBody = true
		}
	}

	if !hasType || !hasID || !hasTitle || !hasBody {
		t.Errorf("Missing expected triples: hasType=%v, hasID=%v, hasTitle=%v, hasBody=%v", 
			hasType, hasID, hasTitle, hasBody)
	}
}

func TestNoteToTriples_DualTyping(t *testing.T) {
	note := markdown.Note{
		ID:    "test-20250101000000",
		Title: "Test",
		Type:  "Note",
	}

	triples := NoteToTriples(note, "http://example.org")

	typeCount := 0
	hasWeaveNote := false
	hasSchemaArticle := false

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "rdf-syntax-ns#type") {
			typeCount++
			if strings.Contains(s, "weave.dev/vocab#Note") {
				hasWeaveNote = true
			}
			if strings.Contains(s, "schema.org/Article") {
				hasSchemaArticle = true
			}
		}
	}

	if typeCount < 2 {
		t.Errorf("Expected at least 2 type triples, got %d", typeCount)
	}
	if !hasWeaveNote || !hasSchemaArticle {
		t.Errorf("Missing dual types: hasWeaveNote=%v, hasSchemaArticle=%v", hasWeaveNote, hasSchemaArticle)
	}
}

func TestNoteToTriples_WithDates(t *testing.T) {
	created := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	modified := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)

	note := markdown.Note{
		ID:       "test-20250101000000",
		Title:    "Test",
		Type:     "Note",
		Created:  created,
		Modified: modified,
	}

	triples := NoteToTriples(note, "http://example.org")

	hasCreated := false
	hasModified := false

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "dc/terms/created") {
			hasCreated = true
		}
		if strings.Contains(s, "dc/terms/modified") {
			hasModified = true
		}
	}

	if !hasCreated || !hasModified {
		t.Errorf("Missing date triples: hasCreated=%v, hasModified=%v", hasCreated, hasModified)
	}
}

func TestNoteToTriples_WithTags(t *testing.T) {
	note := markdown.Note{
		ID:    "test-20250101000000",
		Title: "Test",
		Type:  "Note",
		Tags:  []string{"golang", "rdf"},
	}

	triples := NoteToTriples(note, "http://example.org")

	subjectCount := 0
	conceptCount := 0
	labelCount := 0

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "dc/terms/subject") {
			subjectCount++
		}
		if strings.Contains(s, "skos/core#Concept") {
			conceptCount++
		}
		if strings.Contains(s, "skos/core#prefLabel") {
			labelCount++
		}
	}

	if subjectCount != 2 {
		t.Errorf("Expected 2 subject triples, got %d", subjectCount)
	}
	if conceptCount != 2 {
		t.Errorf("Expected 2 concept type triples, got %d", conceptCount)
	}
	if labelCount != 2 {
		t.Errorf("Expected 2 label triples, got %d", labelCount)
	}
}

func TestNoteToTriples_WithLinks(t *testing.T) {
	note := markdown.Note{
		ID:    "test-20250101000000",
		Title: "Test",
		Type:  "Note",
		Links: []links.Link{
			{ID: "other-note", Type: "related", Label: ""},
			{ID: "another-note", Type: "linksTo", Label: ""},
		},
	}

	triples := NoteToTriples(note, "http://example.org")

	hasRelated := false
	hasLinksTo := false

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "skos/core#related") {
			hasRelated = true
		}
		if strings.Contains(s, "weave.dev/vocab#linksTo") {
			hasLinksTo = true
		}
	}

	if !hasRelated || !hasLinksTo {
		t.Errorf("Missing link types: hasRelated=%v, hasLinksTo=%v", hasRelated, hasLinksTo)
	}
}

func TestMapRelationshipType_Known(t *testing.T) {
	tests := []struct {
		relType string
		want    string
	}{
		{"related", "http://www.w3.org/2004/02/skos/core#related"},
		{"broader", "http://www.w3.org/2004/02/skos/core#broader"},
		{"narrower", "http://www.w3.org/2004/02/skos/core#narrower"},
		{"seeAlso", "http://www.w3.org/2000/01/rdf-schema#seeAlso"},
		{"linksTo", "http://weave.dev/vocab#linksTo"},
	}

	for _, tt := range tests {
		got := mapRelationshipType(tt.relType)
		if got != tt.want {
			t.Errorf("mapRelationshipType(%q) = %q, want %q", tt.relType, got, tt.want)
		}
	}
}

func TestMapRelationshipType_Unknown(t *testing.T) {
	got := mapRelationshipType("customType")
	want := "http://weave.dev/vocab#customType"

	if got != want {
		t.Errorf("mapRelationshipType(customType) = %q, want %q", got, want)
	}
}

func TestVaultToTriples_Deduplication(t *testing.T) {
	notes := []markdown.Note{
		{
			ID:    "note-1-20250101000000",
			Title: "Note 1",
			Type:  "Note",
			Tags:  []string{"golang", "testing"},
		},
		{
			ID:    "note-2-20250101000001",
			Title: "Note 2",
			Type:  "Note",
			Tags:  []string{"golang", "rdf"},
		},
	}

	triples := VaultToTriples(notes, "http://example.org")

	// Count how many times "golang" tag is defined
	golangTypeCount := 0
	golangLabelCount := 0

	for _, triple := range triples {
		s := triple.String()
		if strings.Contains(s, "/tags/golang") {
			if strings.Contains(s, "skos/core#Concept") {
				golangTypeCount++
			}
			if strings.Contains(s, "skos/core#prefLabel") {
				golangLabelCount++
			}
		}
	}

	if golangTypeCount != 1 {
		t.Errorf("Expected 1 type definition for 'golang' tag, got %d", golangTypeCount)
	}
	if golangLabelCount != 1 {
		t.Errorf("Expected 1 label for 'golang' tag, got %d", golangLabelCount)
	}
}

func TestSanitizeBaseURI_Empty(t *testing.T) {
	got := sanitizeBaseURI("")
	want := "http://localhost"

	if got != want {
		t.Errorf("sanitizeBaseURI(\"\") = %q, want %q", got, want)
	}
}

func TestSanitizeBaseURI_Valid(t *testing.T) {
	input := "http://example.org"
	got := sanitizeBaseURI(input)

	if got != input {
		t.Errorf("sanitizeBaseURI(%q) = %q, want %q", input, got, input)
	}
}

func TestNoteToTriples_EmptyBaseURI(t *testing.T) {
	note := markdown.Note{
		ID:    "test-20250101000000",
		Title: "Test",
		Type:  "Note",
	}

	triples := NoteToTriples(note, "")

	if len(triples) == 0 {
		t.Fatal("Expected triples with empty baseURI, got none")
	}

	// Should use default baseURI
	hasLocalhost := false
	for _, triple := range triples {
		if strings.Contains(triple.String(), "localhost") {
			hasLocalhost = true
			break
		}
	}

	if !hasLocalhost {
		t.Error("Expected default baseURI (localhost) to be used")
	}
}
