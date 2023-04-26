package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"testing"
)

func TestChatGPT3Dot5Turbo(t *testing.T) {
	gpt := NewChatGptTool("", "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")
	message := []Gpt3Dot5Message{
		{
			Role:    "system",
			Content: "你是一个精通开发的资深工程师，熟悉全栈技术，任何问题都难不倒你",
		},
		{
			Role:    "user",
			Content: "帮我使用golang开发一个在线客服系统",
		},
	}
	res, err := gpt.ChatGPT3Dot5Turbo(message)
	log.Println(res, err)
}
func TestChatGPT3Dot5Stream(t *testing.T) {
	gpt := NewChatGptTool("", "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")
	message := []Gpt3Dot5Message{
		{
			Role:    "system",
			Content: "你是一个精通开发的资深工程师，熟悉全栈技术，任何问题都难不倒你",
		},
		{
			Role:    "user",
			Content: "帮我使用golang开发一个在线客服系统",
		},
	}
	stream, _ := gpt.ChatGPT3Dot5TurboStream(message)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("\nStream finished")
			break
		} else if err != nil {
			log.Printf("\nStream error: %v\n", err)
			break
		} else {
			log.Println(response.Choices[0].Delta.Content, err)
		}
	}
	stream.Close()
}
func TestGetEmbedding(t *testing.T) {
	gpt := NewChatGptTool("https://openai.api2d.net", "fk188528-JAPbwe87SKzXwGBroAIdcOLfSC1bAMVU")
	response, err := gpt.GetEmbedding("测试", "text-embedding-ada-002")
	if err != nil {
		t.Errorf("Error GetEmbedding: %v", err)
	}

	var embeddingResponse EmbeddingResponse
	json.Unmarshal([]byte(response), &embeddingResponse)
	log.Println(embeddingResponse)
}
