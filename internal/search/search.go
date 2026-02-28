package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/DeDude/weave2/internal/markdown"
	"github.com/DeDude/weave2/internal/notes"
)

type Query struct {
	Term string
}

type Result struct {
	Note  markdown.Note
	Score int
}

func Search(vaultPath string, query Query) ([]Result, error) {
	allNotes, errs := notes.List(vaultPath)
	if len(errs) > 0 {
		return nil, fmt.Errorf("search failed with %d errors: %w", len(errs), errs[0])
	}

	var results []Result
	term := strings.ToLower(query.Term)

	for _, note := range allNotes {
		score := scoreNote(note, term)
		if score > 0 {
			results = append(results, Result{
				Note:  note,
				Score: score,
			})
		}
	}

	sortResults(results)
	return results, nil
}

func scoreNote(note markdown.Note, term string) int {
	score := 0
	
	score += strings.Count(strings.ToLower(note.Title), term)
	score += strings.Count(strings.ToLower(note.Body), term)
	
	for _, tag := range note.Tags {
		score += strings.Count(strings.ToLower(tag), term)
	}
	
	for _, link := range note.Links {
		score += strings.Count(strings.ToLower(link.ID), term)
	}
	
	return score
}

func sortResults(results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}
