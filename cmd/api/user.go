package api

import (
	"fmt"

	"encoding/json"
	"net/http"

	db "github.com/MateSousa/cloud-build/initializers"
	"github.com/MateSousa/cloud-build/models"
	"golang.org/x/crypto/bcrypt"
)

type User models.User

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var u *User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
		// return 400 if the request body is not valid
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if user exists return error
	err = u.FindByEmail(u)
	if err == nil {
		fmt.Println("User already exists")
		// return 409 if the user already exists
		w.WriteHeader(http.StatusConflict)
		return
	}

	// hash the password
	hashedPassword, err := GenerateHashPassword(u.Password)
	if err != nil {
		fmt.Println(err)
		// return 500 if the password could not be hashed
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set the hashed password
	u.Password = hashedPassword

	result := db.DB.Create(&u)
	if result.Error != nil {
		fmt.Println(result.Error)
		// return 500 if the user could not be created
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// return message to user saying that the user was created in json
	json.NewEncoder(w).Encode(u)
}

func (u *User) FindByEmail(payload *User) error {
	result := db.DB.Where("email = ?", payload.Email).First(u)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *User) CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
