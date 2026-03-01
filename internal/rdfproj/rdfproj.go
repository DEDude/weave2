package rdfproj

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DeDude/weave2/internal/markdown"
	rdf "github.com/deiu/rdf2go"
)

// Vocabulary URIs
const (
	rdfType         = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	dctermsID       = "http://purl.org/dc/terms/identifier"
	dctermsTitle    = "http://purl.org/dc/terms/title"
	dctermsCreated  = "http://purl.org/dc/terms/created"
	dctermsModified = "http://purl.org/dc/terms/modified"
	dctermsSubject  = "http://purl.org/dc/terms/subject"
	schemaText      = "http://schema.org/text"
	skosConcept     = "http://www.w3.org/2004/02/skos/core#Concept"
	skosPrefLabel   = "http://www.w3.org/2004/02/skos/core#prefLabel"
	xsdDateTime     = "http://www.w3.org/2001/XMLSchema#dateTime"

	defaultBaseURI = "http://localhost"
)

var relationshipPredicates = map[string]string{
	"linksTo":  "http://weave.dev/vocab#linksTo",
	"related":  "http://www.w3.org/2004/02/skos/core#related",
	"broader":  "http://www.w3.org/2004/02/skos/core#broader",
	"narrower": "http://www.w3.org/2004/02/skos/core#narrower",
	"seeAlso":  "http://www.w3.org/2000/01/rdf-schema#seeAlso",
}

var noteTypeMap = map[string][]string{
	"Note": {
		"http://weave.dev/vocab#Note",
		"http://schema.org/Article",
	},
}

func NoteToTriples(note markdown.Note, baseURI string) []*rdf.Triple {
	baseURI = sanitizeBaseURI(baseURI)
	var triples []*rdf.Triple
	noteURI := makeNoteURI(baseURI, note.ID)

	// rdf:type (dual typing)
	types := noteTypeMap[note.Type]
	if types == nil {
		types = []string{"http://weave.dev/vocab#" + note.Type}
	}
	for _, t := range types {
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(rdfType),
			rdf.NewResource(t),
		))
	}

	// dcterms:identifier
	triples = append(triples, rdf.NewTriple(
		rdf.NewResource(noteURI),
		rdf.NewResource(dctermsID),
		rdf.NewLiteral(note.ID),
	))

	// dcterms:title
	triples = append(triples, rdf.NewTriple(
		rdf.NewResource(noteURI),
		rdf.NewResource(dctermsTitle),
		rdf.NewLiteral(note.Title),
	))

	// schema:text (body)
	if note.Body != "" {
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(schemaText),
			rdf.NewLiteral(note.Body),
		))
	}

	// dcterms:created
	if !note.Created.IsZero() {
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(dctermsCreated),
			rdf.NewLiteralWithDatatype(
				note.Created.Format("2006-01-02T15:04:05Z07:00"),
				rdf.NewResource(xsdDateTime),
			),
		))
	}

	// dcterms:modified
	if !note.Modified.IsZero() {
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(dctermsModified),
			rdf.NewLiteralWithDatatype(
				note.Modified.Format("2006-01-02T15:04:05Z07:00"),
				rdf.NewResource(xsdDateTime),
			),
		))
	}

	// Tags
	for _, tag := range note.Tags {
		tagURI := makeTagURI(baseURI, tag)

		// dcterms:subject
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(dctermsSubject),
			rdf.NewResource(tagURI),
		))

		// rdf:type skos:Concept
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(tagURI),
			rdf.NewResource(rdfType),
			rdf.NewResource(skosConcept),
		))

		// skos:prefLabel
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(tagURI),
			rdf.NewResource(skosPrefLabel),
			rdf.NewLiteral(tag),
		))
	}

	// Links
	for _, link := range note.Links {
		predicate := mapRelationshipType(link.Type)
		targetURI := makeNoteURI(baseURI, link.ID)
		triples = append(triples, rdf.NewTriple(
			rdf.NewResource(noteURI),
			rdf.NewResource(predicate),
			rdf.NewResource(targetURI),
		))
	}

	return triples
}

func VaultToTriples(notes []markdown.Note, baseURI string) []*rdf.Triple {
	baseURI = sanitizeBaseURI(baseURI)
	var triples []*rdf.Triple
	seenTags := make(map[string]bool)

	for _, note := range notes {
		noteTriples := NoteToTriples(note, baseURI)

		for _, triple := range noteTriples {
			if isTagDefinitionTriple(triple, baseURI) {
				key := triple.Subject.String() + triple.Predicate.String()
				if !seenTags[key] {
					seenTags[key] = true
					triples = append(triples, triple)
				}
			} else {
				triples = append(triples, triple)
			}
		}
	}

	return triples
}

func mapRelationshipType(relType string) string {
	if pred, ok := relationshipPredicates[relType]; ok {
		return pred
	}
	return fmt.Sprintf("http://weave.dev/vocab#%s", relType)
}

func makeNoteURI(baseURI, id string) string {
	return fmt.Sprintf("%s/notes/%s", baseURI, url.PathEscape(id))
}

func makeTagURI(baseURI, tag string) string {
	return fmt.Sprintf("%s/tags/%s", baseURI, url.PathEscape(tag))
}

func sanitizeBaseURI(baseURI string) string {
	if baseURI == "" {
		return defaultBaseURI
	}
	return baseURI
}

func isTagDefinitionTriple(triple *rdf.Triple, baseURI string) bool {
	subj := triple.Subject.String()
	pred := triple.Predicate.String()

	isTagSubject := strings.Contains(subj, "/tags/")
	isTagPredicate := pred == "<"+rdfType+">" || pred == "<"+skosPrefLabel+">"

	return isTagSubject && isTagPredicate
}
