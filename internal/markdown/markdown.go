package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Note struct {
	ID       string
	Title    string
	Body     string
	Tags     []string
	Created  time.Time
	Modified time.Time
	Links    []string
}

type frontmatter struct {
	ID       string    `yaml:"id"`
	Title    string    `yaml:"title"`
	Tags     []string  `yaml:"tags,omitempty"`
	Created  time.Time `yaml:"created,omitempty"`
	Modified time.Time `yaml:"modified,omitempty"`
	Links    []string  `yaml:"links,omitempty"`
}

func Write(n Note) ([]byte, error) {
	fm := frontmatter{
		ID:       n.ID,
		Title:    n.Title,
		Tags:     n.Tags,
		Created:  n.Created,
		Modified: n.Modified,
		Links:    n.Links,
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshal frontmatter: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(fmBytes)
	buf.WriteString("---\n")
	if n.Body != "" {
		buf.WriteString("\n")
		buf.WriteString(n.Body)
		if !strings.HasSuffix(n.Body, "\n") {
			buf.WriteString("\n")
		}
	} else {
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}

func Read(data []byte) (Note, error) {
	const delim = "---"
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(content, delim+"\n") {
		return Note{}, errors.New("missing frontmatter")
	}

	idx := strings.Index(content[len(delim)+1:], "\n"+delim+"\n")
	if idx == -1 {
		return Note{}, errors.New("malformed frontmatter: closing --- not found")
	}

	end := len(delim) + 1 + idx
	fmText := content[len(delim)+1 : end+1]
	body := content[end+len(delim)+2:]
	if strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	if strings.HasSuffix(body, "\n") {
		body = body[:len(body)-1]
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(fmText), &fm); err != nil {
		return Note{}, fmt.Errorf("unmarshal frontmatter: %w", err)
	}

	n := Note{
		ID:       fm.ID,
		Title:    fm.Title,
		Tags:     fm.Tags,
		Created:  fm.Created,
		Modified: fm.Modified,
		Links:    fm.Links,
		Body:     body,
	}
	return n, nil
}
