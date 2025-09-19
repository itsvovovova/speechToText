package db

import (
	"database/sql"
	"errors"
	"fmt"
	"speechToText/src/service"
	"speechToText/src/types"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func AddAuthData(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	if _, err := db.Exec(query, username, password); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckAuthData(username string, password string) (bool, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}

func AddAudioTask(taskID string, username string, audio string) error {
	query := "INSERT INTO tasks (username, task_id, audio, status) VALUES ($1, $2, $3, $4)"
	if _, err := db.Exec(query, username, taskID, audio, "in progress"); err != nil {
		return err
	}
	return nil
}

func GetStatusTask(taskID string) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM tasks WHERE task_id = $1", taskID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func GetResultTask(taskID string) (string, error) {
	var result sql.NullString
	err := db.QueryRow("SELECT result FROM tasks WHERE task_id = $1", taskID).Scan(&result)
	if err != nil {
		return "", err
	}
	if result.Valid {
		return result.String, nil
	}

	return "in progress", nil
}

func AddResultTask(taskID string, text string) error {
	service.LogDebug("ADD RESULT TASK IS WORKING!")
	service.LogDebug("TEXT: %s", text)
	query := "UPDATE tasks SET result = $2, status = 'completed' WHERE task_id = $1"
	_, err := db.Exec(query, taskID, text)
	return err
}

func ExistUsername(username string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking username existence: %w", err)
	}

	return exists, nil
}

func ExistTask(taskID string, username string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = $1 AND username = $2)"

	err := db.QueryRow(query, taskID, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetTasksWithPagination(username string, page, pageSize int) ([]types.TaskInfo, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	countQuery := "SELECT COUNT(*) FROM tasks WHERE username = $1"
	err := db.QueryRow(countQuery, username).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT task_id, username, status, created_at 
		FROM tasks 
		WHERE username = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, username, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []types.TaskInfo
	for rows.Next() {
		var task types.TaskInfo
		var createdAt time.Time

		err := rows.Scan(&task.TaskID, &task.Username, &task.Status, &createdAt)
		if err != nil {
			return nil, 0, err
		}

		task.Created = createdAt.Format(time.RFC3339)
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func GetAllTasksWithPagination(page, pageSize int) ([]types.TaskInfo, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	countQuery := "SELECT COUNT(*) FROM tasks"
	err := db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT task_id, username, status, created_at 
		FROM tasks 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []types.TaskInfo
	for rows.Next() {
		var task types.TaskInfo
		var createdAt time.Time

		err := rows.Scan(&task.TaskID, &task.Username, &task.Status, &createdAt)
		if err != nil {
			return nil, 0, err
		}

		task.Created = createdAt.Format(time.RFC3339)
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}
