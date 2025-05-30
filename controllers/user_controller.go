package controllers

import (
	"html/template"
	"log"
	"net/http"

	"goforum/middleware"
	"goforum/models"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get user's threads
	threadFilters := models.ThreadFilters{
		AuthorID: user.ID,
		Limit:    20,
		Page:     1,
	}
	userThreads, _, _ := models.GetThreads(threadFilters)

	// Get user's messages
	messageFilters := models.MessageFilters{
		UserID: user.ID,
		Limit:  20,
		Page:   1,
	}
	userMessages, _, _ := models.GetMessagesByUser(user.ID, messageFilters)

	// Get user stats
	threadCount, _ := models.GetThreadCountByUser(user.ID)
	messageCount, _ := models.GetMessageCountByUser(user.ID)

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/user/profile.html"))
	data := map[string]interface{}{
		"Title":        "Profile - u/" + user.Username,
		"Page":         "profile",
		"User":         user,
		"UserThreads":  userThreads,
		"UserMessages": userMessages,
		"ThreadCount":  threadCount,
		"MessageCount": messageCount,
	}

	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tmpl := template.Must(template.New("").Funcs(TemplateFuncMap).ParseFiles("views/layouts/main.html", "views/user/settings.html"))
	data := map[string]interface{}{
		"Title": "Settings - GoForum",
		"Page":  "settings",
		"User":  user,
	}

	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
