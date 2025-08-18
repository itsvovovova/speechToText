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
	hashPassword, err := db.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exist, err := db.ExistUsername(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}
	if err := db.AddAuthData(user.Username, hashPassword); err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(rvalue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	service.LogInfo("=== LOGIN START ===")
	defer service.LogInfo("=== LOGIN END ===")

	user, err := service.ReadAuthRequest(r)
	if err != nil {
		service.LogError("ReadAuthRequest error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	service.LogDebug("User data: %+v", user)

	exist, err := db.CheckAuthData(user.Username, user.Password)
	if err != nil {
		service.LogError("CheckAuthData error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	service.LogDebug("Auth check result: %v", exist)

	if !exist {
		service.LogDebug("User not authenticated")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	service.LogDebug("Starting session")
	session, err := cache.SessionManager.SessionStart(ctx, w, r)
	if err != nil {
		service.LogError("Session error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	service.LogDebug("Session created: %s", session.SessionId)

	if err := session.Set(ctx, "username", user.Username); err != nil {
		service.LogError("Session set error: %v", err)
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
		service.LogError("JSON marshal error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(rvalue); err != nil {
		service.LogError("Write error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	service.LogInfo("Login successful")
}
