package model

// Define the Item struct
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Define your Error struct
type MyError struct {
	msg string
}
