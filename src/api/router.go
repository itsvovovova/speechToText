package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"speechToText/src/cache"
	"speechToText/src/consumer"
	"speechToText/src/db"
	"speechToText/src/service"
	"speechToText/src/types"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func validateAudioURL(audioURL string) error {
	if audioURL == "" {
		return fmt.Errorf("audio URL is required")
	}
	u, err := url.Parse(audioURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("invalid audio URL: must be an http or https URL")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, v any) {
	response, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		service.LogError("Write error: %v", err)
	}
}

// Audio godoc
// @Summary Upload audio for processing
// @Description Sends audio URL for speech to text conversion
// @Tags audio
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.AudioRequest true "Audio URL"
// @Success 200 {object} types.GetInfoResponse "Task ID created"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /audio [post]
func Audio(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	session, err := cache.SessionManager.SessionGet(r.Context(), r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var request types.AudioRequest
	if err = json.Unmarshal(data, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validateAudioURL(request.Audio); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskID, err := consumer.CreateTask(username, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, types.GetInfoResponse{Task_id: taskID})
}

// Status godoc
// @Summary Get task status
// @Description Returns current audio processing status
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param task_id query string true "Task ID"
// @Success 200 {object} types.GetStatusResponse "Task status"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /status [get]
func Status(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionGet(r.Context(), r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		http.Error(w, "task_id is required", http.StatusBadRequest)
		return
	}
	exist, err := db.ExistTask(taskID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	status, err := db.GetStatusTask(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, types.GetStatusResponse{Status: status})
}

// Result godoc
// @Summary Get processing result
// @Description Returns audio to text conversion result
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param task_id query string true "Task ID"
// @Success 200 {object} types.GetResultResponse "Processing result"
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /result [get]
func Result(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionGet(r.Context(), r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		http.Error(w, "task_id is required", http.StatusBadRequest)
		return
	}
	exist, err := db.ExistTask(taskID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	result, err := db.GetResultTask(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, types.GetResultResponse{Result: result})
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
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
func Tasks(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionGet(r.Context(), r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil || username == "" {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	writeJSON(w, types.TaskListResponse{
		Tasks: tasks,
		Pagination: types.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task by ID; only the owner can delete
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string "Task deleted"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	session, err := cache.SessionManager.SessionGet(r.Context(), r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, err := session.Get(r.Context(), "username")
	if err != nil || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "task id is required", http.StatusBadRequest)
		return
	}
	exist, err := db.ExistTask(taskID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := db.DeleteTask(taskID, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"result": "ok"})
}
