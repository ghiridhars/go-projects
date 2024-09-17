package handlers

import (
	"basic/database"
	"basic/model"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
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

func Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("Login---------------START")
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	slog.Info("Login---------------END")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	slog.Info("Callback---------------START")

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
	fmt.Printf("%+v\n", userInfo)
	slog.Info("user :: %+v", userInfo)
	if err := json.NewDecoder(userInfoResponse.Body).Decode(&userInfo); err != nil {
		log.Println("couldn't parse user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a new user struct and save the user information to the database
	user := model.User{
		ID:         userInfo.ID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		OAuthToken: token.AccessToken,
	}
	userExists := database.GetUserByMail(user.Email)

	if userExists.Email != "" {
		slog.Error("Already in User table. Not Inserting")
	} else if err := database.CreateUserDAO(user); err != "Successfully Inserted" {
		log.Fatal("couldn't save user to the database: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Display the user's information on the page
	fmt.Fprintf(w, "Welcome, %s! Your email is %s", userInfo.Name, userInfo.Email)
	slog.Info("Callback---------------END")
}
