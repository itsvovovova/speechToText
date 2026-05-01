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

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AddAuthData(username string, password string) error {
	_, err := s.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *Store) CheckAuthData(username string, password string) (bool, error) {
	var hashedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil, nil
}

func (s *Store) AddAudioTask(taskID string, username string, audio string) error {
	_, err := s.db.Exec(
		"INSERT INTO tasks (username, task_id, audio, status) VALUES ($1, $2, $3, $4)",
		username, taskID, audio, "in progress",
	)
	return err
}

func (s *Store) GetStatusTask(taskID string) (string, error) {
	var status string
	err := s.db.QueryRow("SELECT status FROM tasks WHERE task_id = $1", taskID).Scan(&status)
	return status, err
}

func (s *Store) GetResultTask(taskID string) (string, error) {
	var result sql.NullString
	err := s.db.QueryRow("SELECT result FROM tasks WHERE task_id = $1", taskID).Scan(&result)
	if err != nil {
		return "", err
	}
	if result.Valid {
		return result.String, nil
	}
	return "in progress", nil
}

func (s *Store) AddResultTask(taskID string, text string) error {
	service.LogDebug("ADD RESULT TASK IS WORKING!")
	service.LogDebug("TEXT: %s", text)
	_, err := s.db.Exec("UPDATE tasks SET result = $2, status = 'completed' WHERE task_id = $1", taskID, text)
	return err
}

func (s *Store) UpdateTaskFailed(taskID string) error {
	_, err := s.db.Exec("UPDATE tasks SET status = 'failed' WHERE task_id = $1", taskID)
	return err
}

func (s *Store) ExistUsername(username string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking username existence: %w", err)
	}
	return exists, nil
}

func (s *Store) ExistTask(taskID string, username string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id = $1 AND username = $2)",
		taskID, username,
	).Scan(&exists)
	return exists, err
}

func (s *Store) DeleteTask(taskID string, username string) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE task_id = $1 AND username = $2", taskID, username)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (s *Store) GetTasksWithPagination(username string, page, pageSize int) ([]types.TaskInfo, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE username = $1", username).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(`
		SELECT task_id, username, status, created_at
		FROM tasks
		WHERE username = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		username, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []types.TaskInfo
	for rows.Next() {
		var task types.TaskInfo
		var createdAt time.Time
		if err := rows.Scan(&task.TaskID, &task.Username, &task.Status, &createdAt); err != nil {
			return nil, 0, err
		}
		task.Created = createdAt.Format(time.RFC3339)
		tasks = append(tasks, task)
	}
	return tasks, total, rows.Err()
}
