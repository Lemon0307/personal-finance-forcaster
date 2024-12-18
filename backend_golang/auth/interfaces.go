package auth

import "net/http"

type AuthenticationService interface {
	Login(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
}
