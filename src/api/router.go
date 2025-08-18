package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"speechToText/src/cache"
	"speechToText/src/consumer"
	"speechToText/src/db"
	"speechToText/src/service"
	"speechToText/src/types"
)

func Audio(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionStart(r.Context(), w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var request types.AudioRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskID, err := consumer.CreateTask(username, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(taskID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Status(w http.ResponseWriter, r *http.Request) {
	service.LogInfo("=== STATUS START ===")
	defer service.LogInfo("=== STATUS END ===")

	w.Header().Set("Content-Type", "application/json")
	session, err := cache.SessionManager.SessionStart(r.Context(), w, r)
	if err != nil {
		service.LogError("Session error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	service.LogDebug("Session found: %s", session.SessionId)

	username, err := session.Get(r.Context(), "username")
	if err != nil {
		service.LogError("Username error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	service.LogDebug("Username: %s", username)

	status, err := db.GetStatusTask(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.LogDebug("No tasks found for user")
			if _, err := w.Write([]byte("no tasks")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		service.LogError("DB error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	service.LogDebug("Status: %s", status)
	if _, err := w.Write([]byte(status)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Result(w http.ResponseWriter, r *http.Request) {
	service.LogInfo("=== RESULT START ===")
	defer service.LogInfo("=== RESULT END ===")

	w.Header().Set("Content-Type", "application/json")
	session, err := cache.SessionManager.SessionStart(r.Context(), w, r)
	if err != nil {
		service.LogError("Session error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	service.LogDebug("Session found: %s", session.SessionId)

	username, err := session.Get(r.Context(), "username")
	if err != nil {
		service.LogError("Username error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	service.LogDebug("Username: %s", username)

	result, err := db.GetResultTask(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.LogDebug("No results found for user")
			if _, err := w.Write([]byte("no results")); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		service.LogError("DB error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	service.LogDebug("Result: %s", result)
	if _, err := w.Write([]byte(result)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
