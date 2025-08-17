package db

import (
	"database/sql"
	"fmt"
	"log"
	"speechToText/src/config"
)

var db = InitDB()

func InitDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.CurrentConfig.Database.Host,
		config.CurrentConfig.Database.Port,
		config.CurrentConfig.Database.Username,
		config.CurrentConfig.Database.Password,
		config.CurrentConfig.Database.Name)
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
		log.Fatal("Db connection error:", err)
		return nil
	}
	var query1 = `
    CREATE TABLE IF NOT EXISTS users (
        username VARCHAR(1000) NOT NULL,
        password VARCHAR(1000) NOT NULL,
    );`
	_, err = db.Exec(query1)
	if err != nil {
		log.Fatal("Db create table error:", err)
	}
	var query2 = `
    CREATE TABLE IF NOT EXISTS tasks (
        username VARCHAR(1000) NOT NULL,
        task_id VARCHAR(1000) NOT NULL,
        audio VARCHAR(1000) NOT NULL,
        status VARCHAR(1000) NOT NULL,
        result VARCHAR(1000),
    );`
	_, err = db.Exec(query2)
	if err != nil {
		log.Fatal("Db error:", err)
	}
	return db
}
