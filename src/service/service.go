package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"speechToText/src/consumer"
	"speechToText/src/db"
	"speechToText/src/types"
)

func ReadAuthRequest(r *http.Request) (types.AuthRequest, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return types.AuthRequest{}, err
	}
	var authData types.AuthRequest
	err = json.Unmarshal(data, &authData)
	if err != nil {
		return types.AuthRequest{}, err
	}
	return authData, nil
}

func CreateTask(username string, request types.AudioRequest) (string, error) {
	taskID := uuid.New().String()
	if err := db.AddAudioTask(taskID, username, request.Audio); err != nil {
		return "", err
	}
	err := consumer.SendMessage(taskID, "queue", request.Audio, "url")
	if err != nil {
		return "", err
	}
	return taskID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
