package handlers

import (
	"newco-go-reporting-service/internal/ai/dto"
	"newco-go-reporting-service/internal/ai/services"

	"github.com/gofiber/fiber/v2"
)

type ExecutiveAIHandler struct {
	Ollama *services.OllamaService
}

func NewExecutiveAIHandler(
	ollama *services.OllamaService,
) *ExecutiveAIHandler {
	return &ExecutiveAIHandler{
		Ollama: ollama,
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

	userPrompt := `
Generate a short executive summary for NewCo operations.

Current request:
Start date: ` + request.StartDate + `
End date: ` + request.EndDate + `

For now, this is a connectivity test. Summarize that the AI service is connected and ready.
`

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

	response := dto.ExecutiveSummaryResponse{
		Message: "AI executive summary generated successfully",
		Summary: content,
		KeyInsights: []string{
			"AI service connection is working.",
		},
		Risks: []string{},
		Recommendations: []string{
			"Next step is to connect approved reporting facts from Go.",
		},
		DataNotes: []string{
			"This response is a connectivity test only.",
			"The LLM did not access the database directly.",
		},
	}

	return c.JSON(response)
}
