package services

import "newco-go-reporting-service/internal/ai/dto"

type AIToolDefinition struct {
	Name        string
	Intent      string
	Description string
}

var ApprovedAITools = map[string]AIToolDefinition{
	"executive_summary": {
		Name:        "executive_summary",
		Intent:      "executive_summary",
		Description: "General executive operational summary using approved dashboard metrics.",
	},
	"ingredient_category_usage": {
		Name:        "ingredient_category_usage",
		Intent:      "ingredient_category_usage",
		Description: "Ingredient category usage and consumption totals.",
	},
	"top_recipe_variance": {
		Name:        "top_recipe_variance",
		Intent:      "recipe_variance",
		Description: "Recipes with highest prediction versus actual variance.",
	},
	"branch_summary": {
		Name:        "branch_summary",
		Intent:      "branch_performance",
		Description: "Branch/site performance and workload comparison.",
	},
	"daily_plan_summary": {
		Name:        "daily_plan_summary",
		Intent:      "daily_plan_summary",
		Description: "Daily plan totals, finalized plans, drafts, and planning progress.",
	},
	"site_staff_load": {
		Name:        "site_staff_load",
		Intent:      "site_staff_load",
		Description: "Staffing levels, consumptions, workload per staff member, and operational load analysis.",
	},
	"management_action_summary": {
		Name:        "management_action_summary",
		Intent:      "management_action_summary",
		Description: "Executive management actions, operational priorities, staffing concerns, planning risks, and variance risks.",
	},
	"operational_market_intelligence": {
		Name:        "operational_market_intelligence",
		Intent:      "operational_market_intelligence",
		Description: "Combines internal operational intelligence with external market intelligence to identify strategic risks, market pressures, and management actions.",
	},
	"planning_risk_summary": {
		Name:        "planning_risk_summary",
		Intent:      "planning_risk_summary",
		Description: "Planning readiness, finalized/draft daily plans, requisition risk, and management attention areas.",
	},
	"ingredient_variance_risk": {
		Name:        "ingredient_variance_risk",
		Intent:      "ingredient_variance_risk",
		Description: "Ingredient planned versus actual consumption variance risk and management attention areas.",
	},
	"internet_search": {
		Name:        "internet_search",
		Intent:      "internet_search",
		Description: "Search the internet for current market information, industry trends, prices, regulations, external risks, and business intelligence.",
	},
	"conversational": {
		Name:        "conversational",
		Intent:      "conversation",
		Description: "Greetings, acknowledgements, thanks, and general conversational replies that do not require a business intelligence tool.",
	},
	"none": {
		Name:        "none",
		Intent:      "general_advice",
		Description: "General advisory response without internal data tools.",
	},
}

func IsApprovedAITool(toolName string) bool {
	_, ok := ApprovedAITools[toolName]
	return ok
}

func GetSafeFallbackClassification() dto.AIIntentClassification {
	return dto.AIIntentClassification{
		Intent:               "general_advice",
		ToolName:             "none",
		NeedsChart:           false,
		ChartType:            "",
		Reason:               "invalid or unapproved tool returned by classifier",
		ConfidenceScore:      0,
		ClassificationSource: "safety_fallback",
	}
}
