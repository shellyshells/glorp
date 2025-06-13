package utils

import (
	"html/template"
	"log"
	"net/http"
)

// ErrorPageData represents the data structure for the error page
type ErrorPageData struct {
	ErrorCode    int
	ErrorTitle   string
	ErrorMessage string
	ErrorDetails string
}

// ShowErrorPage displays a custom error page
func ShowErrorPage(w http.ResponseWriter, r *http.Request, code int, message string, details string) {
	// Set the status code
	w.WriteHeader(code)

	// Get the error title based on the status code
	var title string
	switch code {
	case http.StatusNotFound:
		title = "Page Not Found"
	case http.StatusForbidden:
		title = "Access Denied"
	case http.StatusInternalServerError:
		title = "Server Error"
	default:
		title = "Error"
	}

	// Create the error page data
	data := ErrorPageData{
		ErrorCode:    code,
		ErrorTitle:   title,
		ErrorMessage: message,
		ErrorDetails: details,
	}

	// Parse and execute the error template
	tmpl := template.Must(template.New("").ParseFiles("views/layouts/main.html", "views/errors/error.html"))
	if err := tmpl.ExecuteTemplate(w, "main.html", data); err != nil {
		log.Printf("Error executing error template: %v", err)
		// If template execution fails, fall back to a simple error message
		http.Error(w, message, code)
	}
}
