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
)

// Audio godoc
// @Summary Загрузка аудио для обработки
// @Description Отправляет аудио файл для преобразования в текст
// @Tags audio
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.AudioRequest true "Аудио данные"
// @Success 200 {object} types.GetInfoResponse "ID задачи создана"
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 401 {string} string "Не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
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
// @Summary Получение статуса задачи
// @Description Возвращает текущий статус обработки аудио
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.GetInfoResponse true "ID задачи"
// @Success 200 {object} types.GetStatusResponse "Статус задачи"
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 401 {string} string "Не авторизован"
// @Failure 403 {string} string "Доступ запрещен"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
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
// @Summary Получение результата обработки
// @Description Возвращает результат преобразования аудио в текст
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.GetInfoResponse true "ID задачи"
// @Success 200 {object} types.GetResultResponse "Результат обработки"
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 401 {string} string "Не авторизован"
// @Failure 403 {string} string "Доступ запрещен"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
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
