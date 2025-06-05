package controllers

import (
	"goforum/models"
	"html/template"
	"strconv"
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
		"slice": func(items interface{}, start, end int) interface{} {
			switch v := items.(type) {
			case []models.Thread:
				if start >= len(v) {
					return []models.Thread{}
				}
				if end > len(v) {
					end = len(v)
				}
				return v[start:end]
			case []models.Message:
				if start >= len(v) {
					return []models.Message{}
				}
				if end > len(v) {
					end = len(v)
				}
				return v[start:end]
			default:
				return items
			}
		},
		"isUserOnline": func(user *models.User) bool {
			if user == nil || !user.ShowOnline {
				return false
			}
			// Consider user online if last activity was within 5 minutes
			threshold := time.Now().Add(-5 * time.Minute)
			return user.LastActivity.After(threshold)
		},
		"getUserInitial": func(user *models.User) string {
			if user == nil || len(user.Username) == 0 {
				return "U"
			}
			return strings.ToUpper(string(user.Username[0]))
		},
		"getAvatarStyle": func(user *models.User) string {
			if user == nil || user.AvatarStyle == "" {
				return "default"
			}
			return user.AvatarStyle
		},
		"getAvatarClass": func(style string) string {
			switch style {
			case "red":
				return "avatar-red"
			case "blue":
				return "avatar-blue"
			case "green":
				return "avatar-green"
			case "purple":
				return "avatar-purple"
			case "orange":
				return "avatar-orange"
			case "pink":
				return "avatar-pink"
			case "teal":
				return "avatar-teal"
			case "admin":
				return "avatar-admin"
			default:
				return "avatar-default"
			}
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
	// FIX: Convert int to string properly
	return strconv.Itoa(count) + " " + unit + "s ago"
}
