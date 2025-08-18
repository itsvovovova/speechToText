package consumer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"speechToText/src/db"
	"speechToText/src/service"
	"speechToText/src/types"

	"speechToText/src/config"

	listen "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func CreateTask(username string, request types.AudioRequest) (string, error) {
	taskID := uuid.New().String()
	if err := db.AddAudioTask(taskID, username, request.Audio); err != nil {
		return "", err
	}
	err := SendMessage(taskID, "queue", request.Audio, config.CurrentConfig.RabbitMQ.Url)
	if err != nil {
		return "", err
	}
	return taskID, nil
}

func ConvertToText(audioUrl string) (string, error) {
	ctx := context.Background()

	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:    "nova-2",
		Language: "ru",
	}

	c := client.NewREST(config.CurrentConfig.Deepgram.ApiKey, &interfaces.ClientOptions{})
	dg := listen.New(c)

	res, err := dg.FromURL(ctx, audioUrl, options)
	if err != nil {
		service.LogError("FromURL failed. Err: %v", err)
		return "", err
	}

	data, err := json.Marshal(res)
	if err != nil {
		service.LogError("json.Marshal failed. Err: %v", err)
		return "", err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		service.LogError("prettyjson.Marshal failed. Err: %v", err)
		return "", err
	}
	return string(prettyJson), nil
}
