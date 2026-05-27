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

	systemPrompt = systemPrompt + `

IMPORTANT RESPONSE RULES:

- Never expose chain-of-thought reasoning.
- Never expose internal analysis steps.
- Never narrate how you are thinking.
- Never mention "approved facts", "interpretation rules", "the user said", or "I should".
- Respond directly as an executive operations assistant.
- Use polished, concise management-ready language.
- Use bullets when summarizing several operational points.
- Never invent operational values or metrics.
`

	userPrompt = "/no_think\nReturn only the final answer. Do not show reasoning, analysis, planning, or thinking. Write a polished executive-ready response only.\n\n" + userPrompt
	think := false

	requestBody := dto.OllamaChatRequest{
		Model:  s.Model,
		Stream: false,
		Think:  &think,
		Options: map[string]interface{}{
			"num_predict": 500,
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
	cleaned := strings.TrimSpace(content)

	// Remove explicit think tags
	cleaned = strings.ReplaceAll(cleaned, "<think>", "")
	cleaned = strings.ReplaceAll(cleaned, "</think>", "")

	lines := strings.Split(content, "\n")

	reasoningMarkers := []string{
		"steps:",
		"response structure:",
		"operational facts summary:",
		"key kpis:",
		"approach:",
		"how to structure:",
		"we have to",
		"we can do:",
		"the response must",
		"using \"branch\"",
		"using \"batch\"",
	}

	lowerContent := strings.ToLower(content)

	for _, marker := range reasoningMarkers {
		index := strings.Index(lowerContent, marker)

		if index != -1 {
			content = content[index+len(marker):]
			lowerContent = strings.ToLower(content)
		}
	}

	filteredLines := []string{}

	skipPrefixes := []string{
		"okay",
		"let's",
		"first",
		"i need to",
		"looking at",
		"the user is asking",
		"the key here",
		"wait",
		"hmm",
		"so the answer should",
		"the question is",
		"the approved facts show",
		"from the data",
		"the main point",
		"we are given",
		"we must",
		"we cannot",
		"the instruction",
		"the problem says",
		"however",
		"but the user",
		"since the approved facts",
		"the approved facts do not",
		"therefore",
		"possible next steps",
		"total active users",
		"highlights:",
		"top recipe variance",
		"need to structure",
		"mention that",
		"start with",
		"the variance",
		"this is",
		"key points from",
		"- reporting period:",
		"- kpis:",
		"- highlights:",
		"- top_recipe_variance:",
		"- use business terminology",
		"- be concise",
		"- use bullets",
		"- do not invent",
		"- do not call",
		"important:",
		"we are to focus on:",
		"we have to be careful:",
		"for the top_recipe_variance",
		"the reporting period",
		"branch becomes",
		"consumption becomes",
		"highlights mention",
		"let me",
		"the ingredient categories",
		"also,",
		"so ",
		"wait,",
		"maybe ",
		"steps:",
		"operational facts summary:",
		"approach:",
		"how to structure:",
		"key kpis:",
		"3.",
		"2.",
		"1.",
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		lower := strings.ToLower(trimmed)

		shouldSkip := false

		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(lower, prefix) {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			continue
		}

		filteredLines = append(filteredLines, trimmed)
	}

	finalResponse := strings.TrimSpace(
		strings.Join(filteredLines, "\n"),
	)

	// Safety fallback
	if finalResponse == "" {
		return "I could not generate a clean operational response."
	}
	return finalResponse
}
