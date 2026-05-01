package main

import (
	"log"
	"math"
	"os"
	"speechToText/src/api"
	"speechToText/src/cache"
	"speechToText/src/config"
	"speechToText/src/consumer"
	"speechToText/src/db"
	"testing"
)

var (
	testStore    *db.Store
	testHandlers *api.Handlers
)

func TestMain(m *testing.M) {
	sqlDB, err := db.NewDB(config.CurrentConfig.Database)
	if err != nil {
		log.Printf("DB not available, DB-dependent tests will be skipped: %v", err)
	} else {
		testStore = db.NewStore(sqlDB)
		defer sqlDB.Close()
	}

	sessionProvider := cache.NewRedisSessionProvider(config.CurrentConfig.Redis.Host)
	defer sessionProvider.Close()
	sessionManager := cache.NewRedisSessionManager("session_id", sessionProvider, int64(math.Pow10(5)))

	producer := consumer.NewProducer(config.CurrentConfig.RabbitMQ.Url)
	defer producer.Close()

	testHandlers = api.NewHandlers(testStore, sessionManager, producer)

	os.Exit(m.Run())
}
