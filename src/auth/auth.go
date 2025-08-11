package auth

import (
	"net/http"
	"speechToText/src/db"
	"speechToText/src/service"
)

func register(w http.ResponseWriter, r *http.Request) {
	user, err := service.ReadAuthRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.AddAuthData(user.Username, user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
