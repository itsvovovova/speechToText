package types

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session interface {
	GetSessionId() string
	Get(key interface{}) (interface{}, error)
	Delete(key interface{}) error
	Set(key, value interface{}) error
}
