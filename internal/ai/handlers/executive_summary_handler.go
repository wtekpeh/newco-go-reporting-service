package handlers

import (
	"encoding/json"
	"newco-go-reporting-service/internal/ai/dto"
	"newco-go-reporting-service/internal/ai/services"

	"github.com/gofiber/fiber/v2"
)

type ExecutiveAIHandler struct {
	Ollama         *services.OllamaService
	ContextBuilder *services.ExecutiveContextBuilder
}

func NewExecutiveAIHandler(
	ollama *services.OllamaService,
	contextBuilder *services.ExecutiveContextBuilder,
) *ExecutiveAIHandler {

	return &ExecutiveAIHandler{
		Ollama:         ollama,
		ContextBuilder: contextBuilder,
	}
}

func (h *ExecutiveAIHandler) GenerateSummary(c *fiber.Ctx) error {
	var request dto.ExecutiveSummaryRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	systemPrompt := `
You are NewCo's executive operations assistant.
Use only the facts provided by the system.
Do not mention internal reasoning.
Do not include chain-of-thought.
Respond in concise executive language.
`

	executiveContext, err := h.ContextBuilder.BuildExecutiveSummaryContext(
		c.UserContext(),
		request,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to build executive AI context",
			"error":   err.Error(),
		})
	}

	userPrompt := `
Generate a STRICT JSON response only.

Do not include markdown.
Do not include explanations.
Do not include code fences.
Do not include introductory text.

Required response format:

Return ONLY valid JSON.

Do NOT return:
- markdown
- bullet points
- explanations
- headings
- labels like "Summary:"
- code fences

The response MUST start with {
and MUST end with }

Example valid response:

{
  "summary": "Operations remained stable.",
  "key_insights": [
    "Accra Site processed the highest batches."
  ],
  "risks": [
    "One draft daily plan remains pending."
  ],
  "recommendations": [
    "Finalize pending draft daily plans."
  ]
}

Approved operational facts:

Interpretation rules:
- If total_actual is 0, assume actual usage has not yet been recorded unless explicitly stated otherwise.
- Do not automatically treat missing actual usage as a discrepancy or operational failure.
- Focus only on confirmed operational risks.

` + executiveContext

	content, err := h.Ollama.Chat(
		c.UserContext(),
		systemPrompt,
		userPrompt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate AI summary",
			"error":   err.Error(),
		})
	}

	type AIExecutiveStructuredResponse struct {
		Summary         string   `json:"summary"`
		KeyInsights     []string `json:"key_insights"`
		Risks           []string `json:"risks"`
		Recommendations []string `json:"recommendations"`
	}

	var aiResult AIExecutiveStructuredResponse

	err = json.Unmarshal([]byte(content), &aiResult)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":         "failed to parse AI response",
			"error":           err.Error(),
			"raw_ai_response": content,
		})
	}

	response := dto.ExecutiveSummaryResponse{
		Message:         "AI executive summary generated successfully",
		Summary:         aiResult.Summary,
		KeyInsights:     aiResult.KeyInsights,
		Risks:           aiResult.Risks,
		Recommendations: aiResult.Recommendations,
		DataNotes: []string{
			"Response generated from approved operational facts only.",
			"LLM did not directly access the database.",
		},
	}

	return c.JSON(response)
}
