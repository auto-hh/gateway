package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/internal/service/oauth"
	"gateway/internal/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid/v5"
)

const SessionIdAgeTime = int(24 * time.Hour / time.Second)

type OAuthHandler struct {
	service *oauth.Service
	client  *http.Client
}

func NewOAuthHandler(Service *oauth.Service) *OAuthHandler {
	client := &http.Client{}
	return &OAuthHandler{
		service: Service,
		client:  client,
	}
}

func (o *OAuthHandler) Begin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthBeginHandler is ready to work")

	sessionIdCookie, err := r.Cookie(utils.CookieKeySessionId)

	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sessionId string

	if sessionIdCookie == nil || sessionIdCookie.Value == "" {
		sessionId, err = generateSessionId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     utils.CookieKeySessionId,
			Value:    sessionId,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   SessionIdAgeTime,
		})

	} else {
		sessionId = sessionIdCookie.Value
	}

	var redirectUrl *url.URL

	redirectUrl, err = o.service.BuildCodeRequest(r.Context(), sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectUrl.String(), http.StatusFound)
}

func (o *OAuthHandler) Complete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthEndHandler is ready to work")

	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sessionId string

	if sessionIdCookie == nil || sessionIdCookie.Value == "" {
		http.Error(w, "oauth.Complete: invalid sessionId", http.StatusUnauthorized)
		return
	} else {
		sessionId = sessionIdCookie.Value
	}

	code := r.URL.Query().Get("code")

	request, err := o.service.BuildTokenRequest(r.Context(), sessionId, r.URL.Query().Get("state"), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := o.client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "oauth.Complete: something went wrong in authorization", http.StatusUnauthorized)
		return
	}

	var results struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = o.service.SetTokens(
		request.Context(),
		sessionId,
		results.AccessToken,
		results.RefreshToken,
		time.Duration(results.ExpiresIn)*time.Second,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "authorized",
		Value:    "true",
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   SessionIdAgeTime,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func generateSessionId() (string, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("oauth.generateSessionId: %w", err)
	}
	sessionID := id.String()

	return sessionID, nil
}
