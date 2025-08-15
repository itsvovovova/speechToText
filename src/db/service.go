package db

import (
	"log"
)

func AddAuthData(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	if _, err := db.Exec(query, username, password); err != nil {
		return err
	}
	return nil
}

func CheckAuthData(username string, password string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND password = $2)", username, password).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists, nil
}

func AddAudioTask(username string, audio string) error {
	query := "INSERT INTO audio (username, task,) VALUES ($1, $2)"
	if _, err := db.Exec(query, username, audio); err != nil {
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
