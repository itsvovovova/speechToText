package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"speechToText/src/api"
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
			name: "Valid audio request",
			requestBody: types.AudioRequest{
				Audio: "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav",
			},
			expectedStatus: 200,
			expectTaskID:   true,
		},
		{
			name: "Empty audio",
			requestBody: types.AudioRequest{
				Audio: "",
			},
			expectedStatus: 400,
			expectTaskID:   false,
		},
		{
			name: "Invalid URL (not http/https)",
			requestBody: types.AudioRequest{
				Audio: "not-a-url",
			},
			expectedStatus: 400,
			expectTaskID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/audio", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			api.Audio(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
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
		{
			name:           "Valid status request",
			queryParams:    "?task_id=test-task-id",
			expectedStatus: 200,
		},
		{
			name:           "Missing task_id",
			queryParams:    "",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/status"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			api.Status(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
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
		{
			name:           "Valid result request",
			queryParams:    "?task_id=test-task-id",
			expectedStatus: 200,
		},
		{
			name:           "Missing task_id",
			queryParams:    "",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/result"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			api.Result(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
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
		{
			name:           "Valid tasks request",
			queryParams:    "?page=1&page_size=10",
			expectedStatus: 200,
			expectTasks:    true,
		},
		{
			name:           "Invalid page parameter",
			queryParams:    "?page=0&page_size=10",
			expectedStatus: 200,
			expectTasks:    true,
		},
		{
			name:           "Invalid page_size parameter",
			queryParams:    "?page=1&page_size=0",
			expectedStatus: 200,
			expectTasks:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/tasks"+tt.queryParams, nil)

			rr := httptest.NewRecorder()
			api.Tasks(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == 200 && tt.expectTasks {
				var response types.TaskListResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Pagination.Page == 0 {
					t.Errorf("Expected Pagination.Page in response but not found")
				}
			}
		})
	}
}
