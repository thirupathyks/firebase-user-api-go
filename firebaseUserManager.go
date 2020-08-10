package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thirupathyks/firebase-user-api-go/models"
)

var webAPIKey = os.Getenv("FIREBASE_WEB_API_KEY")

var client *auth.Client
var err error

//CreateUserHandler - Adds the user to Firebase
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Input
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(false).
		Password(user.Password).
		DisplayName(user.DisplayName).
		Disabled(false)
	u, err := client.CreateUser(context.Background(), params)
	if err != nil {
		fmt.Fprintf(w, "error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", u)
	fmt.Fprintf(w, "Successfully created user: %s\n", user.DisplayName)
}

//UpdateUserHandler - Updates the user info in Firebase
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Input
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	u, err := client.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		log.Fatalf("error getting user by email %s: %v\n", user.Email, err)
	}
	log.Printf("Successfully fetched user data: %v\n", u)
	params := (&auth.UserToUpdate{}).
		Email(user.Email).
		EmailVerified(user.EmailVerified).
		Password(user.Password).
		DisplayName(user.DisplayName).
		Disabled(user.Disabled)
	u1, err := client.UpdateUser(context.Background(), u.UID, params)
	if err != nil {
		log.Fatalf("error updating user: %v\n", err)
	}
	fmt.Fprintf(w, "Successfully updated user: %s\n", u1)
}

//SignInUserHandler - SignIns the user using REST API
func SignInUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Input
	var signInUser models.SignInUserRequest
	json.NewDecoder(r.Body).Decode(&signInUser)
	signInUser.ReturnSecureToken = false
	//Invoke the Firebase REST API to signInWithPassword
	if len(webAPIKey) > 0 {
		url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + webAPIKey
		jsonStr, err := json.Marshal(signInUser)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.Status == "400 Bad Request" {
			http.Error(w, string(body), 400)
		} else {
			fmt.Fprintf(w, "%s", string(body))
		}
	}
}

func accessServicesSingleApp() (*auth.Client, error) {
	// Initialize default app
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return client, err
}

func main() {
	client, err = accessServicesSingleApp()
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/createuser", CreateUserHandler)
	r.HandleFunc("/updateuser", UpdateUserHandler)
	r.HandleFunc("/signinuser", SignInUserHandler)
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":8081", handlers.CORS(headers, methods, origins)(r))
}
