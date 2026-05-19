package external

import (
	"context"
	"cvbuilder/config"
	"fmt"
	"log"

	openrouter "github.com/revrost/go-openrouter"
)

type LLM struct {
	client *openrouter.Client
	c      *config.Config
}

func InitLLM(c *config.Config) *LLM {
	client := openrouter.NewClient(
		c.OpenrouterToken,
		openrouter.WithXTitle("offer_hustler"),
		openrouter.WithHTTPReferer("https://myapp.com"),
	)

	return &LLM{
		client: client,
		c:      c,
	}
}

func (llm *LLM) ChatCompletion(message string, model string) (string, error) {
	log.Printf("req start")
	resp, err := llm.client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: model,
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

func (llm *LLM) ChechHealth() error {
	_, err := llm.client.ListModels(context.TODO())

	if err != nil {
		log.Println("Model check health failed")
		return err
	} else {
		log.Println("Model check health success")
		return nil
	}
}
