package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"speechToText/src/types"
)

func LogDebug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func LogInfo(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func ReadAuthRequest(r *http.Request) (types.AuthRequest, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return types.AuthRequest{}, err
	}
	LogDebug("Received data: %s", string(data))
	var authData types.AuthRequest
	err = json.Unmarshal(data, &authData)
	if err != nil {
		return types.AuthRequest{}, err
	}
	LogDebug("Parsed - Username='%s', Password='%s'", authData.Username, authData.Password)
	if authData.Username == "" || authData.Password == "" {
		LogDebug("Validation failed - empty fields")
		return types.AuthRequest{}, fmt.Errorf("username and password are required")
	}
	return authData, nil
}
