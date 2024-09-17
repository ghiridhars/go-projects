package model

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OAuthToken string `json:"token"`
	Email      string `json:"email"`
}

// Define your Error struct
type MyError struct {
	msg string
}
