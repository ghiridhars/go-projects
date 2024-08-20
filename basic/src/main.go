package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	// "postgres://ghiri:develop@postgres-go:5432/items?sslmode=disable"
	slog.Info(connStr)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("success")

	if err = db.Ping(); err != nil {
		slog.Error("Postgres ping error : (%v)", err)
	}

	// create items if not there
	query := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`
	_, err = db.Exec(query)

	slog.Info("SUccess")
	if err != nil {
		log.Fatal(err)
	}
}

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
	initDB()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	// Serve static files (e.g., HTMX if used locally)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/createItem", createItem).Methods("POST")
	router.HandleFunc("/deleteItem/{id}", deleteItem).Methods("DELETE")
	router.HandleFunc("/udpateItem/{id}", updateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
