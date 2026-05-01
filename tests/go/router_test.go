package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"speechToText/src/types"
	"testing"
)

func TestAudio(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    types.AudioRequest
		expectedStatus int
		expectTaskID   bool
	}{
		{
			name:           "Valid audio request",
			requestBody:    types.AudioRequest{Audio: "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"},
			expectedStatus: 200,
			expectTaskID:   true,
		},
		{
			name:           "Empty audio",
			requestBody:    types.AudioRequest{Audio: ""},
			expectedStatus: 401,
		},
		{
			name:           "Invalid URL (not http/https)",
			requestBody:    types.AudioRequest{Audio: "not-a-url"},
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == 200 && testStore == nil {
				t.Skip("DB not available")
			}
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/audio", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			testHandlers.Audio(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == 200 && tt.expectTaskID {
				var response types.GetInfoResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Task_id == "" {
					t.Errorf("Expected TaskID in response but not found")
				}
			}
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{name: "Valid status request", queryParams: "?task_id=test-task-id", expectedStatus: 401},
		{name: "Missing task_id", queryParams: "", expectedStatus: 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/status"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			testHandlers.Status(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestResult(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{name: "Valid result request", queryParams: "?task_id=test-task-id", expectedStatus: 401},
		{name: "Missing task_id", queryParams: "", expectedStatus: 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/result"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			testHandlers.Result(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestTasks(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectTasks    bool
	}{
		{name: "Valid tasks request", queryParams: "?page=1&page_size=10", expectedStatus: 401},
		{name: "Invalid page parameter", queryParams: "?page=0&page_size=10", expectedStatus: 401},
		{name: "Invalid page_size parameter", queryParams: "?page=1&page_size=0", expectedStatus: 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/tasks"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			testHandlers.Tasks(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
