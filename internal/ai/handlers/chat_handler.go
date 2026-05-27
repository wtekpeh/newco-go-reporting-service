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
	PromptBuilder    *services.PromptBuilder
}

type AIStructuredResponse struct {
	Title string `json:"title"`

	Observation string `json:"observation,omitempty"`

	OperationalMeaning string `json:"operational_meaning,omitempty"`

	Recommendation string `json:"recommendation,omitempty"`

	PlanningStatus string `json:"planning_status,omitempty"`

	OperationalRisk string `json:"operational_risk,omitempty"`

	ManagementAttention string `json:"management_attention,omitempty"`

	OperationalExplanation string `json:"operational_explanation,omitempty"`

	OperationalAdvice string `json:"operational_advice,omitempty"`

	ExpectedImprovement string `json:"expected_improvement,omitempty"`

	ManagementConsideration string `json:"management_consideration,omitempty"`

	OperationalImpact string `json:"operational_impact,omitempty"`

	NextAction string `json:"next_action,omitempty"`
}

func formatStructuredAIResponse(rawContent string) string {
	cleanContent := strings.TrimSpace(rawContent)

	startIndex := strings.Index(cleanContent, "{")
	endIndex := strings.LastIndex(cleanContent, "}")

	if startIndex >= 0 && endIndex > startIndex {
		cleanContent = cleanContent[startIndex : endIndex+1]
	}

	if !strings.HasPrefix(cleanContent, "{") {
		return strings.TrimSpace(rawContent)
	}

	var structured AIStructuredResponse

	err := json.Unmarshal(
		[]byte(cleanContent),
		&structured,
	)
	if err != nil {
		return cleanContent
	}

	parts := []string{}

	if structured.Title != "" {
		parts = append(parts, structured.Title)
	}

	if structured.Observation != "" {
		parts = append(parts, "Observation: "+structured.Observation)
	}

	if structured.OperationalMeaning != "" {
		parts = append(parts, "Operational Meaning: "+structured.OperationalMeaning)
	}

	if structured.OperationalExplanation != "" {
		parts = append(parts, "Operational Explanation: "+structured.OperationalExplanation)
	}

	if structured.PlanningStatus != "" {
		parts = append(parts, "Planning Status: "+structured.PlanningStatus)
	}

	if structured.OperationalRisk != "" {
		parts = append(parts, "Operational Risk: "+structured.OperationalRisk)
	}

	if structured.ManagementAttention != "" {
		parts = append(parts, "Management Attention: "+structured.ManagementAttention)
	}

	if structured.Recommendation != "" {
		parts = append(parts, "Recommendation: "+structured.Recommendation)
	}

	if structured.OperationalImpact != "" {
		parts = append(parts, "Operational Impact: "+structured.OperationalImpact)
	}

	if structured.NextAction != "" {
		parts = append(parts, "Next Action: "+structured.NextAction)
	}

	if structured.OperationalAdvice != "" {
		parts = append(parts, "Operational Advice: "+structured.OperationalAdvice)
	}

	if structured.ExpectedImprovement != "" {
		parts = append(parts, "Expected Improvement: "+structured.ExpectedImprovement)
	}

	if structured.ManagementConsideration != "" {
		parts = append(parts, "Management Consideration: "+structured.ManagementConsideration)
	}

	return strings.Join(parts, "\n\n")
}

func NewAIChatHandler(
	intentClassifier *services.IntentClassifier,
	ollama *services.OllamaService,
	reportService *reportservices.ReportService,
	memoryStore *services.ChatMemoryStore,
	contextBuilder *services.ExecutiveContextBuilder,
	promptBuilder *services.PromptBuilder,
) *AIChatHandler {

	return &AIChatHandler{
		IntentClassifier: intentClassifier,
		Ollama:           ollama,
		ReportService:    reportService,
		MemoryStore:      memoryStore,
		ContextBuilder:   contextBuilder,
		PromptBuilder:    promptBuilder,
	}
}

func determineConversationFocus(intent string, message string) string {
	lowerMessage := strings.ToLower(message)

	switch intent {
	case "branch_performance":
		if strings.Contains(lowerMessage, "overloaded") {
			return "overloaded_site"
		}

		if strings.Contains(lowerMessage, "best") ||
			strings.Contains(lowerMessage, "performing") {
			return "best_performing_site"
		}

		if strings.Contains(lowerMessage, "inactive") {
			return "inactive_sites"
		}

		if strings.Contains(lowerMessage, "staff") ||
			strings.Contains(lowerMessage, "redistribute") {
			return "staff_distribution"
		}
	}

	return ""
}

func determineReasoningMode(message string) string {
	lowerMessage := strings.ToLower(message)

	if strings.Contains(lowerMessage, "why") ||
		strings.Contains(lowerMessage, "explain") ||
		strings.Contains(lowerMessage, "how so") ||
		strings.Contains(lowerMessage, "that so") ||
		strings.Contains(lowerMessage, "tell me more") ||
		strings.Contains(lowerMessage, "more detail") {
		return "explanation"
	}

	if strings.Contains(lowerMessage, "risk") ||
		strings.Contains(lowerMessage, "danger") ||
		strings.Contains(lowerMessage, "problem") ||
		strings.Contains(lowerMessage, "issue") ||
		strings.Contains(lowerMessage, "concern") {
		return "risk_analysis"
	}

	if strings.Contains(lowerMessage, "recommend") ||
		strings.Contains(lowerMessage, "advice") ||
		strings.Contains(lowerMessage, "what should") ||
		strings.Contains(lowerMessage, "what next") ||
		strings.Contains(lowerMessage, "what do we do") ||
		strings.Contains(lowerMessage, "management do") ||
		strings.Contains(lowerMessage, "do next") {
		return "recommendation"
	}
	if strings.Contains(lowerMessage, "redistribute") ||
		strings.Contains(lowerMessage, "staff") ||
		strings.Contains(lowerMessage, "move staff") {
		return "operational_advice"
	}

	return "detection"
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

	if classification.ToolName == "planning_risk_summary" {
		request.StartDate = ""
		request.EndDate = ""
	}

	conversationFocus := determineConversationFocus(
		classification.Intent,
		request.Message,
	)

	if conversationFocus == "" && len(recentTurns) > 0 {
		lastTurn := recentTurns[len(recentTurns)-1]
		conversationFocus = lastTurn.Focus
	}

	conversationReasoningMode := determineReasoningMode(
		request.Message,
	)

	if conversationReasoningMode == "detection" && len(recentTurns) > 0 {
		lastTurn := recentTurns[len(recentTurns)-1]

		if lastTurn.ReasoningMode != "" {
			conversationReasoningMode = lastTurn.ReasoningMode
		}
	}

	assistantResponse := "No tool execution implemented yet."

	chartSuggestions := []dto.AIChartSuggestion{}

	chartData := map[string]interface{}{}

	operationalContext := ""

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

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate branch AI response",
				"error":   err.Error(),
			})
		}

		if (assistantResponse == "" ||
			assistantResponse == "No tool execution implemented yet.") &&
			len(branchSummary) > 0 {

			topSite := branchSummary[0]

			if conversationFocus == "best_performing_site" {

				if conversationReasoningMode == "explanation" {

					assistantResponse = fmt.Sprintf(
						"%s is performing best because it recorded the highest consumption activity during the selected period, with %d consumptions compared with the other sites.",
						topSite.BranchName,
						topSite.TotalBatches,
					)

				} else {

					assistantResponse = fmt.Sprintf(
						"%s is currently performing best with %d consumptions recorded during the selected period.",
						topSite.BranchName,
						topSite.TotalBatches,
					)
				}

			} else if conversationFocus == "overloaded_site" {

				if conversationReasoningMode == "recommendation" {

					assistantResponse = fmt.Sprintf(
						"Management should review workload distribution immediately. %s is handling %d consumptions with only %d staff, while the other sites recorded little or no activity. The next step is to check whether some work or support can be moved to underutilized sites.",
						topSite.BranchName,
						topSite.TotalBatches,
						topSite.StaffCount,
					)

				} else if conversationReasoningMode == "explanation" {

					assistantResponse = fmt.Sprintf(
						"%s appears overloaded because it handled %d consumptions with only %d staff, while the other sites recorded little or no activity during the selected period.",
						topSite.BranchName,
						topSite.TotalBatches,
						topSite.StaffCount,
					)

				} else {

					assistantResponse = fmt.Sprintf(
						"%s appears operationally overloaded with %d consumptions handled by %d staff.",
						topSite.BranchName,
						topSite.TotalBatches,
						topSite.StaffCount,
					)
				}

			} else {

				assistantResponse = fmt.Sprintf(
					"%s is the most active site with %d consumptions recorded during the selected period.",
					topSite.BranchName,
					topSite.TotalBatches,
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

	if classification.ToolName == "planning_risk_summary" {

		filters := reportdto.ReportFiltersDTO{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}

		planningRiskSummary, err := h.ReportService.GetAIPlanningRiskSummary(
			filters,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch planning risk summary",
				"error":   err.Error(),
			})
		}

		jsonBytes, err := json.MarshalIndent(
			planningRiskSummary,
			"",
			"  ",
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to marshal planning risk data",
				"error":   err.Error(),
			})
		}

		systemPrompt, userPrompt := h.PromptBuilder.BuildOperationalReasoningPrompt(
			"risk_analysis",
			request.Message,
			string(jsonBytes),
		)

		content, err := h.Ollama.Chat(
			c.UserContext(),
			systemPrompt,
			userPrompt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate planning risk response",
				"error":   err.Error(),
			})
		}

		assistantResponse = formatStructuredAIResponse(content)

		chartData["planning_risk_summary"] = planningRiskSummary

		assistantResponse = fmt.Sprintf(
			"Planning Risk Summary\n\nPlanning Status: %d total plans, with %d finalized and %d still in draft.\n\nOperational Risk: %s\n\nManagement Attention: %s",
			planningRiskSummary.TotalPlans,
			planningRiskSummary.FinalizedPlans,
			planningRiskSummary.DraftPlans,
			strings.Join(planningRiskSummary.RiskReasons, " "),
			strings.Join(planningRiskSummary.ManagementAttention, " "),
		)

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "line",
			Title:     "Planning Risk Trends",
			Dataset:   "planning_risk_summary",
			XField:    "date",
			YField:    "draft_count",
		})
	}

	if classification.ToolName == "ingredient_variance_risk" {

		filters := reportdto.ReportFiltersDTO{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}

		if request.BranchID != nil {
			filters.BranchID = strconv.FormatInt(*request.BranchID, 10)
		}

		ingredientVarianceRisk, err := h.ReportService.GetAIIngredientVarianceRisk(
			filters,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to fetch ingredient variance risk",
				"error":   err.Error(),
			})
		}

		chartData["ingredient_variance_risk"] = ingredientVarianceRisk

		if len(ingredientVarianceRisk.Items) == 0 {
			assistantResponse = "No ingredient variance risk was detected from finalized consumptions with recorded actual values."
		} else {

			topItem := ingredientVarianceRisk.Items[0]

			lowerMessage := strings.ToLower(request.Message)

			isExplanationFollowUp :=
				strings.Contains(lowerMessage, "why") ||
					strings.Contains(lowerMessage, "how so") ||
					strings.Contains(lowerMessage, "explain")

			isRecommendationFollowUp :=
				strings.Contains(lowerMessage, "what should management do") ||
					strings.Contains(lowerMessage, "what should we do") ||
					strings.Contains(lowerMessage, "what do we do") ||
					strings.Contains(lowerMessage, "what next") ||
					strings.Contains(lowerMessage, "recommend") ||
					strings.Contains(lowerMessage, "recommendation")

			if isExplanationFollowUp {

				assistantResponse = fmt.Sprintf(
					"%s shows %.2f%% variance between planned and actual usage. Planned value was %.2f %s while actual usage was %.2f %s across %d finalized consumptions. This indicates that operational execution differed from planning assumptions.",
					topItem.Ingredient,
					topItem.VariancePercent,
					topItem.TotalPlannedValue,
					topItem.Unit,
					topItem.TotalActualValue,
					topItem.Unit,
					topItem.TotalConsumptions,
				)

			} else if isRecommendationFollowUp {

				assistantResponse = fmt.Sprintf(
					"Management should review %s because it has %.2f%% variance between planned and actual usage across %d finalized consumptions. The first action is to compare recent planning assumptions against actual kitchen execution, then confirm whether the variance came from portion changes, recording errors, supplier/unit differences, or recipe scaling issues. Since the risk level is %s, this should be monitored and corrected before it becomes a recurring planning control issue.",
					topItem.Ingredient,
					topItem.VariancePercent,
					topItem.TotalConsumptions,
					topItem.RiskLevel,
				)
			} else {

				assistantResponse = fmt.Sprintf(
					"Ingredient Variance Risk\n\nHighest Risk Level: %s\n\nTop Ingredient: %s has %.2f%% variance between planned and actual usage.\n\nManagement Attention: %s",
					ingredientVarianceRisk.HighestRiskLevel,
					topItem.Ingredient,
					topItem.VariancePercent,
					strings.Join(ingredientVarianceRisk.ManagementAttention, " "),
				)
			}
		}

		chartSuggestions = append(chartSuggestions, dto.AIChartSuggestion{
			ChartType: "bar",
			Title:     "Ingredient Variance Risk",
			Dataset:   "ingredient_variance_risk",
			XField:    "ingredient",
			YField:    "variance_percent",
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

		_, aiContext, err := h.ContextBuilder.BuildExecutiveSummaryContext(
			c.UserContext(),
			executiveRequest,
		)

		operationalContext = services.BuildExecutiveSummaryNarrative(
			aiContext,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to build executive summary context",
				"error":   err.Error(),
			})
		}

		if conversationReasoningMode == "explanation" {

			assistantResponse = services.BuildExecutiveSummaryExplanation(
				aiContext,
			)

		} else if conversationReasoningMode == "recommendation" {

			assistantResponse = services.BuildExecutiveSummaryRecommendation(
				aiContext,
			)

		} else {

			assistantResponse = operationalContext
		}

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

			Focus: conversationFocus,

			ReasoningMode: conversationReasoningMode,

			StartDate: request.StartDate,
			EndDate:   request.EndDate,

			BranchID: request.BranchID,

			AssistantResponse: assistantResponse,

			OperationalContext: operationalContext,

			CreatedAt: time.Now(),
		},
	)

	return c.JSON(response)
}
