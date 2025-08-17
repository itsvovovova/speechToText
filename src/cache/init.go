package cache

import (
	"math"
	"speechToText/src/config"
)

var sessionProvider = NewRedisSessionProvider(config.CurrentConfig.Redis.Host)
var SessionManager = NewRedisSessionManager("session_id", sessionProvider, int64(math.Pow10(5)))
