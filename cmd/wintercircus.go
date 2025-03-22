package cmd

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type ClubWinterCircusProvider struct {
	Endpoint string
}

func (provider ClubWinterCircusProvider) ScrapeEvents() ([]Event, error) {
	c := colly.NewCollector()
	var events []Event

	// Select all list items that contain events
	c.OnHTML("li.wp-block-post", func(e *colly.HTMLElement) {
		// Extract the title
		title := e.ChildText("h2.list_title")

		// Extract date components
		dateString := e.ChildText("div.post-date-container time")

		var date Date

		// Parse date from format like "22 mrt"
		if parts := strings.Split(dateString, " "); len(parts) == 2 {
			// Parse day
			if day, err := strconv.Atoi(parts[0]); err == nil {
				date.Day = day
			}

			// Parse month from Dutch abbreviation
			month := parseMonthDutch(parts[1])
			date.Month = month

			// Set year from datetime attribute
			datetimeAttr := e.ChildAttr("time", "datetime")
			if len(datetimeAttr) >= 4 {
				yearStr := datetimeAttr[:4]
				if year, err := strconv.Atoi(yearStr); err == nil {
					date.Year = year
				}
			}
		}

		// Create a new Event struct and append it to the events slice
		event := Event{
			Title:       title,
			Description: "", // Would need a second request to get description
			Date:        date,
			Location:    "Club Wintercircus",
		}

		// If we want to fetch full description, we could make a second request to the detail page
		// This would require an additional collector or using the link variable above

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

// Helper function to parse Dutch month abbreviations
func parseMonthDutch(monthStr string) int {
	switch strings.ToLower(monthStr) {
	case "jan":
		return 1
	case "feb":
		return 2
	case "mrt":
		return 3
	case "apr":
		return 4
	case "mei":
		return 5
	case "jun":
		return 6
	case "jul":
		return 7
	case "aug":
		return 8
	case "sep":
		return 9
	case "okt":
		return 10
	case "nov":
		return 11
	case "dec":
		return 12
	default:
		return 0
	}
}

// Constructor for ClubWinterCircusProvider
func NewClubWinterCircusProvider() ClubWinterCircusProvider {
	return ClubWinterCircusProvider{
		Endpoint: "https://www.clubwintercircus.be/",
	}
}
