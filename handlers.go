package main

import (
	"encoding/json"
	"net/http"
)

// Login logs in user.
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Parse form data.
	user, err := server.GetUserByEmail(r.FormValue("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Password != r.FormValue("password") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, _ := CreateAccessToken(map[string]interface{}{
		"userID":    user.ID,
		"userEmail": user.Email,
	})
	refreshToken := CreateRefreshToken()
	if _, err = server.StoreRefreshToken(refreshToken, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// Welcome lists down all todos.
func Welcome(w http.ResponseWriter, r *http.Request) {
	if !ValidToken(r) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	todos, err := server.GetAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// RenewJWT renews JWT.
func RenewJWT(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	r, err := server.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := server.GetUserByRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, _ := CreateAccessToken(map[string]interface{}{
		"userID":    user.ID,
		"userEmail": user.Email,
	})
	newRefreshToken := CreateRefreshToken()
	if _, err = server.StoreRefreshToken(newRefreshToken, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// Logout logs out the user.
func Logout(w http.ResponseWriter, r *http.Request) {
	// TODO
}
