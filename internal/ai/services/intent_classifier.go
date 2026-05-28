package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	llm         *LLMService
	MemoryStore *ChatMemoryStore
}

type LLMIntentResponse struct {
	Intent             string `json:"intent"`
	ToolName           string `json:"tool_name"`
	ReasoningMode      string `json:"reasoning_mode"`
	NeedsChart         bool   `json:"needs_chart"`
	PreferredChartType string `json:"preferred_chart_type"`
}

const minimumRuleConfidence = 0.10

func NewIntentClassifier(
	llm *LLMService,
	memoryStore *ChatMemoryStore,
) *IntentClassifier {

	return &IntentClassifier{
		llm:         llm,
		MemoryStore: memoryStore,
	}
}

var intentRules = []IntentRule{
	{
		Intent:     "ingredient_variance_risk",
		ToolName:   "ingredient_variance_risk",
		ChartType:  "bar",
		NeedsChart: true,
		Keywords: []string{
			"ingredient variance",
			"ingredient risk",
			"ingredient usage risk",
			"planned versus actual",
			"planned vs actual",
			"actual versus planned",
			"actual vs planned",
			"usage variance",
			"consumption variance",
			"variance risk",
			"which ingredient is off",
			"which ingredients are off",
			"ingredients have the highest",
			"highest planned vs actual",
		},
	},

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
		Intent:     "ingredient_variance_risk",
		ToolName:   "ingredient_variance_risk",
		ChartType:  "bar",
		NeedsChart: true,
		Keywords: []string{
			"ingredient variance",
			"ingredient risk",
			"ingredient usage risk",
			"planned versus actual",
			"planned vs actual",
			"actual versus planned",
			"actual vs planned",
			"usage variance",
			"consumption variance",
			"variance risk",
			"which ingredient is off",
			"which ingredients are off",
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

	llmClassification, err := c.classifyWithLLM(
		ctx,
		message,
		recentTurns,
	)

	if err == nil && IsApprovedAITool(llmClassification.ToolName) {
		return &dto.AIIntentClassification{
			Intent:               llmClassification.Intent,
			ToolName:             llmClassification.ToolName,
			ReasoningMode:        llmClassification.ReasoningMode,
			NeedsChart:           llmClassification.NeedsChart,
			ChartType:            "",
			Reason:               "classified by LLM router",
			ConfidenceScore:      0.95,
			ClassificationSource: "llm_router",
		}, nil
	}

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

3. ingredient_variance_risk
tool_name: ingredient_variance_risk
Use when the user asks about ingredient planned versus actual usage, ingredient variance, usage variance, ingredient control risk, or ingredients that are off.

4. recipe_variance
tool_name: top_recipe_variance
Use when the user asks about prediction accuracy, recipe variance, recipes with high deviation, actual vs predicted, or recipe planning errors.

5. branch_performance
tool_name: branch_summary
Use when the user asks about branch/site performance, most active branch, site comparison, or branch workload.

6. daily_plan_summary
tool_name: daily_plan_summary
Use when the user asks about daily plans, finalized plans, draft plans, planning progress, or requisition planning.

7. planning_risk_summary
tool_name: planning_risk_summary
Use when the user asks about planning risks, planning concerns, upcoming plan readiness, daily plan risks, or requisition risks.

8. ingredient_variance_risk
tool_name: ingredient_variance_risk
Use when the user asks about ingredient planned versus actual usage, ingredient variance, usage variance, ingredient control risk, or ingredients that are off.

9. general_advice
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

	content, err := c.llm.Chat(
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

func (c *IntentClassifier) classifyWithLLM(
	ctx context.Context,
	message string,
	recentTurns []dto.AIConversationTurn,
) (LLMIntentResponse, error) {

	conversationContext := ""

	if len(recentTurns) > 0 {

		lastTurn := recentTurns[len(recentTurns)-1]

		conversationContext = fmt.Sprintf(
			`
Previous Intent: %s
Previous Tool: %s
Previous Assistant Response:
%s
`,
			lastTurn.Intent,
			lastTurn.ToolName,
			lastTurn.AssistantResponse,
		)
	}

	systemPrompt := `
You are an operational AI routing assistant.

Your job is ONLY to determine:
- operational intent
- approved tool
- reasoning mode
- whether a chart is useful
- preferred chart type

You MUST ONLY return valid JSON.

Approved tools:
- executive_summary
- branch_summary
- ingredient_variance_risk
- planning_risk_summary
- none

Reasoning modes:
- summary
- explanation
- recommendation

Chart types:
- bar
- line
- pie
- none

If the user asks for:
- chart
- graph
- bar graph
- line graph
- visualize
- compare visually

then:
- needs_chart should be true
- choose the best chart type

Examples:

User:
"Which site is overloaded?"

Response:
{
  "intent": "branch_performance",
  "tool_name": "branch_summary",
  "reasoning_mode": "summary",
  "needs_chart": true,
  "preferred_chart_type": "bar"
}

User:
"why so"

Response:
{
  "intent": "follow_up",
  "tool_name": "branch_summary",
  "reasoning_mode": "explanation",
  "needs_chart": false
}

Return ONLY JSON.
`

	userPrompt := fmt.Sprintf(
		`
Conversation Context:
%s

Current User Message:
%s
`,
		conversationContext,
		message,
	)

	response, err := c.llm.Chat(
		ctx,
		systemPrompt,
		userPrompt,
	)

	if err != nil {
		return LLMIntentResponse{}, err
	}

	var parsedResponse LLMIntentResponse

	err = json.Unmarshal(
		[]byte(response),
		&parsedResponse,
	)

	if err != nil {
		return LLMIntentResponse{}, err
	}

	return parsedResponse, nil
}
