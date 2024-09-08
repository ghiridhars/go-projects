package main

import (
	"basic/database"
	"basic/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func homePage(w http.ResponseWriter, r *http.Request) {
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

func main() {
	database.InitializeDb()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	// Serve static files (e.g., HTMX if used locally)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/items", handlers.GetItems).Methods("GET")
	router.HandleFunc("/createItem", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/deleteItem/{id}", handlers.DeleteItem).Methods("DELETE")
	router.HandleFunc("/udpateItem/{id}", handlers.UpdateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
