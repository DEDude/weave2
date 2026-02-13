package links

import (
	"reflect"
	"testing"
)

func TestFormatLink_DefaultTypeNoLabel(t *testing.T) {
	got  := FormatLink("note-123", "", "")
	want := "[[note-123]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestFormatLink_ExplicitDefaultTypeNoLabel(t *testing.T) {
	got  := FormatLink("note-123", DefaultLinkType, "")
	want := "[[note-123]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestFormatLink_DefaultTypeWithLabel(t *testing.T) {
	got  := FormatLink("note-123", "", "My Note")
	want := "[[note-123|My Note]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestFormatLink_ExplicitDefaultTypeWithLabel(t *testing.T) {
	got  := FormatLink("note-123", DefaultLinkType, "My Note")
	want := "[[note-123|My Note]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestFormatLink_TypedNoLabel(t *testing.T) {
	got  := FormatLink("note-123", "related", "")
	want := "[[related::note-123]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestFormatLink_TypedWithLabel(t *testing.T) {
	got  := FormatLink("note-123", "related", "Reference")
	want := "[[related::note-123|Reference]]"

	if got != want {
		t.Fatalf("FormatLink() = %q, want %q", got, want)
	}
}

func TestParseLinks_Empty(t *testing.T) {
	got := ParseLinks("")
	if len(got) != 0 {
		t.Fatalf("len(ParseLinks()) = %d, want 0", len(got))
	}
}

func TestParseLinks_DefaultNoLabel(t *testing.T) {
	body := "See [[note-123]] for details."
	got  := ParseLinks(body)

	want := []Link{
		{ID: "note-123", Type: DefaultLinkType, Label: ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	} 
}

func TestParseLinks_DefaultWithLabel(t *testing.T) {
	body := "See [[note-123|My Note]] for details."
	got  := ParseLinks(body)

	want := []Link{
		{ID: "note-123", Type: DefaultLinkType, Label: "My Note"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	}
}

func TestParseLinks_TypedNoLabel(t *testing.T) {
	body := "Related: [[related::note-123]]"
	got := ParseLinks(body)

	want := []Link{
		{ID: "note-123", Type: "related", Label: ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	}
}

func TestParseLinks_TypedWithLabel(t *testing.T) {
	body := "Related: [[related::note-123|Reference]]"
	got := ParseLinks(body)

	want := []Link{
		{ID: "note-123", Type: "related", Label: "Reference"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	}
}

func TestParseLinks_Multiple(t *testing.T) {
	body := "A [[a-1]] B [[related::b-2|Bee]] C [[c-3|See]]"
	got  := ParseLinks(body)

	want := []Link{
		{ID: "a-1", Type: DefaultLinkType, Label: ""},
		{ID: "b-2", Type: "related", Label: "Bee"},
		{ID: "c-3", Type: DefaultLinkType, Label: "See"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	}
}

func TestParseLinks_SkipsMalformed(t *testing.T) {
	body := "Bad [[no-close and [[::missingid]] and [[ok-1]]"
	got  := ParseLinks(body)

	want := []Link{
		{ID: "ok-1", Type: DefaultLinkType, Label: ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseLinks() = %#v, want %#v", got, want)
	}
}
