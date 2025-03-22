package cmd

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type PalmariumProvider struct {
	Endpoint string
}

func (provider PalmariumProvider) ScrapeEvents() ([]Event, error) {
	c := colly.NewCollector()
	var events []Event

	// Select all articles with the class "event-item has-long-title"
	c.OnHTML("article.event-item.has-long-title", func(e *colly.HTMLElement) {
		// Extract the title
		title := e.ChildText("h3")

		// Extract the description
		description := e.ChildText("div.event-content-description p")

		// Extract date components
		dateString := e.ChildText("span.event-date-date") // e.g., "15.05"
		yearString := dateString[len(dateString)-4:]

		dateString = dateString[:len(dateString)-4]

		var date Date

		// Parse day and month from "15.05" format
		if parts := strings.Split(dateString, "."); len(parts) == 2 {
			if day, err := strconv.Atoi(parts[0]); err == nil {
				date.Day = day
			}
			if month, err := strconv.Atoi(parts[1]); err == nil {
				date.Month = month
			}
		}

		// Parse year
		if year, err := strconv.Atoi(yearString); err == nil {
			date.Year = year
		}

		// Create a new Event struct and append it to the events slice
		event := Event{
			Title:       title,
			Description: description,
			Date:        date,
			Location:    "Palmarium, Plantentuin",
		}
		events = append(events, event)
	})

	// Log requests being made
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	err := c.Visit(provider.Endpoint)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Constructor for PalmariumProvider
func NewPalmariumProvider() PalmariumProvider {
	return PalmariumProvider{
		Endpoint: "https://www.democrazy.be/projects/palmarium/",
	}
}
