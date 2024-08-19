package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "postgres://youruser:yourpassword@localhost:5432/yourdb?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// create items if not there
	query := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my first GO page")
}

func getitems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		items = append(items, item)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func main() {
	// initDB()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/items", getitems).Methods("GET")
	router.HandleFunc("/items", createItem).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
