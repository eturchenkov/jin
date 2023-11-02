package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"os"
)

func buildAgent() func(string) string {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	return func(prompt string) string {
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				Temperature: 0.3,
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return ""
		}

		return resp.Choices[0].Message.Content
	}
}
