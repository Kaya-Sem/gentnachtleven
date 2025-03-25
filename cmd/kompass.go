package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type KompassProvider struct {
	Endpoint string
}

func (provider KompassProvider) ScrapeEvents() ([]Event, error) {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.AllowedDomains("kompassklub.com"),
	)

	// Use a map to track unique events and visited URLs
	uniqueEvents := make(map[string]Event)
	visitedURLs := make(map[string]bool)

	// Find all event URLs on the listing page
	c.OnHTML("div.jet-listing-grid__item a", func(e *colly.HTMLElement) {
		eventURL := e.Attr("href")
		if eventURL == "" || visitedURLs[eventURL] {
			return
		}

		// Mark URL as visited
		visitedURLs[eventURL] = true

		// Create a new collector for each event detail page
		detailCollector := c.Clone()
		detailCollector.OnHTML("body", func(detail *colly.HTMLElement) {
			// Extract title
			title := detail.ChildText("h1.elementor-heading-title.elementor-size-default")
			if title == "" {
				title = "no title found"
			}

			// Extract date
			dateText := detail.ChildText("div.jet-listing-dynamic-field__content")
			parsedDate := parseDate(dateText)

			// Extract description
			description := detail.ChildText(".event-description, .tribe-events-content")

			// Create a unique key for the event to prevent duplicates
			eventKey := fmt.Sprintf("%s_%d_%d_%d", title, parsedDate.Year, parsedDate.Month, parsedDate.Day)

			// Create event
			event := Event{
				Title:       strings.TrimSpace(title),
				Description: strings.TrimSpace(description),
				Location:    "Kompass Klub",
				Date:        parsedDate,
			}

			// Add to unique events only if not already present
			if _, exists := uniqueEvents[eventKey]; !exists {
				uniqueEvents[eventKey] = event
				fmt.Printf("Added event: %s on %d/%d/%d\n", title, parsedDate.Day, parsedDate.Month, parsedDate.Year)
			}
		})

		// Visit the detail page
		detailCollector.Visit(eventURL)
	})

	// Start scraping the main page
	err := c.Visit(provider.Endpoint)
	if err != nil {
		return nil, err
	}

	// Convert unique events map to slice
	events := make([]Event, 0, len(uniqueEvents))
	for _, event := range uniqueEvents {
		events = append(events, event)
	}

	return events, nil

}
func parseDate(dateText string) Date {
	// Trim any whitespace
	dateText = strings.TrimSpace(dateText)

	// List of potential date formats to try
	formats := []string{
		"2 January 2006",
		"2 Jan 2006",
		"Jan 2, 2006",
		"2006-01-02",
		"02-01-2006",
		"01/02/2006",
	}

	// Custom month name mapping to handle potential variations
	monthMapping := map[string]string{
		"January":  "January",
		"Jan":      "January",
		"February": "February",
		"Feb":      "February",
		"March":    "March",
		"Mar":      "March",
		// Add more mappings as needed
	}

	// First, try standard time.Parse with various formats
	for _, format := range formats {
		t, err := time.Parse(format, dateText)
		if err == nil {
			return Date{
				Day:   t.Day(),
				Month: int(t.Month()),
				Year:  t.Year(),
			}
		}
	}

	// If standard parsing fails, try a more custom approach
	parts := strings.Fields(dateText)
	if len(parts) >= 3 {
		// Try to parse day, month name, and year
		dayStr := parts[0]
		monthStr := parts[1]
		yearStr := parts[2]

		// Normalize month name
		normalizedMonth, ok := monthMapping[monthStr]
		if !ok {
			normalizedMonth = monthStr
		}

		// Combine the parts in a standard format
		combinedDateStr := fmt.Sprintf("%s %s %s", dayStr, normalizedMonth, yearStr)

		t, err := time.Parse("2 January 2006", combinedDateStr)
		if err == nil {
			return Date{
				Day:   t.Day(),
				Month: int(t.Month()),
				Year:  t.Year(),
			}
		}
	}

	// Fallback to current date if parsing fails
	now := time.Now()
	return Date{
		Day:   now.Day(),
		Month: int(now.Month()),
		Year:  now.Year(),
	}
}

// Constructor for KompassProvider
func NewKompassProvider() KompassProvider {
	return KompassProvider{
		Endpoint: "https://kompassklub.com/event-list/",
	}
}
