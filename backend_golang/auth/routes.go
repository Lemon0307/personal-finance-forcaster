package auth

import "github.com/gorilla/mux"

func AuthenticationRoutes(router *mux.Router, authService AuthenticationService) {
	router.HandleFunc("/sign-up", authService.SignUp).Methods("POST")
	router.HandleFunc("/login", authService.Login).Methods("POST")
}