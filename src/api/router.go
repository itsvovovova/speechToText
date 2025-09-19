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
	"strconv"
)

// Audio godoc
// @Summary Upload audio for processing
// @Description Sends audio file for speech to text conversion
// @Tags audio
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.AudioRequest true "Audio data"
// @Success 200 {object} types.GetInfoResponse "Task ID created"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /audio [post]
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
	taskResponse := types.GetInfoResponse{
		Task_id: taskID,
	}
	response, err := json.Marshal(taskResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Status godoc
// @Summary Get task status
// @Description Returns current audio processing status
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.GetInfoResponse true "Task ID"
// @Success 200 {object} types.GetStatusResponse "Task status"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /status [get]
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var request types.GetInfoResponse
	err = json.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exist, err := db.ExistTask(request.Task_id, username)
	if !exist {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
	status, err := db.GetStatusTask(request.Task_id)
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
	statusResponse := types.GetStatusResponse{
		Status: status,
	}
	response, err := json.Marshal(statusResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Result godoc
// @Summary Get processing result
// @Description Returns audio to text conversion result
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.GetInfoResponse true "Task ID"
// @Success 200 {object} types.GetResultResponse "Processing result"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /result [get]
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var request types.GetInfoResponse
	err = json.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exist, err := db.ExistTask(request.Task_id, username)
	if !exist {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
	result, err := db.GetResultTask(request.Task_id)
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
	resultResponse := types.GetResultResponse{
		Result: result,
	}
	response, err := json.Marshal(resultResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Tasks godoc
// @Summary Get tasks list with pagination
// @Description Returns user tasks list with pagination
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} types.TaskListResponse "Tasks list"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
func Tasks(w http.ResponseWriter, r *http.Request) {
	service.LogInfo("=== TASKS START ===")
	defer service.LogInfo("=== TASKS END ===")

	w.Header().Set("Content-Type", "application/json")
	session, err := cache.SessionManager.SessionStart(r.Context(), w, r)
	if err != nil {
		service.LogError("Session error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := session.Get(r.Context(), "username")
	if err != nil {
		service.LogError("Username error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	tasks, total, err := db.GetTasksWithPagination(username, page, pageSize)
	if err != nil {
		service.LogError("DB error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	response := types.TaskListResponse{
		Tasks: tasks,
		Pagination: types.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		service.LogError("JSON marshal error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(responseJSON); err != nil {
		service.LogError("Write error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
