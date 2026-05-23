package services

import (
	"context"
	"encoding/json"
	"strings"

	"newco-go-reporting-service/internal/ai/dto"
)

type IntentClassifier struct {
	Ollama      *OllamaService
	MemoryStore *ChatMemoryStore
}

func NewIntentClassifier(
	ollama *OllamaService,
	memoryStore *ChatMemoryStore,
) *IntentClassifier {

	return &IntentClassifier{
		Ollama:      ollama,
		MemoryStore: memoryStore,
	}
}

func (c *IntentClassifier) Classify(
	ctx context.Context,
	sessionID string,
	message string,
) (*dto.AIIntentClassification, error) {

	lowerMessage := strings.ToLower(message)

	recentTurns := c.MemoryStore.RecentTurns(sessionID)

	if len(recentTurns) > 0 {
		lastTurn := recentTurns[len(recentTurns)-1]

		if lowerMessage == "why" ||
			lowerMessage == "why?" ||
			lowerMessage == "explain" ||
			lowerMessage == "explain why" ||
			strings.Contains(lowerMessage, "what should management do") ||
			strings.Contains(lowerMessage, "redistribute staff") ||
			strings.Contains(lowerMessage, "overloaded") ||
			strings.Contains(lowerMessage, "recommend") ||
			strings.Contains(lowerMessage, "advice") {

			return &dto.AIIntentClassification{
				Intent:     lastTurn.Intent,
				ToolName:   lastTurn.ToolName,
				NeedsChart: true,
				ChartType:  "",
				Reason:     "follow-up question inherited previous intent and tool",
			}, nil
		}
	}

	if strings.Contains(lowerMessage, "variance") ||
		strings.Contains(lowerMessage, "prediction") ||
		strings.Contains(lowerMessage, "deviation") ||
		strings.Contains(lowerMessage, "actual vs") {

		return &dto.AIIntentClassification{
			Intent:     "recipe_variance",
			ToolName:   "top_recipe_variance",
			NeedsChart: true,
			ChartType:  "bar",
			Reason:     "matched recipe variance keywords",
		}, nil
	}

	if strings.Contains(lowerMessage, "ingredient") ||
		strings.Contains(lowerMessage, "protein") ||
		strings.Contains(lowerMessage, "carbohydrate") ||
		strings.Contains(lowerMessage, "oil") ||
		strings.Contains(lowerMessage, "category") {

		return &dto.AIIntentClassification{
			Intent:     "ingredient_category_usage",
			ToolName:   "ingredient_category_usage",
			NeedsChart: true,
			ChartType:  "bar",
			Reason:     "matched ingredient usage keywords",
		}, nil
	}

	if strings.Contains(lowerMessage, "site") ||
		strings.Contains(lowerMessage, "branch") ||
		strings.Contains(lowerMessage, "performing") ||
		strings.Contains(lowerMessage, "performance") ||
		strings.Contains(lowerMessage, "workload") {

		return &dto.AIIntentClassification{
			Intent:     "branch_performance",
			ToolName:   "branch_summary",
			NeedsChart: true,
			ChartType:  "bar",
			Reason:     "matched site performance keywords",
		}, nil
	}

	if strings.Contains(lowerMessage, "daily plan") ||
		strings.Contains(lowerMessage, "planning") ||
		strings.Contains(lowerMessage, "draft") ||
		strings.Contains(lowerMessage, "finalized") {

		return &dto.AIIntentClassification{
			Intent:     "daily_plan_summary",
			ToolName:   "daily_plan_summary",
			NeedsChart: true,
			ChartType:  "pie",
			Reason:     "matched daily plan keywords",
		}, nil
	}

	if strings.Contains(lowerMessage, "summary") ||
		strings.Contains(lowerMessage, "summarize") ||
		strings.Contains(lowerMessage, "overview") ||
		strings.Contains(lowerMessage, "management") ||
		strings.Contains(lowerMessage, "executive") {

		return &dto.AIIntentClassification{
			Intent:     "executive_summary",
			ToolName:   "executive_summary",
			NeedsChart: true,
			ChartType:  "line",
			Reason:     "matched executive summary keywords",
		}, nil
	}

	conversationContext := ""

	for _, turn := range recentTurns {

		conversationContext += `
Previous User Message:
` + turn.UserMessage + `

Previous Intent:
` + turn.Intent + `

Previous Tool:
` + turn.ToolName + `

Previous Assistant Response:
` + turn.AssistantResponse + `
`
	}

	systemPrompt := `
You are NewCo's AI intent classifier.

Your job is ONLY to classify the user's question into one approved intent and tool.

Return STRICT JSON only.

Do not answer the user.
Do not explain.
Do not include markdown.
Do not include code fences.
`

	userPrompt := `
Previous conversation context:
` + conversationContext + `

Classify this user message:

"` + message + `"

Approved intents and tools:

1. executive_summary
tool_name: executive_summary
Use when the user asks for a general executive overview, management summary, operational summary, or business performance summary.

2. ingredient_category_usage
tool_name: ingredient_category_usage
Use when the user asks about ingredients, ingredient categories, consumption by category, protein/carbohydrate/oil usage, or ingredient quantities.

3. recipe_variance
tool_name: top_recipe_variance
Use when the user asks about prediction accuracy, variance, recipes with high deviation, actual vs predicted, or recipe planning errors.

4. branch_performance
tool_name: branch_summary
Use when the user asks about branch/site performance, most active branch, site comparison, or branch workload.

5. daily_plan_summary
tool_name: daily_plan_summary
Use when the user asks about daily plans, finalized plans, draft plans, planning progress, or requisition planning.

6. general_advice
tool_name: none
Use when the user asks for general operational advice that does not require data.

Required JSON shape:

{
  "intent": "string",
  "tool_name": "string",
  "needs_chart": false,
  "chart_type": "",
  "reason": "string"
}
`

	content, err := c.Ollama.Chat(
		ctx,
		systemPrompt,
		userPrompt,
	)
	if err != nil {
		return nil, err
	}

	var result dto.AIIntentClassification

	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
