package controllers

import (
	"html/template"
	"strings"
	"time"
)

var TemplateFuncMap template.FuncMap

func init() {
	TemplateFuncMap = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"timeAgo": func(t time.Time) string {
			return timeAgo(t)
		},
		"truncate": func(text string, length int) string {
			if len(text) <= length {
				return text
			}
			return text[:length] + "..."
		},
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("Jan 02, 2006 15:04")
		},
		"pluralize": func(count int, singular, plural string) string {
			if count == 1 {
				return singular
			}
			return plural
		},
		"join": func(elements []string, separator string) string {
			return strings.Join(elements, separator)
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
	}
}

func timeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return formatDuration(minutes, "minute")
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return formatDuration(hours, "hour")
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		return formatDuration(days, "day")
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		return formatDuration(months, "month")
	} else {
		years := int(diff.Hours() / (24 * 365))
		return formatDuration(years, "year")
	}
}

func formatDuration(count int, unit string) string {
	if count == 1 {
		return "1 " + unit + " ago"
	}
	return string(rune(count)) + " " + unit + "s ago"
}
