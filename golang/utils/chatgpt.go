package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"net/http"
)

type ChatGptTool struct {
	Secret string
	Url    string
	Client *openai.Client
}
type Gpt3Dot5Message openai.ChatCompletionMessage

var (
	OPENAI_API_KEY  = ""
	OPENAI_API_BASE = "https://api.openai.com"
)

func NewChatGptTool() *ChatGptTool {
	config := openai.DefaultConfig(OPENAI_API_KEY)
	if OPENAI_API_BASE != "" {
		config.BaseURL = OPENAI_API_BASE + "/v1"
	}
	client := openai.NewClientWithConfig(config)
	//client := openai.NewClient(secret)
	return &ChatGptTool{
		Secret: OPENAI_API_KEY,
		Client: client,
		Url:    OPENAI_API_BASE,
	}
}

/**
调用gpt3.5接口
*/
func (this *ChatGptTool) ChatGPT3Dot5Turbo(messages []Gpt3Dot5Message) (string, error) {
	reqMessages := make([]openai.ChatCompletionMessage, 0)
	for _, row := range messages {
		reqMessage := openai.ChatCompletionMessage{
			Role:    row.Role,
			Content: row.Content,
			Name:    row.Name,
		}
		reqMessages = append(reqMessages, reqMessage)
	}
	resp, err := this.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    reqMessages,
			Temperature: 0,
		},
	)

	if err != nil {
		log.Println("ChatGPT3Dot5Turbo error: ", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

/**
调用gpt3.5流式接口
*/
func (this *ChatGptTool) ChatGPT3Dot5TurboStream(messages []Gpt3Dot5Message) (*openai.ChatCompletionStream, error) {
	c := this.Client
	ctx := context.Background()
	reqMessages := make([]openai.ChatCompletionMessage, 0)
	for _, row := range messages {
		reqMessage := openai.ChatCompletionMessage{
			Role:    row.Role,
			Content: row.Content,
			Name:    row.Name,
		}
		reqMessages = append(reqMessages, reqMessage)
	}
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 2000,
		Messages:  reqMessages,
		Stream:    true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return stream, err
	}

	//for {
	//	response, err := stream.Recv()
	//	if errors.Is(err, io.EOF) {
	//		log.Println("\nStream finished")
	//		break
	//	} else if err != nil {
	//		log.Printf("\nStream error: %v\n", err)
	//		break
	//	} else {
	//		log.Println(response.Choices[0].Delta.Content, err)
	//	}
	//}
	return stream, nil
}

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
		Object    string    `json:"object"`
	} `json:"data"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func (this *ChatGptTool) GetEmbedding(input string, model string) (string, error) {
	// 构建请求体
	requestBody := EmbeddingRequest{
		Input: input,
		Model: model,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 构建 HTTP 请求
	url := this.Url + "/v1/embeddings"
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBodyBytes))
	if err != nil {
		log.Println("embeddings error:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+this.Secret)

	// 发送请求并获取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("embeddings error:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 解析响应体
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("embeddings error:", err)
		return string(responseBodyBytes), err
	}
	return string(responseBodyBytes), nil
}
