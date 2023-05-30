package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MateSousa/cloud-build/initializers"
	"github.com/MateSousa/cloud-build/models"
)

// auth is a struct that holds the configuration for the auth api

func Login(w http.ResponseWriter, r *http.Request) {
	var user *User
	var token models.Token
	DB := initializers.DB
	Redis := initializers.Redis

	// Get the email and password from the request body json and store them in the user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	var existingUser *User
	// Check if the user exists in the database
	result := DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		// Return a 404 if the user does not exist
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the password is correct
	errHash := user.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		fmt.Println("Password is incorrect")
		// Return a 401 if the password is incorrect
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new token
	generatedToken, err := token.GenerateToken(models.User{ID: existingUser.ID})
	if err != nil {
		fmt.Println(err)
		// Return a 500 if the token could not be generated
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Store the token in redis
	err = Redis.Set(generatedToken, user.ID, 0).Err()
	if err != nil {
		fmt.Println(err)
		// Return a 500 if the token could not be stored in redis
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the session token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    generatedToken,
		HttpOnly: true,
	})

	// Return the token as json
	json.NewEncoder(w).Encode(generatedToken)

	fmt.Println("Login successful")

	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var token models.Token
	Redis := initializers.Redis

	// Get the token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println(err)
		// Return a 400 if the cookie does not exist
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(cookie.Value)

	// Delete the token from redis
	err = Redis.Del(cookie.Value).Err()
	if err != nil {
		fmt.Println(err)
		// Return a 500 if the token could not be deleted from redis
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the cookie from the response
	c := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, c)

	// Delete the token from the database.
	err = token.Destroy(cookie.Value)
	if err != nil {
		fmt.Println(err)
		// Return a 500 if the token could not be deleted from the database
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Logout successful")

	return
}
