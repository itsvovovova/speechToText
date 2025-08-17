package api

import (
	"encoding/json"
	"net/http"
	"speechToText/src/cache"
	"speechToText/src/db"
	"speechToText/src/service"
)

func Register(w http.ResponseWriter, r *http.Request) {
	user, err := service.ReadAuthRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	HashPassword, err := db.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := db.AddAuthData(user.Username, HashPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response = struct {
		Result string `json:"result"`
	}{
		Result: "ok",
	}
	rvalue, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(rvalue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	user, err := service.ReadAuthRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exist, err := db.CheckAuthData(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !exist {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	session, err := cache.SessionManager.SessionStart(ctx, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := session.Set(ctx, "username", user.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Result string `json:"result"`
		Token  string `json:"token"`
	}{
		Result: "ok",
		Token:  session.SessionId,
	}
	rvalue, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(rvalue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
