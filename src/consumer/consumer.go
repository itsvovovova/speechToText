package consumer

import (
	"context"
	"speechToText/src/db"
	"speechToText/src/service"
	"speechToText/src/types"

	"github.com/google/uuid"

	"speechToText/src/config"

	listen "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
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
	service.LogDebug("AUDIO URL: %s", audioUrl)
	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:    "nova-2",
		Language: "en",
	}

	c := client.NewREST(config.CurrentConfig.Deepgram.ApiKey, &interfaces.ClientOptions{})
	dg := listen.New(c)

	res, err := dg.FromURL(ctx, audioUrl, options)
	if err != nil {
		service.LogError("FromURL failed. Err: %v", err)
		return "", err
	}
	rtext := res.Results.Channels[0].Alternatives[0].Transcript

	return rtext, nil
}
