package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"speechToText/src/config"
	"speechToText/src/db"
	"speechToText/src/types"
)

// Producer manages a lazily-initialised RabbitMQ producer connection.
type Producer struct {
	url  string
	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewProducer(url string) *Producer {
	return &Producer{url: url}
}

func (p *Producer) channel() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil || p.conn.IsClosed() {
		conn, err := amqp.Dial(p.url)
		if err != nil {
			return nil, err
		}
		p.conn = conn
		p.ch = nil
	}

	if p.ch == nil {
		ch, err := p.conn.Channel()
		if err != nil {
			return nil, err
		}
		p.ch = ch
	}

	return p.ch, nil
}

func (p *Producer) resetChannel() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ch = nil
}

func (p *Producer) Send(taskID string, queueName string, audioUrl string) error {
	ch, err := p.channel()
	if err != nil {
		return err
	}

	if _, err = ch.QueueDeclare(queueName, false, false, false, false, nil); err != nil {
		p.resetChannel()
		return err
	}

	body, err := json.Marshal(types.AudioMessage{TaskID: taskID, Audio: audioUrl})
	if err != nil {
		return err
	}

	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		p.resetChannel()
	}
	return err
}

func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil && !p.conn.IsClosed() {
		return p.conn.Close()
	}
	return nil
}

// Consumer processes messages from a RabbitMQ queue.
type Consumer struct {
	store *db.Store
}

func NewConsumer(store *db.Store) *Consumer {
	return &Consumer{store: store}
}

func (c *Consumer) Receive(queueName string, ctx context.Context) error {
	connection, err := amqp.Dial(config.CurrentConfig.RabbitMQ.Url)
	if err != nil {
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	messages, err := channel.ConsumeWithContext(ctx, queue.Name, "", false, false, true, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-messages:
				if !ok {
					fmt.Println("Channel closed")
					return
				}
				if err := c.processMessage(msg); err != nil {
					fmt.Println("error:", err)
					if err := msg.Nack(false, true); err != nil {
						return
					}
				} else {
					if err := msg.Ack(false); err != nil {
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-ctx.Done()
	return nil
}

func (c *Consumer) processMessage(data amqp.Delivery) error {
	var audio types.AudioMessage
	if err := json.Unmarshal(data.Body, &audio); err != nil {
		return err
	}
	text, err := ConvertToText(audio.Audio)
	if err != nil {
		_ = c.store.UpdateTaskFailed(audio.TaskID)
		return err
	}
	if err := c.store.AddResultTask(audio.TaskID, text); err != nil {
		_ = c.store.UpdateTaskFailed(audio.TaskID)
		return err
	}
	return nil
}
