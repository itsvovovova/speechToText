package types

import amqp "github.com/rabbitmq/amqp091-go"

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

type AudioRequest struct {
	Audio string
}

type AudioMessage struct {
	Audio  string
	TaskID string
}

type QueueRabbitMQ struct {
	Queue      *amqp.Queue
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

type GetInfoResponse struct {
	Task_id string `json:"task_id"`
}

type GetResultResponse struct {
	Result string `json:"result"`
}

type GetStatusResponse struct {
	Status string `json:"status"`
}

type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type TaskListResponse struct {
	Tasks      []TaskInfo         `json:"tasks"`
	Pagination PaginationResponse `json:"pagination"`
}

type TaskInfo struct {
	TaskID   string `json:"task_id"`
	Username string `json:"username"`
	Status   string `json:"status"`
	Created  string `json:"created"`
}
