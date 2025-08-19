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
