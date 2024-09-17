package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// commented out all this server side theme switching and restricted it to html invert to simplify the process

// func ToggleButtonTheme(w http.ResponseWriter, r *http.Request) {
// 	theme := r.Header.Get("X-Theme")
// 	if theme == "" {
// 		theme = "light"
// 	}

// 	// Toggle between light and dark themes
// 	if theme == "dark" {
// 		w.Header().Set("X-Theme", "light")
// 		w.Write([]byte(`<button hx-get="/toggle-theme" hx-trigger="click" hx-swap="outerHTML">Switch to Dark Theme</button>`))
// 	} else {
// 		w.Header().Set("X-Theme", "dark")
// 		w.Write([]byte(`<button hx-get="/toggle-theme" hx-trigger="click" hx-swap="outerHTML">Switch to Light Theme</button>`))
// 	}
// }

// func ToggleTheme(w http.ResponseWriter, r *http.Request) {
// 	// Get the theme from the custom header
// 	theme := r.Header.Get("X-Theme")
// 	slog.Info("::", theme)
// 	if theme == "" {
// 		// Default to light theme if not set
// 		theme = "light"
// 	}
// 	slog.Info("::", theme)

// 	// Toggle between light and dark themes
// 	if theme == "dark" {
// 		slog.Info("will render light theme")
// 		// Respond with light theme
// 		w.Header().Set("X-Theme", "light")
// 		w.Write([]byte(`<i id="theme-toggle-icon" class="fas fa-moon theme-toggle-icon"
// 		hx-get="/toggle-theme" hx-swap="outerHTML" hx-trigger="click"></i>`))
// 	} else {
// 		slog.Info("will render dark theme")
// 		// Respond with dark theme
// 		w.Header().Set("X-Theme", "dark")
// 		w.Write([]byte(`<i id="theme-toggle-icon" class="fas fa-sun theme-toggle-icon"
// 		hx-get="/toggle-theme" hx-swap="outerHTML" hx-trigger="click"></i>`))
// 	}
// }

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Print current working directory
		dir, _ := os.Getwd()
		fmt.Println("Current working directory:", dir)

		// Print the absolute path to the template
		tmplPath, err := filepath.Abs(filepath.Join("templates", "home.html"))
		if err != nil {
			http.Error(w, "Error finding template", http.StatusInternalServerError)
			return
		}
		fmt.Println("Template path:", tmplPath)

		tmpl := template.Must(template.ParseFiles(tmplPath))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		name := r.FormValue("name")
		response := "Hello, " + name + "!"
		w.Write([]byte(response))
	}
}
