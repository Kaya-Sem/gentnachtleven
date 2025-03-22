package cmd

type Date struct {
	Day   int
	Month int
	Year  int
}

// Define the Event struct
type Event struct {
	Title       string
	Description string
	Date        Date
	Location    string
}

// EventProvider interface
type EventProvider interface {
	ScrapeEvents() ([]Event, error)
}
