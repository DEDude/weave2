package links

import "strings"

type Link struct {
	ID    string
	Type  string
	Label string
}

const DefaultLinkType = "linksTo"

func FormatLink(id, relType, label string) string {
	if relType == "" {
		relType = DefaultLinkType
	}

	var core string
	if relType == DefaultLinkType {
		core = id
	} else {
		core = relType + "::" + id
	}

	if label != "" {
		core = core + "|" + label
	}

	return "[[" + core + "]]"
}

func ParseLinks(body string) []Link {
	var out []Link
	start := 0

	for {
		open := strings.Index(body[start:], "[[")
		if open == -1 {
			break
		}
		open += start
		
		close := strings.Index(body[open+2:], "]]")
		if close == -1 {
			break
		}
		close += open + 2

		content := body[open+2 : close]
		if content == "" {
			start = close + 2
			continue
		}
		if strings.Contains(content, "[[") {
			start = close + 2
			continue
		}

		link := parseLinkContent(content)
		if link.ID != "" {
			out = append(out, link)
		}

		start = close + 2
	}

	return out
}

func parseLinkContent(content string) Link {
	parts := strings.SplitN(content, "|", 2)
	left := parts[0]

	label := ""
	if len(parts) == 2 {
		label = parts[1]
	}

	typeAndID := strings.SplitN(left, "::", 2)
	if len(typeAndID) == 2 {
		relType := typeAndID[0]
		id := typeAndID[1]
		if relType == "" || id == "" {
			return Link{}
		}
		return Link{ID: id, Type: relType, Label: label}
	}

	if left == "" {
		return Link{}
	}

	return Link{ID: left, Type: DefaultLinkType, Label: label}
}
