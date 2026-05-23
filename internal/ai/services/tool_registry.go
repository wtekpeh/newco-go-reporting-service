package services

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
