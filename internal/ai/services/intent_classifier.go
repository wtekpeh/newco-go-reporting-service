package services

import (
	"context"
	"encoding/json"
	"strings"

	"newco-go-reporting-service/internal/ai/dto"
)

type IntentRule struct {
	Intent     string
	ToolName   string
	ChartType  string
	Keywords   []string
	NeedsChart bool
}

type IntentClassifier struct {
	Ollama      *OllamaService
	MemoryStore *ChatMemoryStore
}

const minimumRuleConfidence = 0.10

func NewIntentClassifier(
	ollama *OllamaService,
	memoryStore *ChatMemoryStore,
) *IntentClassifier {

	return &IntentClassifier{
		Ollama:      ollama,
		MemoryStore: memoryStore,
	}
}

var intentRules = []IntentRule{
	{
		Intent:     "recipe_variance",
		ToolName:   "top_recipe_variance",
		ChartType:  "bar",
		NeedsChart: true,
		Keywords: []string{
			"variance",
			"prediction",
			"deviation",
			"actual vs",
			"prediction error",
			"inaccurate recipe",
		},
	},

	{
		Intent:     "ingredient_category_usage",
		ToolName:   "ingredient_category_usage",
		ChartType:  "bar",
		NeedsChart: true,
		Keywords: []string{
			"ingredient",
			"protein",
			"carbohydrate",
			"oil",
			"category",
			"consumption category",
		},
	},

	{
		Intent:     "branch_performance",
		ToolName:   "branch_summary",
		ChartType:  "bar",
		NeedsChart: true,
		Keywords: []string{
			"site",
			"branch",
			"performing",
			"performance",
			"workload",
			"overloaded",
			"staff",
			"redistribute",
			"best site",
			"performing best",
			"best performing site",
			"top performing site",
			"highest performing site",
			"which site is best",
			"inactive site",
		},
	},

	{
		Intent:     "daily_plan_summary",
		ToolName:   "daily_plan_summary",
		ChartType:  "pie",
		NeedsChart: true,
		Keywords: []string{
			"daily plan",
			"planning",
			"draft",
			"finalized",
			"requisition",
		},
	},

	{
		Intent:     "planning_risk_summary",
		ToolName:   "planning_risk_summary",
		ChartType:  "line",
		NeedsChart: true,
		Keywords: []string{
			"planning risk",
			"planning issue",
			"planning problem",
			"planning delay",
			"daily plan risk",
			"requisition risk",
			"planning readiness",
			"which plans are risky",
			"operational planning risk",
			"planning concern",
			"upcoming risk",
			"upcoming plans",
			"plans need attention",
			"management attention",
			"planning concerns",
			"next few days",
			"next week",
			"this week",
			"ready for execution",
			"not ready",
		},
	},

	{
		Intent:     "executive_summary",
		ToolName:   "executive_summary",
		ChartType:  "line",
		NeedsChart: true,
		Keywords: []string{
			"summary",
			"summarize",
			"overview",
			"management",
			"executive",
			"business performance",
			"operational summary",
		},
	},
}

func calculateRuleConfidence(score int, totalKeywords int) float64 {
	if totalKeywords == 0 {
		return 0
	}

	confidence := float64(score) / float64(totalKeywords)

	if confidence > 1 {
		return 1
	}

	return confidence
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

		followUpKeywords := []string{
			"why",
			"explain",
			"what do you mean",
			"how so",
			"what should",
			"what can",
			"what next",
			"recommend",
			"advice",
			"action",
			"management do",
			"redistribute",
			"clarify",
			"that so",
			"that's so",
			"tell me more",
			"more detail",
			"what do we do",
		}

		for _, keyword := range followUpKeywords {
			if strings.Contains(lowerMessage, keyword) {
				return &dto.AIIntentClassification{
					Intent:               lastTurn.Intent,
					ToolName:             lastTurn.ToolName,
					NeedsChart:           true,
					ChartType:            "",
					Reason:               "follow-up inherited previous intent and tool",
					ConfidenceScore:      1.0,
					ClassificationSource: "follow_up_memory",
				}, nil
			}
		}
	}

	bestScore := 0
	var bestRule *IntentRule

	for i := range intentRules {
		rule := &intentRules[i]
		score := 0

		for _, keyword := range rule.Keywords {
			if strings.Contains(lowerMessage, keyword) {
				score++
			}
		}

		if score > bestScore {
			bestScore = score
			bestRule = rule
		}
	}

	if bestRule != nil && bestScore > 0 {
		confidenceScore := calculateRuleConfidence(
			bestScore,
			len(bestRule.Keywords),
		)

		if confidenceScore >= minimumRuleConfidence {
			return &dto.AIIntentClassification{
				Intent:               bestRule.Intent,
				ToolName:             bestRule.ToolName,
				NeedsChart:           bestRule.NeedsChart,
				ChartType:            bestRule.ChartType,
				Reason:               "matched scored intent rule with enough confidence",
				ConfidenceScore:      confidenceScore,
				ClassificationSource: "rule_classifier",
			}, nil
		}
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
		if len(recentTurns) > 0 {
			lastTurn := recentTurns[len(recentTurns)-1]

			return &dto.AIIntentClassification{
				Intent:               lastTurn.Intent,
				ToolName:             lastTurn.ToolName,
				NeedsChart:           true,
				ChartType:            "",
				Reason:               "LLM router returned non-JSON; inherited previous intent and tool",
				ConfidenceScore:      0.5,
				ClassificationSource: "llm_router_fallback_memory",
			}, nil
		}

		fallback := GetSafeFallbackClassification()
		return &fallback, nil
	}

	result.ConfidenceScore = 0.75
	result.ClassificationSource = "llm_router"

	if !IsApprovedAITool(result.ToolName) {
		fallback := GetSafeFallbackClassification()
		return &fallback, nil
	}

	return &result, nil
}
