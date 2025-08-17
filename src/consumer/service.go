package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"speechToText/src/db"
	"speechToText/src/types"
)

func connectionAMQP(url string) *amqp.Connection {
	connection, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return connection
}

func declareQueue(RabbitMQUrl string, nameQueue string) (*types.QueueRabbitMQ, error) {
	connection := connectionAMQP(RabbitMQUrl)
	channel, err := connection.Channel()
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	queue, err := channel.QueueDeclare(nameQueue, false, false, false, false, nil)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	return &types.QueueRabbitMQ{
		Queue:      &queue,
		Channel:    channel,
		Connection: connection,
	}, nil
}

func SendMessage(taskID string, nameQueue string, audioUrl string, RabbitMQUrl string) error {
	queueSettings, err := declareQueue(RabbitMQUrl, nameQueue)
	if err != nil {
		return err
	}
	defer func() {
		err := closeQueueRabbitMQ(queueSettings)
		if err != nil {
			fmt.Println(err)
		}
	}()
	body := types.AudioMessage{
		TaskID: taskID,
		Audio:  audioUrl,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = queueSettings.Channel.Publish(
		"",
		queueSettings.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		})
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMessage(nameQueue string, ctx context.Context) error {
	queueSettings, err := declareQueue(nameQueue, nameQueue)
	if err != nil {
		return err
	}
	defer func() {
		err := closeQueueRabbitMQ(queueSettings)
		if err != nil {
			fmt.Println(err)
		}
	}()
	messages, err := queueSettings.Channel.ConsumeWithContext(
		ctx,
		nameQueue,
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
				err := processingData(message)
				if err != nil {
					fmt.Println("error: ", err)
					err = message.Nack(false, true)
					if err != nil {
						return
					}
				} else {
					err := message.Ack(false)
					if err != nil {
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
	err := json.Unmarshal(data.Body, &audio)
	if err != nil {
		return err
	}
	text, err := ConvertToText(audio.Audio)
	if err != nil {
		return err
	}
	if err = db.AddResultTask(audio.TaskID, text); err != nil {
		return err
	}
	return nil
}

func closeQueueRabbitMQ(queue *types.QueueRabbitMQ) error {
	if queue.Channel != nil {
		err := queue.Channel.Close()
		if err != nil {
			return err
		}
	}
	if queue.Connection != nil {
		return queue.Connection.Close()
	}
	return nil
}
