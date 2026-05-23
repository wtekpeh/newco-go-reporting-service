package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"newco-go-reporting-service/internal/ai/dto"
)

type OllamaService struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewOllamaService(
	baseURL string,
	model string,
) *OllamaService {

	return &OllamaService{
		BaseURL: baseURL,
		Model:   model,
		Client:  &http.Client{},
	}
}

func (s *OllamaService) Chat(
	ctx context.Context,
	systemPrompt string,
	userPrompt string,
) (string, error) {

	userPrompt = "/no_think\nReturn only the final answer. Do not show reasoning, analysis, planning, or thinking.\n\n" + userPrompt
	think := false

	requestBody := dto.OllamaChatRequest{
		Model:  s.Model,
		Stream: false,
		Think:  &think,
		Options: map[string]interface{}{
			"num_predict": 80,
			"temperature": 0.1,
		},
		Messages: []dto.OllamaMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/chat", s.BaseURL),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"ollama returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var ollamaResponse dto.OllamaChatResponse

	err = json.Unmarshal(responseBody, &ollamaResponse)
	if err != nil {
		return "", err
	}

	return cleanAssistantContent(ollamaResponse.Message.Content), nil
}

func cleanAssistantContent(content string) string {
	markers := []string{
		"Final answer:",
		"Final Answer:",
		"Answer:",
		"ANSWER:",
	}

	for _, marker := range markers {
		if strings.Contains(content, marker) {
			parts := strings.Split(content, marker)
			return strings.TrimSpace(parts[len(parts)-1])
		}
	}

	reasoningStarters := []string{
		"Okay,",
		"Let's",
		"First,",
		"I need to",
		"Looking at",
	}

	for _, starter := range reasoningStarters {
		if strings.HasPrefix(strings.TrimSpace(content), starter) {
			return ""
		}
	}

	return strings.TrimSpace(content)
}
