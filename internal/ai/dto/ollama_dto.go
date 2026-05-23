package dto

type OllamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Think    *bool                  `json:"think,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Model   string        `json:"model"`
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}
