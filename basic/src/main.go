package main

import (
	"basic/database"
	"basic/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	database.InitializeDb()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handlers.HomePage)
	// router.HandleFunc("/toggle-theme", handlers.ToggleButtonTheme)

	router.HandleFunc("/login", handlers.Login)
	router.HandleFunc("/auth/google/callback", handlers.Callback)

	// Serve static files (e.g., HTMX if used locally)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/items", handlers.GetItems).Methods("GET")
	router.HandleFunc("/createItem", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/deleteItem/{id}", handlers.DeleteItem).Methods("DELETE")
	router.HandleFunc("/udpateItem/{id}", handlers.UpdateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
