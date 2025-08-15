package db

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDB() *sql.DB {
	connStr := "host=localhost port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
		return nil
	}
	var query1 = `
    CREATE TABLE IF NOT EXISTS users (
        username VARCHAR(1000) NOT NULL,
        password VARCHAR(1000) NOT NULL,
    );`
	_, err = db.Exec(query1)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	var query2 = `
    CREATE TABLE IF NOT EXISTS tasks (
        username VARCHAR(1000) NOT NULL,
        task VARCHAR(1000) NOT NULL,
        status VARCHAR(1000) NOT NULL,
        result VARCHAR(1000) NOT NULL,
    );`
	_, err = db.Exec(query2)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	return db
}
