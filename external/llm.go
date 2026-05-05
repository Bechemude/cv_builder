package external

import (
	"context"
	"cvbuilder/config"
	"fmt"

	openrouter "github.com/revrost/go-openrouter"
)

type LLM struct {
	client *openrouter.Client
}

func InitLLM(c *config.Config) *LLM {
	client := openrouter.NewClient(
		c.OpenrouterToken,
		openrouter.WithXTitle("offer_hustler"),
		openrouter.WithHTTPReferer("https://myapp.com"),
	)

	return &LLM{
		client: client,
	}
}

func (llm *LLM) ChatCompletion(message string) (string, error) {
	resp, err := llm.client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: "nvidia/nemotron-3-nano-30b-a3b:free",
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(message),
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v\n", err)
	}

	response := resp.Choices[0].Message.Content.Text

	return response, nil
}

