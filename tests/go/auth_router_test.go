package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"speechToText/src/api"
	"speechToText/src/types"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    types.AuthRequest
		expectedStatus int
		expectedResult string
	}{
		{
			name: "Valid registration",
			requestBody: types.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			expectedStatus: 200,
			expectedResult: "ok",
		},
		{
			name: "Empty username",
			requestBody: types.AuthRequest{
				Username: "",
				Password: "testpass",
			},
			expectedStatus: 400,
		},
		{
			name: "Empty password",
			requestBody: types.AuthRequest{
				Username: "testuser",
				Password: "",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			api.Register(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == 200 {
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response["result"] != tt.expectedResult {
					t.Errorf("handler returned wrong result: got %v want %v",
						response["result"], tt.expectedResult)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    types.AuthRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Valid login",
			requestBody: types.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			expectedStatus: 200,
			expectToken:    true,
		},
		{
			name: "Invalid credentials",
			requestBody: types.AuthRequest{
				Username: "wronguser",
				Password: "wrongpass",
			},
			expectedStatus: 401,
			expectToken:    false,
		},
		{
			name: "Empty username",
			requestBody: types.AuthRequest{
				Username: "",
				Password: "testpass",
			},
			expectedStatus: 400,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			api.Login(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == 200 && tt.expectToken {
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if _, exists := response["token"]; !exists {
					t.Errorf("Expected token in response but not found")
				}
				if response["result"] != "ok" {
					t.Errorf("Expected result 'ok', got %v", response["result"])
				}
			}
		})
	}
}
