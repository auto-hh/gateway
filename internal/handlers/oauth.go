package handlers

import (
	"fmt"
	"net/http"
)

type OAuthHandler struct {
}

func NewOAuthHandler(backendHost, frontendHost string) *OAuthHandler {
	return &OAuthHandler{}
}
func (o *OAuthHandler) Begin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthBeginHandler is ready to work")
	//userId, err := r.Cookie(utils.CookieKeyUserId)
	//if err == nil {
	//
	//}
}

func (o *OAuthHandler) Complete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthEndHandler is ready to work")
}
