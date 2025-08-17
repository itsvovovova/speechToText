package cache

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

func (session RedisSession) Set(ctx context.Context, key string, value interface{}) error {
	return session.Client.Set(ctx, session.SessionId+key, value, 0).Err()
}

func (session RedisSession) Get(ctx context.Context, key string) (string, error) {
	return session.Client.Get(ctx, session.SessionId+key).Result()
}

func (session RedisSession) Delete(ctx context.Context, key string) error {
	return session.Client.Del(ctx, session.SessionId+key).Err()
}

func (session RedisSession) Clear(ctx context.Context) error {
	return session.Client.Del(ctx, session.SessionId+"*").Err()
}

func NewRedisSessionProvider(address string) RedisSessionProvider {
	return RedisSessionProvider{
		Client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (sessionProvider *RedisSessionProvider) SessionRead(session string, ctx context.Context) (*RedisSession, error) {
	exists, err := sessionProvider.Client.Exists(ctx, session).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		err = sessionProvider.Client.Set(ctx, session, "", time.Hour).Err()
		if err != nil {
			return nil, err
		}
	}
	return &RedisSession{
		SessionId: session,
		Client:    sessionProvider.Client,
	}, nil
}

func (sessionProvider *RedisSessionProvider) SessionDelete(session string, ctx context.Context) error {
	return sessionProvider.Client.Del(ctx, session).Err()
}

func NewRedisSessionManager(cookieName string, provider RedisSessionProvider, maxLifeTime int64) *RedisSessionManager {
	return &RedisSessionManager{
		Provider:    &provider,
		Cookie:      cookieName,
		MaxLifetime: time.Duration(maxLifeTime),
	}
}

func (manager *RedisSessionManager) GenerateSessionID() (string, error) {
	return uuid.New().String(), nil
}

func (manager *RedisSessionManager) SessionStart(ctx context.Context, w http.ResponseWriter, r *http.Request) (*RedisSession, error) {
	cookie, err := r.Cookie(manager.Cookie)
	if err != nil || cookie.Value == "" {
		sid, err := manager.GenerateSessionID()
		if err != nil {
			return nil, err
		}
		session, err := manager.Provider.SessionRead(sid, ctx)
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
	} else {
		session, err := manager.Provider.SessionRead(cookie.Value, ctx)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
}
