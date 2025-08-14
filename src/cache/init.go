package cache

import "math"

var sessionProvider = NewRedisSessionProvider(":8080")
var SessionManager = NewRedisSessionManager("session_id", sessionProvider, int64(math.Pow10(5)))
