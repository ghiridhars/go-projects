package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// get all items
func getItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// create an item
func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newItem.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO items (name) VALUES ($1) RETURNING id", newItem.Name).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newItem.ID = fmt.Sprint("%d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}
