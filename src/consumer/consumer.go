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

var (
	producerConn *amqp.Connection
	producerCh   *amqp.Channel
	producerMu   sync.Mutex
)

func getOrCreateProducerChannel(url string) (*amqp.Channel, error) {
	producerMu.Lock()
	defer producerMu.Unlock()

	if producerConn == nil || producerConn.IsClosed() {
		conn, err := amqp.Dial(url)
		if err != nil {
			return nil, err
		}
		producerConn = conn
		producerCh = nil
	}

	if producerCh == nil {
		ch, err := producerConn.Channel()
		if err != nil {
			return nil, err
		}
		producerCh = ch
	}

	return producerCh, nil
}

func resetProducerChannel() {
	producerMu.Lock()
	defer producerMu.Unlock()
	producerCh = nil
}

func SendMessage(taskID string, nameQueue string, audioUrl string, rabbitMQUrl string) error {
	ch, err := getOrCreateProducerChannel(rabbitMQUrl)
	if err != nil {
		return err
	}

	if _, err = ch.QueueDeclare(nameQueue, false, false, false, false, nil); err != nil {
		resetProducerChannel()
		return err
	}

	body := types.AudioMessage{
		TaskID: taskID,
		Audio:  audioUrl,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = ch.Publish("", nameQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bodyJSON,
	})
	if err != nil {
		resetProducerChannel()
	}
	return err
}

func ReceiveMessage(nameQueue string, ctx context.Context) error {
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

	queue, err := channel.QueueDeclare(nameQueue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	messages, err := channel.ConsumeWithContext(
		ctx,
		queue.Name,
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case message, ok := <-messages:
				if !ok {
					fmt.Println("Channel closed")
					return
				}
				if err := processingData(message); err != nil {
					fmt.Println("error:", err)
					if err := message.Nack(false, true); err != nil {
						return
					}
				} else {
					if err := message.Ack(false); err != nil {
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

func processingData(data amqp.Delivery) error {
	var audio types.AudioMessage
	if err := json.Unmarshal(data.Body, &audio); err != nil {
		return err
	}
	text, err := ConvertToText(audio.Audio)
	if err != nil {
		_ = db.UpdateTaskFailed(audio.TaskID)
		return err
	}
	if err := db.AddResultTask(audio.TaskID, text); err != nil {
		_ = db.UpdateTaskFailed(audio.TaskID)
		return err
	}
	return nil
}
