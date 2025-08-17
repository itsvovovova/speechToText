package consumer

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"speechToText/src/types"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connectionAMQP(url string) *amqp.Connection {
	connection, err := amqp.Dial(url)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}
	return connection
}

func SendMessage(taskID string, nameQueue string, audioUrl string) {
	connection := connectionAMQP(audioUrl)
	defer func(Connection *amqp.Connection) {
		err := Connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(channel)
	queue, err := channel.QueueDeclare(nameQueue, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	body := types.AudioMessage{
		TaskID: taskID,
		Audio:  audioUrl,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		})
}
