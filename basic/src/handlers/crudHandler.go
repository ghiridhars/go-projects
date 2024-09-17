package handlers

import (
	"basic/database"
	"basic/model"
	"encoding/json"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

// get all items
func GetItems(w http.ResponseWriter, r *http.Request) {
	slog.Info("getitems---------------START")
	var items []model.Item = database.GetItemsDAO()
	slog.Info("getitems---------------END")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// create an item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateItem---------------START")
	var newItem model.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newItem.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	successIns := database.CreateItemDAO(newItem)

	slog.Info("CreateItem---------------END")

	// newItem.ID = fmt.Sprint("%d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successIns)
}

// Update an item by ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	slog.Info("UpdateItem---------------START")
	params := mux.Vars(r)
	var updatedItem model.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedItem.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	rowsAffected := database.UpadteItemDAO(updatedItem, params["id"])
	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	updatedItem.ID = params["id"]
	slog.Info("UpdateItem---------------END")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	slog.Info("DeleteItem---------------START")
	params := mux.Vars(r)

	rowsAffected := database.DeleteItemDAO(params["id"])

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	slog.Info("DeleteItem---------------END")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
