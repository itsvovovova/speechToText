package consumer

import (
	"context"
	"fmt"
	"speechToText/src/config"
	"speechToText/src/db"
	"speechToText/src/service"
	"speechToText/src/types"

	"github.com/google/uuid"

	listen "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

func CreateTask(username string, request types.AudioRequest) (string, error) {
	taskID := uuid.New().String()
	if err := db.AddAudioTask(taskID, username, request.Audio); err != nil {
		return "", err
	}
	if err := SendMessage(taskID, "queue", request.Audio, config.CurrentConfig.RabbitMQ.Url); err != nil {
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

	if len(res.Results.Channels) == 0 || len(res.Results.Channels[0].Alternatives) == 0 {
		return "", fmt.Errorf("no transcription result returned by Deepgram")
	}

	return res.Results.Channels[0].Alternatives[0].Transcript, nil
}
