package api

import (
	"encoding/json"
	"io"
	"net/http"
	"speechToText/src/cache"
	"speechToText/src/db"
	"speechToText/src/types"
)

func audio(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionStart(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if err := db.AddAudioTask(username, request.Audio); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// send task to queue rabbitmq
}
