package db

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func AddAuthData(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	if _, err := db.Exec(query, username, password); err != nil {
		return err
	}
	return nil
}

func CheckAuthData(username string, password string) (bool, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND password = $2)", username, hashedPassword).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists, nil
}

func AddAudioTask(taskID string, username string, audio string) error {
	query := "INSERT INTO audio (username, task_id, audio_url, status) VALUES ($1, $2, $3)"
	if _, err := db.Exec(query, username, taskID, audio, "in progress"); err != nil {
		return err
	}
	return nil
}

func GetStatusTask(username string) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM task WHERE username = $1", username).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func GetResultTask(username string) (string, error) {
	var status string
	err := db.QueryRow("SELECT result FROM task WHERE username = $1", username).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func AddResultTask(taskID string, text string) error {
	query := "INSERT INTO result(task_id, result) VALUES ($1, $2)"
	if _, err := db.Exec(query, taskID, text); err != nil {
		return err
	}
	return nil
}

func ExistUsername(username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT exists FROM users WHERE username = $1", username).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return exists, nil

