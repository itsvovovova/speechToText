package service

import (
	"encoding/json"
	"io"
	"net/http"
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
