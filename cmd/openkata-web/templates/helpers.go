package templates

import "strings"

// SplitTags splits a comma-separated tags string into trimmed slices.
func SplitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Split(tags, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// TagClass returns a CSS class based on the tag prefix.
func TagClass(tag string) string {
	prefix, _, _ := strings.Cut(tag, ":")
	switch prefix {
	case "tool":
		return "badge badge-orange"
	case "language":
		return "badge badge-purple"
	default:
		return "badge badge-green"
	}
}

// FirstSentence returns text up to and including the first period.
func FirstSentence(s string) string {
	if i := strings.Index(s, "."); i >= 0 {
		return s[:i+1]
	}
	return s
}

// AfterFirstSentence returns text after the first period.
func AfterFirstSentence(s string) string {
	if i := strings.Index(s, "."); i >= 0 {
		return strings.TrimSpace(s[i+1:])
	}
	return ""
}
