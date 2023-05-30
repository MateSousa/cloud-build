package api

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()


// auth is a struct that holds the configuration for the auth api
func AuthHandler(r *mux.Router) {
	r.HandleFunc("/auth/login", Login).Methods("POST")
	r.HandleFunc("/auth/logout", Logout).Methods("DELETE")
}

func UserHandler(r *mux.Router) {
	r.HandleFunc("/user/create", CreateUser).Methods("POST")
}

func Init() {
	api := r.PathPrefix("/api/v1").Subrouter()

	AuthHandler(api)
	UserHandler(api)

	// init the server
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
