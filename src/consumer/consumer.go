package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	listen "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func ConvertToText(audioUrl string) (string, error) {
	client.InitWithDefault()
	ctx := context.Background()

	options := &interfaces.PreRecordedTranscriptionOptions{
		Model:    "nova-2",
		Language: "ru",
	}

	c := client.NewREST("baa405ef4c3ac905a131ff173d19320313f966a2", &interfaces.ClientOptions{})
	dg := listen.New(c)

	res, err := dg.FromURL(ctx, audioUrl, options)
	if err != nil {
		fmt.Printf("FromURL failed. Err: %v\n", err)
		return "", err
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		return "", err
	}

	prettyJson, err := prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		return "", err
	}
	return string(prettyJson), nil
}
