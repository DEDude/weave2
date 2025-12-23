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

var (
	slugSpecialChars = regexp.MustCompile("[^a-z0-9-]+")
	slugMultiHyphens = regexp.MustCompile("-+")
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
	s = slugSpecialChars.ReplaceAllString(s, "")
	s = slugMultiHyphens.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	return s
}

func formatTimestamp(t time.Time) string {
	return t.UTC().Format("20060102150405")
}

func ResolvePath(vaultPath, id string) (string, error) {
	if len(id) < 14 {
		return "", fmt.Errorf("invalid ID: must be at least 14 characters, got %d", len(id))
	}

	timestamp := id[len(id)-14:]

	year := timestamp[0:4]
	month := timestamp[4:6]

	return vaultPath + "/" + year + "/" + month + "/" + id + ".md", nil
}

func Create(vaultPath string, note markdown.Note, timestamp time.Time) (string, error) {
	id := GenerateID(note.Title, timestamp)

	note.ID = id
	note.Created = timestamp
	note.Modified = timestamp

	filePath, err := ResolvePath(vaultPath, id)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}

	data, err := markdown.Write(note)
	if err != nil {
		return "", fmt.Errorf("write markdown: %w", err)
	}

	if err := safeWrite(filePath, data); err != nil {
		return "", err
	}

	return id, nil
}

func Read(vaultPath, id string) (markdown.Note, error) {
	filePath, err := ResolvePath(vaultPath, id)
	if err != nil {
		return markdown.Note{}, fmt.Errorf("resolve path: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return markdown.Note{}, fmt.Errorf("read file: %w", err)
	}

	note, err := markdown.Read(data)
	if err != nil {
		return markdown.Note{}, fmt.Errorf("parse markdown: %w", err)
	}

	return note, nil
}

func Update(vaultPath, id string, note markdown.Note, timestamp time.Time) error {
	existing, err := Read(vaultPath, id)
	if err != nil {
		return fmt.Errorf("read existing note: %w", err)
	}

	note.ID = id
	note.Created = existing.Created
	note.Modified = timestamp

	filePath, err := ResolvePath(vaultPath, id)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	data, err := markdown.Write(note)

	if err != nil {
		return fmt.Errorf("write markdown: %w", err)
	}

	if err := safeWrite(filePath, data); err != nil {
		return err
	}

	return nil
}

func Delete(vaultPath, id string) error {
	filePath, err := ResolvePath(vaultPath, id)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("remove file: %w", err)
	}

	return nil
}

func safeWrite(filePath string, data []byte) error {
	tempPath := filePath + ".tmp"
	
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	
	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("rename file: %w", err)
	}
	
	return nil
}

func List(vaultPath string) ([]markdown.Note, []error) {
	var notes []markdown.Note
	var errors []error

	err := filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", path, err))
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		data, err := os.ReadFile(path)

		if err != nil {
			errors = append(errors, fmt.Errorf("%s: read failed: %w", path, err))
			return nil
		}

		note, err := markdown.Read(data)

		if err != nil {
			errors = append(errors, fmt.Errorf("%s: parse failed: %w", path, err))
			return nil
		}

		notes = append(notes, note)
		return nil
	})

	if err != nil {
		errors = append(errors, fmt.Errorf("walk vault: %w", err))
	}

	return notes, errors
}
