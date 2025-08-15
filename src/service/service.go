package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
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

func CreateTask(request types.AudioRequest, session string) string {
	taskID := uuid.New().String()

	db.AddAudioTask()
	// send to rabbitmq with send_message.go in consumer
	return ""
}
