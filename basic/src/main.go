package main

import (
	"basic/database"
	"basic/handlers"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "randomstatestring"
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

func login(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	// Validate the state parameter
	if r.FormValue("state") != oauthStateString {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get the authorization code from the callback
	code := r.FormValue("code")

	// Exchange the authorization code for an OAuth 2.0 token
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("couldn't exchange code: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a new HTTP client using the OAuth 2.0 token
	client := googleOauthConfig.Client(oauth2.NoContext, token)

	// Fetch the user's information from Google's API
	userInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println("couldn't get user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer userInfoResponse.Body.Close()

	// Parse the user's information
	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResponse.Body).Decode(&userInfo); err != nil {
		log.Println("couldn't parse user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a new user struct and save the user information to the database
	// user := model.User{
	// 	ID:         userInfo.ID,
	// 	Name:       userInfo.Name,
	// 	Email:      userInfo.Email,
	// 	OAuthToken: token.AccessToken,
	// }
	// if err := database.CreateUser(user); err != nil {
	// 	log.Println("couldn't save user to the database: ", err)
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// Display the user's information on the page
	fmt.Fprintf(w, "Welcome, %s! Your email is %s", userInfo.Name, userInfo.Email)
}

func main() {
	database.InitializeDb()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/login", login)
	router.HandleFunc("/auth/google/callback", callback)

	// Serve static files (e.g., HTMX if used locally)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("/items", handlers.GetItems).Methods("GET")
	router.HandleFunc("/createItem", handlers.CreateItem).Methods("POST")
	router.HandleFunc("/deleteItem/{id}", handlers.DeleteItem).Methods("DELETE")
	router.HandleFunc("/udpateItem/{id}", handlers.UpdateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
