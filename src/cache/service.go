package cache

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (session RedisSession) Set(ctx context.Context, key string, value interface{}) error {
	return session.Client.Set(ctx, session.SessionId+key, value, session.TTL).Err()
}

func (session RedisSession) Get(ctx context.Context, key string) (string, error) {
	return session.Client.Get(ctx, session.SessionId+key).Result()
}

func (session RedisSession) Delete(ctx context.Context, key string) error {
	return session.Client.Del(ctx, session.SessionId+key).Err()
}

func NewRedisSessionProvider(address string) RedisSessionProvider {
	return RedisSessionProvider{
		Client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (p *RedisSessionProvider) Close() error {
	return p.Client.Close()
}

func (sessionProvider *RedisSessionProvider) SessionRead(sessionID string, ctx context.Context, ttl time.Duration) (*RedisSession, error) {
	exists, err := sessionProvider.Client.Exists(ctx, sessionID).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		if err = sessionProvider.Client.Set(ctx, sessionID, "", ttl).Err(); err != nil {
			return nil, err
		}
	}
	return &RedisSession{
		SessionId: sessionID,
		Client:    sessionProvider.Client,
		TTL:       ttl,
	}, nil
}

func (sessionProvider *RedisSessionProvider) SessionDelete(sessionID string, ctx context.Context) error {
	return sessionProvider.Client.Del(ctx, sessionID).Err()
}

func NewRedisSessionManager(cookieName string, provider RedisSessionProvider, maxLifeTime int64) *RedisSessionManager {
	return &RedisSessionManager{
		Provider:    &provider,
		Cookie:      cookieName,
		MaxLifetime: time.Duration(maxLifeTime) * time.Second,
	}
}

func (manager *RedisSessionManager) GenerateSessionID() (string, error) {
	return uuid.New().String(), nil
}

// SessionGet reads an existing session without creating a new one.
// Returns an error if no valid session cookie is present.
func (manager *RedisSessionManager) SessionGet(ctx context.Context, r *http.Request) (*RedisSession, error) {
	cookie, err := r.Cookie(manager.Cookie)
	if err != nil || cookie.Value == "" {
		return nil, fmt.Errorf("no session cookie")
	}
	exists, err := manager.Provider.Client.Exists(ctx, cookie.Value).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, fmt.Errorf("session not found")
	}
	return &RedisSession{
		SessionId: cookie.Value,
		Client:    manager.Provider.Client,
		TTL:       manager.MaxLifetime,
	}, nil
}

// SessionStart reads an existing session or creates a new one if none exists.
// Use only when a session must be created (e.g. on login).
func (manager *RedisSessionManager) SessionStart(ctx context.Context, w http.ResponseWriter, r *http.Request) (*RedisSession, error) {
	cookie, err := r.Cookie(manager.Cookie)
	if err != nil || cookie.Value == "" {
		sid, err := manager.GenerateSessionID()
		if err != nil {
			return nil, err
		}
		session, err := manager.Provider.SessionRead(sid, ctx, manager.MaxLifetime)
		if err != nil {
			return nil, err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     manager.Cookie,
			Value:    sid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.MaxLifetime.Seconds()),
		})
		return session, nil
	}
	return manager.Provider.SessionRead(cookie.Value, ctx, manager.MaxLifetime)
}

// SessionDestroy deletes the session from Redis and clears the cookie.
func (manager *RedisSessionManager) SessionDestroy(ctx context.Context, w http.ResponseWriter, sessionID string) error {
	// Delete all keys belonging to this session (root key + all data keys like sessionID+"username")
	keys, err := manager.Provider.Client.Keys(ctx, sessionID+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		if err := manager.Provider.Client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     manager.Cookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	return nil
}
