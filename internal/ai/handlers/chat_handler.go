package handlers

import (
	"encoding/json"
	"fmt"
	"newco-go-reporting-service/internal/ai/dto"
	"newco-go-reporting-service/internal/ai/services"
	"strconv"
	"strings"
	"time"

	reportdto "newco-go-reporting-service/internal/dto"
	reportservices "newco-go-reporting-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AIChatHandler struct {
	IntentClassifier *services.IntentClassifier
	Ollama           *services.OllamaService
	ReportService    *reportservices.ReportService
	MemoryStore      *services.ChatMemoryStore
	ContextBuilder   *services.ExecutiveContextBuilder
}

func NewAIChatHandler(
	intentClassifier *services.IntentClassifier,
	ollama *services.OllamaService,
	reportService *reportservices.ReportService,
	memoryStore *services.ChatMemoryStore,
	contextBuilder *services.ExecutiveContextBuilder,
) *AIChatHandler {

	return &AIChatHandler{
		IntentClassifier: intentClassifier,
		Ollama:           ollama,
		ReportService:    reportService,
		MemoryStore:      memoryStore,
		ContextBuilder:   contextBuilder,
	}
}

func (h *AIChatHandler) Chat(c *fiber.Ctx) error {
	var request dto.AIChatRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if request.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "message is required",
		})
	}

	recentTurns := h.MemoryStore.RecentTurns(
		request.SessionID,
	)

	if len(recentTurns) > 0 {

		lastTurn := recentTurns[len(recentTurns)-1]

		if request.StartDate == "" {
			request.StartDate = lastTurn.StartDate
		}

		if request.EndDate == "" {
			request.EndDate = lastTurn.EndDate
		}

		if request.BranchID == nil {
			request.BranchID = lastTurn.BranchID
		}
	}

	classification, err := h.IntentClassifier.Classify(
		c.UserContext(),
		request.SessionID,
		request.Message,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to classify message",
			"error":   err.Error(),
		})
	}

	if !services.IsApprovedAITool(classification.ToolName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "unapproved AI tool selected",
			"tool":    classification.ToolName,
		})
	}

	assistantResponse := "No tool execution implemented yet."

	chartSuggestions := []dto.AIChartSuggestion{}

	chartData := map[string]interface{}{}

	if classification.ToolName == "top_recipe_variance" {

		branchID := ""

		if request.BranchID != nil {
			branchID = strconv.FormatInt(*request.BranchID, 10)
		}

		varianceResult, err := h.ReportService.GetTopRecipeVariance(
			c.UserContext(),
			request.StartDate,
			request.EndDate,
			branchID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch recipe variance report",
				"error":   err.Error(),
			})
		}

		jsonBytes, err := json.MarshalIndent(
			varianceResult,
			"",
			"  ",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to marshal variance data",
				"error":   err.Error(),
			})
		}

		systemPrompt := `
You are NewCo's operational AI assistant.

Do not show your reasoning.
Do not explain your thinking process.
Give only the final answer.

Use ONLY the approved operational facts provided.

Do not invent data.
Do not mention internal reasoning.
Keep the answer concise and conversational.
`

		userPrompt := `
User question:
` + request.Message + `

Approved operational facts:
` + string(jsonBytes)

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate AI response",
				"error":   err.Error(),
			})
		}

		assistantResponse = content

		chartData["top_recipe_variance"] = varianceResult

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "bar",
			Title:     "Top Recipe Variance",
			Dataset:   "top_recipe_variance",
			XField:    "recipe_name",
			YField:    "average_variance_g",
		})
	}

	if classification.ToolName == "ingredient_category_usage" {

		if request.StartDate == "" || request.EndDate == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "start_date and end_date are required for ingredient analysis",
			})
		}

		startDate, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid start_date",
				"error":   err.Error(),
			})
		}

		endDate, err := time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid end_date",
				"error":   err.Error(),
			})
		}

		ingredientReport, err := h.ReportService.GetIngredientCategoryDaily(
			c.UserContext(),
			startDate,
			endDate,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch ingredient category report",
				"error":   err.Error(),
			})
		}

		jsonBytes, err := json.MarshalIndent(
			ingredientReport,
			"",
			"  ",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to marshal ingredient report",
				"error":   err.Error(),
			})
		}

		systemPrompt := `
You are NewCo's ingredient operations assistant.

Do not show your reasoning.
Do not explain your thinking process.
Give only the final answer.

Use ONLY approved operational facts.

Do not invent numbers.
Be concise and conversational.
`

		userPrompt := `
User question:
` + request.Message + `

Interpretation rules:
- If total_actual_value is 0, assume actual usage has not yet been recorded.
- Do not treat missing actuals as operational failure automatically.
- Do not combine different units together.
- Compare g, ml, and pc separately.
- Never sum grams and milliliters together.
- Keep unit comparisons operationally accurate.

Approved operational facts:
` + string(jsonBytes)

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate ingredient AI response",
				"error":   err.Error(),
			})
		}

		assistantResponse = content

		chartData["ingredient_category_usage"] = ingredientReport

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "bar",
			Title:     "Ingredient Category Usage",
			Dataset:   "ingredient_category_usage",
			XField:    "category_name",
			YField:    "total_final_value",
		})
	}

	if classification.ToolName == "branch_summary" {

		filters := reportdto.ReportFiltersDTO{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}

		if request.BranchID != nil {
			filters.BranchID = strconv.FormatInt(*request.BranchID, 10)
		}

		branchSummary, err := h.ReportService.AIBranchPerformance(filters)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch branch summary",
				"error":   err.Error(),
			})
		}

		jsonBytes, err := json.MarshalIndent(
			branchSummary,
			"",
			"  ",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to marshal branch summary",
				"error":   err.Error(),
			})
		}

		systemPrompt := `
You are NewCo's site operations assistant.

Do not show your reasoning.
Do not explain your thinking process.
Give only the final answer.

Use ONLY approved operational facts.

Use business terminology:
- Say "consumption" instead of "batch".
- Say "site" instead of "branch".

Use correct grammar:
- Say "1 consumption" for one.
- Say "consumptions" for two or more.

Do not invent numbers.
Be concise and conversational.
`

		userPrompt := `
User question:
` + request.Message + `

Interpretation rules:
- Compare branches operationally.
- Focus on workload, batch activity, staffing, and operational balance.
- Highlight overloaded or dominant branches carefully.

Approved operational facts:
` + string(jsonBytes)

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate branch AI response",
				"error":   err.Error(),
			})
		}

		assistantResponse = content

		if assistantResponse == "" && len(branchSummary) > 0 {
			topSite := branchSummary[0]

			messageLower := strings.ToLower(request.Message)

			if strings.Contains(messageLower, "why") {

				assistantResponse = fmt.Sprintf(
					"%s is performing best because it handled %d consumptions with %d staff, while the other sites recorded little or no operational activity during the selected period.",
					topSite.BranchName,
					topSite.TotalBatches,
					topSite.StaffCount,
				)

			} else if strings.Contains(messageLower, "overloaded") {

				assistantResponse = fmt.Sprintf(
					"%s appears operationally overloaded because it is handling nearly all consumptions with only %d staff.",
					topSite.BranchName,
					topSite.StaffCount,
				)

			} else if strings.Contains(messageLower, "management") ||
				strings.Contains(messageLower, "redistribute") ||
				strings.Contains(messageLower, "recommend") ||
				strings.Contains(messageLower, "advice") {

				assistantResponse = fmt.Sprintf(
					"Management should review staff distribution across sites. %s is currently handling most operational activity with only %d staff, while other sites remain underutilized.",
					topSite.BranchName,
					topSite.StaffCount,
				)

			} else {

				assistantResponse = fmt.Sprintf(
					"%s is performing best with %d consumptions and %d staff.",
					topSite.BranchName,
					topSite.TotalBatches,
					topSite.StaffCount,
				)
			}
		}

		chartData["branch_summary"] = branchSummary

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "bar",
			Title:     "Site Performance",
			Dataset:   "branch_summary",
			XField:    "branch_name",
			YField:    "total_batches",
		})
	}

	if classification.ToolName == "daily_plan_summary" {

		filters := reportdto.ReportFiltersDTO{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}

		dailyPlanSummary, err := h.ReportService.DailyPlanSummary(
			filters,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch daily plan summary",
				"error":   err.Error(),
			})
		}

		jsonBytes, err := json.MarshalIndent(
			dailyPlanSummary,
			"",
			"  ",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to marshal daily plan summary",
				"error":   err.Error(),
			})
		}

		systemPrompt := `
You are NewCo's planning operations assistant.

Do not show your reasoning.
Do not explain your thinking process.
Give only the final answer.

Use ONLY approved operational facts.

Use business terminology:
- Say "consumption" instead of "batch".
- Say "site" instead of "branch".

Do not invent numbers.
Be concise and conversational.
`

		userPrompt := `
User question:
` + request.Message + `

Interpretation rules:
- Focus on finalized plans, draft plans, and operational planning readiness.
- Highlight operational planning risks carefully.
- Do not exaggerate missing draft plans into failures automatically.

Approved operational facts:
` + string(jsonBytes)

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate daily plan AI response",
				"error":   err.Error(),
			})
		}

		assistantResponse = content

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "pie",
			Title:     "Daily Plan Status",
			Dataset:   "daily_plan_summary",
			XField:    "status",
			YField:    "count",
		})

		chartData["daily_plan_summary"] = dailyPlanSummary
	}

	if classification.ToolName == "executive_summary" {

		executiveRequest := dto.ExecutiveSummaryRequest{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}

		if request.BranchID != nil {
			executiveRequest.BranchID = request.BranchID
		}

		executiveContext, err := h.ContextBuilder.BuildExecutiveSummaryContext(
			c.UserContext(),
			executiveRequest,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to build executive summary context",
				"error":   err.Error(),
			})
		}

		systemPrompt := `
You are NewCo's executive operations assistant.

Do not show your reasoning.
Do not explain your thinking process.
Give only the final answer.

Use ONLY approved operational facts.

Use business terminology:
- Say "consumption" instead of "batch".
- Say "site" instead of "branch".

Do not invent numbers.
Be concise, executive-level, and conversational.
`

		userPrompt := `
User question:
` + request.Message + `

Interpretation rules:
- Focus on operational performance, planning, staffing, and risk.
- Highlight important operational insights carefully.
- Do not exaggerate risks.
- If actual value is 0, treat it as not recorded yet unless explicitly marked as a finalized actual discrepancy.
- Do not call missing actuals a gap, loss, issue, variance, or failure automatically.

Approved operational facts:
` + executiveContext

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate executive AI response",
				"error":   err.Error(),
			})
		}

		assistantResponse = content

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "line",
			Title:     "Executive Operational Trends",
			Dataset:   "executive_summary",
			XField:    "date",
			YField:    "consumptions",
		})
	}

	response := dto.AIChatResponse{
		Message:           "AI chat message processed successfully",
		SessionID:         request.SessionID,
		Intent:            classification.Intent,
		UsedTools:         []string{classification.ToolName},
		AssistantResponse: assistantResponse,
		ChartSuggestions:  chartSuggestions,
		ChartData:         chartData,
		DataNotes: []string{
			"Only approved AI tools are allowed.",
			"LLM did not directly access the database.",
		},
	}

	h.MemoryStore.AddTurn(
		dto.AIConversationTurn{
			SessionID: request.SessionID,

			UserMessage: request.Message,

			Intent: classification.Intent,

			ToolName: classification.ToolName,

			StartDate: request.StartDate,
			EndDate:   request.EndDate,

			BranchID: request.BranchID,

			AssistantResponse: assistantResponse,

			CreatedAt: time.Now(),
		},
	)

	return c.JSON(response)
}
