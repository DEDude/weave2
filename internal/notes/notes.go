package notes

import (
	"regexp"
	"strings"
	"time"
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
