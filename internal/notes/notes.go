package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/DeDude/weave2/internal/markdown"
)

func GenerateID(title string, timestamp time.Time) string {
	slug := slugify(title)
	ts := formatTimestamp(timestamp)

	if slug == "" {
		return ts
	}

	return slug + "-" + ts
}

func slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")

	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}

func formatTimestamp(t time.Time) string {
	return t.UTC().Format("20060102150405")
}

func ResolvePath(vaultPath, id string) string {
	timestamp := id[len(id)-14:]

	year := timestamp[0:4]
	month := timestamp[4:6]

	return vaultPath + "/" + year + "/" + month + "/" + id + ".md"
}

func Create(vaultPath string, note markdown.Note, timestamp time.Time) (string, error) {
	id := GenerateID(note.Title, timestamp)

	note.ID = id
	note.Created = timestamp
	note.Modified = timestamp

	filePath := ResolvePath(vaultPath, id)

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}

	data, err := markdown.Write(note)
	if err != nil {
		return "", fmt.Errorf("write markdown: %w", err)
	}

	tempPath := filePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath)

		return "", fmt.Errorf("rename file: %w", err)
	}

	return id, nil
}

func Read(vaultPath, id string) (markdown.Note, error) {
	filePath := ResolvePath(vaultPath, id)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return markdown.Note{}, fmt.Errorf("read file: %w", err)
	}

	note, err := markdown.Read(data)
	if err != nil {
		return markdown.Note{}, fmt.Errorf("parse markdownL %w", err)
	}

	return note, nil
}
