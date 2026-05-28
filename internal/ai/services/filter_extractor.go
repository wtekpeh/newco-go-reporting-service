package services

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExtractedAIFilters struct {
	StartDate string
	EndDate   string
}

func ExtractFiltersFromMessage(
	message string,
	now time.Time,
) ExtractedAIFilters {

	lowerMessage := strings.ToLower(message)

	dateRangeRegex := regexp.MustCompile(
		`(\d{4}-\d{2}-\d{2})`,
	)

	matches := dateRangeRegex.FindAllString(
		message,
		-1,
	)

	if len(matches) >= 2 {

		return ExtractedAIFilters{
			StartDate: matches[0],
			EndDate:   matches[1],
		}
	}

	if strings.Contains(lowerMessage, "today") {
		date := now.Format("2006-01-02")

		return ExtractedAIFilters{
			StartDate: date,
			EndDate:   date,
		}
	}

	if strings.Contains(lowerMessage, "yesterday") {
		date := now.AddDate(0, 0, -1).Format("2006-01-02")

		return ExtractedAIFilters{
			StartDate: date,
			EndDate:   date,
		}
	}

	if strings.Contains(lowerMessage, "this week") {

		weekday := int(now.Weekday())

		if weekday == 0 {
			weekday = 7
		}

		startOfWeek := now.AddDate(
			0,
			0,
			-(weekday - 1),
		)

		return ExtractedAIFilters{
			StartDate: startOfWeek.Format("2006-01-02"),
			EndDate:   now.Format("2006-01-02"),
		}
	}

	if strings.Contains(lowerMessage, "last week") {

		weekday := int(now.Weekday())

		if weekday == 0 {
			weekday = 7
		}

		startOfCurrentWeek := now.AddDate(
			0,
			0,
			-(weekday - 1),
		)

		startOfLastWeek := startOfCurrentWeek.AddDate(
			0,
			0,
			-7,
		)

		endOfLastWeek := startOfCurrentWeek.AddDate(
			0,
			0,
			-1,
		)

		return ExtractedAIFilters{
			StartDate: startOfLastWeek.Format("2006-01-02"),
			EndDate:   endOfLastWeek.Format("2006-01-02"),
		}
	}

	if strings.Contains(lowerMessage, "this month") {
		startDate := time.Date(
			now.Year(),
			now.Month(),
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)

		return ExtractedAIFilters{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   now.Format("2006-01-02"),
		}
	}

	lastNMonthsRegex := regexp.MustCompile(
		`last\s+(\d+)\s+months`,
	)

	lastNMonthsMatches := lastNMonthsRegex.FindStringSubmatch(
		lowerMessage,
	)

	if len(lastNMonthsMatches) > 1 {

		monthsBack, err := strconv.Atoi(
			lastNMonthsMatches[1],
		)

		if err == nil {

			startDate := now.AddDate(
				0,
				-monthsBack,
				0,
			)

			return ExtractedAIFilters{
				StartDate: startDate.Format("2006-01-02"),
				EndDate:   now.Format("2006-01-02"),
			}
		}
	}

	if strings.Contains(lowerMessage, "last month") {
		firstDayThisMonth := time.Date(
			now.Year(),
			now.Month(),
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)

		lastDayLastMonth := firstDayThisMonth.AddDate(0, 0, -1)

		firstDayLastMonth := time.Date(
			lastDayLastMonth.Year(),
			lastDayLastMonth.Month(),
			1,
			0,
			0,
			0,
			0,
			now.Location(),
		)

		return ExtractedAIFilters{
			StartDate: firstDayLastMonth.Format("2006-01-02"),
			EndDate:   lastDayLastMonth.Format("2006-01-02"),
		}
	}

	return ExtractedAIFilters{}
}
